package pdu

import (
	"testing"
)

func TestReadUnbind(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"unbind", UNBIND},
		{"unbind_resp", UNBIND_RESP},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := WriteHeader(Header{
				CommandLength:  16,
				CommandID:      tc.commandID,
				CommandStatus:  0,
				SequenceNumber: 1,
			})

			got, err := ReadUnbind(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := Unbind{
				Header: Header{CommandLength: 16, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
			}
			if got != want {
				t.Fatalf("expected %+v, got %+v", want, got)
			}
		})
	}
}

func TestReadWriteUnbind(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"unbind", UNBIND},
		{"unbind_resp", UNBIND_RESP},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			unbind := Unbind{
				Header: Header{CommandLength: 16, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
			}
			data := WriteUnbind(unbind)
			got, err := ReadUnbind(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != unbind {
				t.Fatalf("expected %+v, got %+v", unbind, got)
			}
		})
	}
}

func TestReadUnbind_TooShort(t *testing.T) {
	_, err := ReadUnbind([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
