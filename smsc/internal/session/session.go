package session

import (
	"net"

	"github.com/joeygalvin/smpp-sandbox/smsc/internal/pdu"
)

type Session struct {
	State          State
	Conn           net.Conn
	SystemID       string
	SequenceNumber uint32
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		State: OPEN,
		Conn:  conn,
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

func (s *Session) handleDefault(header pdu.Header) error {
	return s.writeGenericNack(header.SequenceNumber, pdu.ESME_RINVCMDID)
}
