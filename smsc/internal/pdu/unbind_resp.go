package pdu

import "errors"

type UnbindResp struct {
	Header Header
}

func ReadUnbindResp(data []byte) (UnbindResp, error) {
	if len(data) < HEADER_LENGTH {
		return UnbindResp{}, errors.New("data too short to contain a unbind_resp PDU")
	}
	header, err := ReadHeader(data[0:HEADER_LENGTH])
	if err != nil {
		return UnbindResp{}, err
	}

	return UnbindResp{
		Header: header,
	}, nil
}

func WriteUnbindResp(unbindResp UnbindResp) []byte {
	data := WriteHeader(unbindResp.Header)
	return data
}

func NewUnbindResp(sequenceNumber uint32, commandStatus uint32) []byte {
	unbindResp := UnbindResp{
		Header: Header{
			CommandLength:  HEADER_LENGTH,
			CommandID:      UNBIND_RESP,
			CommandStatus:  commandStatus,
			SequenceNumber: sequenceNumber,
		},
	}
	return WriteUnbindResp(unbindResp)
}
