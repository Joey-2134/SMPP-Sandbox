package pdu

import (
	"encoding/binary"
	"errors"
)

const HEADER_LENGTH = 16

type Header struct {
	CommandLength  uint32
	CommandID      uint32
	CommandStatus  uint32
	SequenceNumber uint32
}

func ReadHeader(data []byte) (Header, error) {
	if len(data) < HEADER_LENGTH {
		return Header{}, errors.New("data too short to contain a header")
	}
	header := Header{
		CommandLength:  binary.BigEndian.Uint32(data[0:4]),
		CommandID:      binary.BigEndian.Uint32(data[4:8]),
		CommandStatus:  binary.BigEndian.Uint32(data[8:12]),
		SequenceNumber: binary.BigEndian.Uint32(data[12:16]),
	}
	return header, nil
}

func WriteHeader(header Header) []byte {
	data := make([]byte, HEADER_LENGTH)
	binary.BigEndian.PutUint32(data[0:4], header.CommandLength)
	binary.BigEndian.PutUint32(data[4:8], header.CommandID)
	binary.BigEndian.PutUint32(data[8:12], header.CommandStatus)
	binary.BigEndian.PutUint32(data[12:16], header.SequenceNumber)
	return data
}
