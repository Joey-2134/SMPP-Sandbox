package pdu

import "errors"

type GenericNack struct {
	Header Header
}

func ReadGenericNack(data []byte) (GenericNack, error) {
	if len(data) < HEADER_LENGTH {
		return GenericNack{}, errors.New("data too short to contain a generic_nack PDU")
	}
	header, err := ReadHeader(data[0:HEADER_LENGTH])
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

func NewGenericNack(sequenceNumber uint32, commandStatus uint32) []byte {
	genericNack := GenericNack{
		Header: Header{
			CommandLength:  HEADER_LENGTH,
			CommandID:      GENERIC_NACK,
			CommandStatus:  commandStatus,
			SequenceNumber: sequenceNumber,
		},
	}
	return WriteGenericNack(genericNack)
}
