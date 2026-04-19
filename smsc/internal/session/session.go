package session

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/joeygalvin/smpp-sandbox/smsc/internal/pdu"
	"github.com/joeygalvin/smpp-sandbox/smsc/internal/store"
)

var messageIDCounter atomic.Uint64

type Session struct {
	State          State
	Conn           net.Conn
	SystemID       string
	SequenceNumber uint32
	store          *store.Store
	sessionID      int64
}

func generateMessageID() string {
	return fmt.Sprintf("%d", messageIDCounter.Add(1))
}

func NewSession(conn net.Conn, store *store.Store) *Session {
	return &Session{
		State: OPEN,
		Conn:  conn,
		store: store,
	}
}

func (s *Session) Handle(raw []byte) error {
	header, err := pdu.ReadHeader(raw)
	if err != nil {
		return err
	}

	switch header.CommandID {
	case pdu.BIND_RECEIVER:
		return s.handleBind(header, raw, BOUND_RX, pdu.BIND_RECEIVER_RESP)
	case pdu.BIND_TRANSMITTER:
		return s.handleBind(header, raw, BOUND_TX, pdu.BIND_TRANSMITTER_RESP)
	case pdu.BIND_TRANSCEIVER:
		return s.handleBind(header, raw, BOUND_TRX, pdu.BIND_TRANSCEIVER_RESP)
	case pdu.SUBMIT_SM:
		return s.handleSubmitSM(header, raw)
	case pdu.DELIVER_SM_RESP:
		return s.handleDeliverSMResp(header, raw)
	case pdu.ENQUIRE_LINK:
		return s.handleEnquireLink(header)
	case pdu.UNBIND:
		return s.handleUnbind(header)
	default:
		return s.handleDefault(header)
	}
}

func (s *Session) writeGenericNack(sequenceNumber uint32, commandStatus uint32) error {
	_, err := s.Conn.Write(pdu.NewGenericNack(sequenceNumber, commandStatus))
	return err
}

func (s *Session) handleBind(header pdu.Header, raw []byte, state State, commandID uint32) error {
	if s.State != OPEN {
		return s.writeGenericNack(header.SequenceNumber, pdu.ESME_RINVBNDSTS)
	}
	bindRequest, err := pdu.ReadBindRequest(raw)
	if err != nil {
		return err
	}
	s.SystemID = bindRequest.SystemID
	s.State = state

	id, err := s.store.InsertSession(s.SystemID, s.Conn.RemoteAddr().String(), bindTypeString(state))
	if err != nil {
		return err
	}
	s.sessionID = id

	bindResponse := pdu.BindResponse{
		Header: pdu.Header{
			CommandLength:  uint32(pdu.HEADER_LENGTH + len(s.SystemID) + 1),
			CommandID:      commandID,
			CommandStatus:  pdu.ESME_ROK,
			SequenceNumber: bindRequest.Header.SequenceNumber,
		},
		SystemID: s.SystemID,
	}
	_, err = s.Conn.Write(pdu.WriteBindResponse(bindResponse))
	return err
}

func (s *Session) handleDeliverSMResp(header pdu.Header, raw []byte) error {
	if s.State != BOUND_RX && s.State != BOUND_TRX {
		return s.writeGenericNack(header.SequenceNumber, pdu.ESME_RINVBNDSTS)
	}
	_, err := pdu.ReadDeliverSMResp(raw)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) handleSubmitSM(header pdu.Header, raw []byte) error {
	if s.State != BOUND_TX && s.State != BOUND_TRX {
		return s.writeGenericNack(header.SequenceNumber, pdu.ESME_RINVBNDSTS)
	}
	submitSM, err := pdu.ReadSubmitSM(raw)
	if err != nil {
		return err
	}

	messageID := generateMessageID()
	drRequested := submitSM.RegisteredDelivery == 0x01

	if err := s.store.InsertMessage(s.sessionID, messageID, s.SystemID, submitSM.SourceAddr, submitSM.DestAddr, string(submitSM.Message), drRequested); err != nil {
		return err
	}
	log.Printf("submit_sm from %s: to=%s msg=%s id=%s dr_requested=%t", s.SystemID, submitSM.DestAddr, string(submitSM.Message), messageID, drRequested)

	_, err = s.Conn.Write(pdu.NewSubmitSMResp(header.SequenceNumber, pdu.ESME_ROK, messageID))
	if err != nil {
		return err
	}

	if drRequested {
		s.SequenceNumber++
		_, err = s.Conn.Write(pdu.NewDeliverSM(s.SequenceNumber, submitSM, messageID))
		if err != nil {
			return err
		}
		if err := s.store.MarkDelivered(messageID); err != nil {
			return err
		}
		log.Printf("deliver_sm to %s: to=%s msg=%s id=%s", s.SystemID, submitSM.DestAddr, string(submitSM.Message), messageID)
	}

	return nil
}

func (s *Session) handleEnquireLink(header pdu.Header) error {
	if !isBound(s.State) {
		return s.writeGenericNack(header.SequenceNumber, pdu.ESME_RINVBNDSTS)
	}
	_, err := s.Conn.Write(pdu.NewEnquireLinkResp(header.SequenceNumber))
	return err
}

func (s *Session) handleUnbind(header pdu.Header) error {
	if !isBound(s.State) {
		return s.writeGenericNack(header.SequenceNumber, pdu.ESME_RINVBNDSTS)
	}
	s.State = UNBOUND
	_, err := s.Conn.Write(pdu.NewUnbindResp(header.SequenceNumber, pdu.ESME_ROK))
	if s.sessionID != 0 {
		s.store.CloseSession(s.sessionID)
	}
	s.Conn.Close()
	return err
}

func (s *Session) handleDefault(header pdu.Header) error {
	return s.writeGenericNack(header.SequenceNumber, pdu.ESME_RINVCMDID)
}

func bindTypeString(state State) string {
	switch state {
	case BOUND_TX:
		return "TX"
	case BOUND_RX:
		return "RX"
	case BOUND_TRX:
		return "TRX"
	default:
		return ""
	}
}
