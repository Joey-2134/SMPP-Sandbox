package pdu

import "errors"

// see SMPP 3.4 spec section 4.1.4
type BindReceiverResp struct {
	Header   Header
	SystemId string
}

func ReadBindReceiverResp(data []byte) (BindReceiverResp, error) {
	if len(data) < 16 {
		return BindReceiverResp{}, errors.New("data too short to contain a bind_receiver_resp PDU")
	}
	header, err := ReadHeader(data[0:16])
	if err != nil {
		return BindReceiverResp{}, err
	}

	systemId, _, err := ReadCOctetString(data[16:])
	if err != nil {
		return BindReceiverResp{}, err
	}

	return BindReceiverResp{
		Header:   header,
		SystemId: systemId,
	}, nil
}

func WriteBindReceiverResp(bindReceiverResp BindReceiverResp) []byte {
	data := WriteHeader(bindReceiverResp.Header)
	data = append(data, []byte(bindReceiverResp.SystemId)...)
	data = append(data, 0x00)
	return data
}
