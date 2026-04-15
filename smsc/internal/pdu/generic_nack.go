package pdu

import "errors"

type GenericNack struct {
	Header Header
}

func ReadGenericNack(data []byte) (GenericNack, error) {
	if len(data) < 16 {
		return GenericNack{}, errors.New("data too short to contain a generic_nack PDU")
	}
	header, err := ReadHeader(data[0:16])
	if err != nil {
		return GenericNack{}, err
	}

	return GenericNack{
		Header: header,
	}, nil
}

func WriteGenericNack(genericNack GenericNack) []byte {
	data := WriteHeader(genericNack.Header)
	return data
}
