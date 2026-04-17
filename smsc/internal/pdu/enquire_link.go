package pdu

import "errors"

type EnquireLink struct {
	Header Header
}

func ReadEnquireLink(data []byte) (EnquireLink, error) {
	if len(data) < HEADER_LENGTH {
		return EnquireLink{}, errors.New("data too short to contain a enquire_link PDU")
	}
	header, err := ReadHeader(data[0:HEADER_LENGTH])
	if err != nil {
		return EnquireLink{}, err
	}

	return EnquireLink{
		Header: header,
	}, nil
}

func WriteEnquireLink(enquireLink EnquireLink) []byte {
	data := WriteHeader(enquireLink.Header)
	return data
}
