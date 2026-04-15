package pdu

import (
	"testing"
)

func TestReadDeliverSMResp(t *testing.T) {
	data := WriteHeader(Header{
		CommandLength:  0,
		CommandID:      DELIVER_SM_RESP,
		CommandStatus:  0,
		SequenceNumber: 1,
	})
	data = append(data, []byte("msg-001\x00")...)

	got, err := ReadDeliverSMResp(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := DeliverSMResp{
		Header:    Header{CommandLength: 0, CommandID: DELIVER_SM_RESP, CommandStatus: 0, SequenceNumber: 1},
		MessageID: "msg-001",
	}
	if got != want {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

func TestReadWriteDeliverSMResp(t *testing.T) {
	deliverSMResp := DeliverSMResp{
		Header:    Header{CommandLength: 0, CommandID: DELIVER_SM_RESP, CommandStatus: 0, SequenceNumber: 1},
		MessageID: "msg-001",
	}

	data := WriteDeliverSMResp(deliverSMResp)
	got, err := ReadDeliverSMResp(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != deliverSMResp {
		t.Fatalf("expected %+v, got %+v", deliverSMResp, got)
	}
}

func TestReadDeliverSMResp_TooShort(t *testing.T) {
	_, err := ReadDeliverSMResp([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
