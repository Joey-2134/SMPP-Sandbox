package pdu

import (
	"testing"
)

func TestReadEnquireLink(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"enquire_link", ENQUIRE_LINK},
		{"enquire_link_resp", ENQUIRE_LINK_RESP},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := WriteHeader(Header{
				CommandLength:  16,
				CommandID:      tc.commandID,
				CommandStatus:  0,
				SequenceNumber: 1,
			})

			got, err := ReadEnquireLink(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := EnquireLink{
				Header: Header{CommandLength: 16, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
			}
			if got != want {
				t.Fatalf("expected %+v, got %+v", want, got)
			}
		})
	}
}

func TestReadWriteEnquireLink(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"enquire_link", ENQUIRE_LINK},
		{"enquire_link_resp", ENQUIRE_LINK_RESP},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enquireLink := EnquireLink{
				Header: Header{CommandLength: 16, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
			}
			data := WriteEnquireLink(enquireLink)
			got, err := ReadEnquireLink(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != enquireLink {
				t.Fatalf("expected %+v, got %+v", enquireLink, got)
			}
		})
	}
}

func TestReadEnquireLink_TooShort(t *testing.T) {
	_, err := ReadEnquireLink([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
