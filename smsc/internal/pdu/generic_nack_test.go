package pdu

import (
	"testing"
)

func TestReadGenericNack(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"generic_nack", GENERIC_NACK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := WriteHeader(Header{
				CommandLength:  16,
				CommandID:      tc.commandID,
				CommandStatus:  0,
				SequenceNumber: 1,
			})

			got, err := ReadGenericNack(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := GenericNack{
				Header: Header{CommandLength: 16, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
			}
			if got != want {
				t.Fatalf("expected %+v, got %+v", want, got)
			}
		})
	}
}

func TestReadWriteGenericNack(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"generic_nack", GENERIC_NACK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			genericNack := GenericNack{
				Header: Header{CommandLength: 16, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
			}
			data := WriteGenericNack(genericNack)
			got, err := ReadGenericNack(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != genericNack {
				t.Fatalf("expected %+v, got %+v", genericNack, got)
			}
		})
	}
}

func TestReadGenericNack_TooShort(t *testing.T) {
	_, err := ReadGenericNack([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
