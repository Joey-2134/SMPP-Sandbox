package pdu

import "errors"

// see SMPP 3.4 spec section 4.1
type BindResponse struct {
	Header   Header
	SystemID string
}

func ReadBindResponse(data []byte) (BindResponse, error) {
	if len(data) < 16 {
		return BindResponse{}, errors.New("data too short to contain a bind_response PDU")
	}
	header, err := ReadHeader(data[0:16])
	if err != nil {
		return BindResponse{}, err
	}

	systemID, _, err := ReadCOctetString(data[16:])
	if err != nil {
		return BindResponse{}, err
	}

	return BindResponse{
		Header:   header,
		SystemID: systemID,
	}, nil
}

func WriteBindResponse(bindResponse BindResponse) []byte {
	data := WriteHeader(bindResponse.Header)
	data = append(data, []byte(bindResponse.SystemID)...)
	data = append(data, 0x00)
	return data
}
