package pdu

import (
	"testing"
)

func TestReadSubmitSMResp(t *testing.T) {
	data := WriteHeader(Header{
		CommandLength:  0,
		CommandID:      SUBMIT_SM_RESP,
		CommandStatus:  0,
		SequenceNumber: 1,
	})
	data = append(data, []byte("msg-001\x00")...)

	got, err := ReadSubmitSMResp(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := SubmitSMResp{
		Header:    Header{CommandLength: 0, CommandID: SUBMIT_SM_RESP, CommandStatus: 0, SequenceNumber: 1},
		MessageID: "msg-001",
	}
	if got != want {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

func TestReadWriteSubmitSMResp(t *testing.T) {
	submitSMResp := SubmitSMResp{
		Header:    Header{CommandLength: 0, CommandID: SUBMIT_SM_RESP, CommandStatus: 0, SequenceNumber: 1},
		MessageID: "msg-001",
	}

	data := WriteSubmitSMResp(submitSMResp)
	got, err := ReadSubmitSMResp(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != submitSMResp {
		t.Fatalf("expected %+v, got %+v", submitSMResp, got)
	}
}

func TestReadSubmitSMResp_TooShort(t *testing.T) {
	_, err := ReadSubmitSMResp([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
