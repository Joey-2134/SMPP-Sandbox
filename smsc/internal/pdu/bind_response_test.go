package pdu

import (
	"testing"
)

func TestReadBindResponse(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"bind_receiver_resp", BIND_RECEIVER_RESP},
		{"bind_transmitter_resp", BIND_TRANSMITTER_RESP},
		{"bind_transceiver_resp", BIND_TRANSCEIVER_RESP},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte{}
			data = append(data, WriteHeader(Header{
				CommandLength:  0,
				CommandID:      tc.commandID,
				CommandStatus:  0,
				SequenceNumber: 1,
			})...)
			data = append(data, []byte("testSmsc\x00")...)

			got, err := ReadBindResponse(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := BindResponse{
				Header:   Header{CommandLength: 0, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
				SystemID: "testSmsc",
			}
			if got != want {
				t.Fatalf("expected %+v, got %+v", want, got)
			}
		})
	}
}

func TestReadWriteBindResponse(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"bind_receiver_resp", BIND_RECEIVER_RESP},
		{"bind_transmitter_resp", BIND_TRANSMITTER_RESP},
		{"bind_transceiver_resp", BIND_TRANSCEIVER_RESP},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bindResponse := BindResponse{
				Header:   Header{CommandLength: 0, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
				SystemID: "testSmsc",
			}
			data := WriteBindResponse(bindResponse)
			got, err := ReadBindResponse(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != bindResponse {
				t.Fatalf("expected %+v, got %+v", bindResponse, got)
			}
		})
	}
}

func TestReadBindResponse_TooShort(t *testing.T) {
	_, err := ReadBindResponse([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
