package pdu

import (
	"testing"
)

func TestReadBindReceiverResp(t *testing.T) {
	data := []byte{}
	data = append(data, WriteHeader(Header{
		CommandLength:  0,
		CommandID:      BIND_RECEIVER,
		CommandStatus:  0,
		SequenceNumber: 1,
	})...)
	data = append(data, []byte("testSystemId\x00")...)

	got, err := ReadBindReceiverResp(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := BindReceiverResp{
		Header:   Header{CommandLength: 0, CommandID: BIND_RECEIVER, CommandStatus: 0, SequenceNumber: 1},
		SystemId: "testSystemId",
	}
	if got != want {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

func TestReadWriteBindReceiverResp(t *testing.T) {
	bindReceiverResp := BindReceiverResp{
		Header:   Header{CommandLength: 0, CommandID: BIND_RECEIVER, CommandStatus: 0, SequenceNumber: 1},
		SystemId: "testSystemId",
	}
	data := WriteBindReceiverResp(bindReceiverResp)
	got, err := ReadBindReceiverResp(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != bindReceiverResp {
		t.Fatalf("expected %+v, got %+v", bindReceiverResp, got)
	}
}
