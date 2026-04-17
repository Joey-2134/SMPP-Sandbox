package pdu

import "errors"

// see SMPP 3.4 spec section 4.5.2
type DeliverSMResp struct {
	Header    Header
	MessageID string // max 65
}

func ReadDeliverSMResp(data []byte) (DeliverSMResp, error) {
	if len(data) < HEADER_LENGTH {
		return DeliverSMResp{}, errors.New("data too short to contain a deliver_sm_resp PDU")
	}
	header, err := ReadHeader(data[0:HEADER_LENGTH])
	if err != nil {
		return DeliverSMResp{}, err
	}

	messageID, _, err := ReadCOctetString(data[HEADER_LENGTH:])
	if err != nil {
		return DeliverSMResp{}, err
	}

	return DeliverSMResp{
		Header:    header,
		MessageID: messageID,
	}, nil
}

func WriteDeliverSMResp(s DeliverSMResp) []byte {
	data := WriteHeader(s.Header)
	data = append(data, []byte(s.MessageID)...)
	data = append(data, 0x00)
	return data
}
