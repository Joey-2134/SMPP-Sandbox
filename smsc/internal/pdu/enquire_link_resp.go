package pdu

import "errors"

type EnquireLinkResp struct {
	Header Header
}

func ReadEnquireLinkResp(data []byte) (EnquireLinkResp, error) {
	if len(data) < HEADER_LENGTH {
		return EnquireLinkResp{}, errors.New("data too short to contain a enquire_link_resp PDU")
	}
	header, err := ReadHeader(data[0:HEADER_LENGTH])
	if err != nil {
		return EnquireLinkResp{}, err
	}

	return EnquireLinkResp{
		Header: header,
	}, nil
}

func WriteEnquireLinkResp(enquireLinkResp EnquireLinkResp) []byte {
	data := WriteHeader(enquireLinkResp.Header)
	return data
}

func NewEnquireLinkResp(sequenceNumber uint32) []byte {
	enquireLinkResp := EnquireLinkResp{
		Header: Header{
			CommandLength:  HEADER_LENGTH,
			CommandID:      ENQUIRE_LINK_RESP,
			CommandStatus:  ESME_ROK,
			SequenceNumber: sequenceNumber,
		},
	}
	return WriteEnquireLinkResp(enquireLinkResp)
}
