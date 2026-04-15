package pdu

import (
	"testing"
)

func TestReadBindReceiver(t *testing.T) {
	data := []byte{}
	data = append(data, WriteHeader(Header{
		CommandLength:  0,
		CommandID:      BIND_RECEIVER,
		CommandStatus:  0,
		SequenceNumber: 1,
	})...)
	data = append(data, []byte("testSystemId\x00")...)
	data = append(data, []byte("password\x00")...)
	data = append(data, []byte("type\x00")...)
	data = append(data, 0x34, 0x00, 0x00, 0x00)

	got, err := ReadBindReceiver(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := BindReceiver{
		Header:           Header{CommandLength: 0, CommandID: BIND_RECEIVER, CommandStatus: 0, SequenceNumber: 1},
		SystemID:         "testSystemId",
		Password:         "password",
		SystemType:       "type",
		InterfaceVersion: 0x34,
		AddrTon:          0x00,
		AddrNpi:          0x00,
		AddressRange:     "",
	}
	if got != want {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

func TestReadWriteBindReceiver(t *testing.T) {
	bindReceiver := BindReceiver{
		Header:           Header{CommandLength: 0, CommandID: BIND_RECEIVER, CommandStatus: 0, SequenceNumber: 1},
		SystemID:         "testSystemId",
		Password:         "password",
		SystemType:       "type",
		InterfaceVersion: 0x34,
		AddrTon:          0x00,
		AddrNpi:          0x00,
		AddressRange:     "",
	}
	data := WriteBindReceiver(bindReceiver)
	got, err := ReadBindReceiver(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != bindReceiver {
		t.Fatalf("expected %+v, got %+v", bindReceiver, got)
	}
}

func TestReadBindReceiver_TooShort(t *testing.T) {
	_, err := ReadBindReceiver([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
