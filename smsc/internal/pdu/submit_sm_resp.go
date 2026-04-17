package pdu

import "errors"

// see SMPP 3.4 spec section 4.4.2
type SubmitSMResp struct {
	Header    Header
	MessageID string // max 65
}

func ReadSubmitSMResp(data []byte) (SubmitSMResp, error) {
	if len(data) < HEADER_LENGTH {
		return SubmitSMResp{}, errors.New("data too short to contain a submit_sm_resp PDU")
	}
	header, err := ReadHeader(data[0:HEADER_LENGTH])
	if err != nil {
		return SubmitSMResp{}, err
	}

	messageID, _, err := ReadCOctetString(data[HEADER_LENGTH:])
	if err != nil {
		return SubmitSMResp{}, err
	}

	return SubmitSMResp{
		Header:    header,
		MessageID: messageID,
	}, nil
}

func WriteSubmitSMResp(s SubmitSMResp) []byte {
	data := WriteHeader(s.Header)
	data = append(data, []byte(s.MessageID)...)
	data = append(data, 0x00)
	return data
}

func NewSubmitSMResp(sequenceNumber uint32, commandStatus uint32, messageID string) []byte {
	submitSMResp := SubmitSMResp{
		Header: Header{
			CommandLength:  HEADER_LENGTH + uint32(len(messageID)) + 1,
			CommandID:      SUBMIT_SM_RESP,
			CommandStatus:  commandStatus,
			SequenceNumber: sequenceNumber,
		},
		MessageID: messageID,
	}
	return WriteSubmitSMResp(submitSMResp)
}
