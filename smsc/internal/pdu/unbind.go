package pdu

import "errors"

type Unbind struct {
	Header Header
}

func ReadUnbind(data []byte) (Unbind, error) {
	if len(data) < 16 {
		return Unbind{}, errors.New("data too short to contain a unbind PDU")
	}
	header, err := ReadHeader(data[0:16])
	if err != nil {
		return Unbind{}, err
	}

	return Unbind{
		Header: header,
	}, nil
}

func WriteUnbind(unbind Unbind) []byte {
	data := WriteHeader(unbind.Header)
	return data
}
