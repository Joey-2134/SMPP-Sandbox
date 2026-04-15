package pdu

import (
	"testing"
)

func TestReadBindRequest(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"bind_receiver", BIND_RECEIVER},
		{"bind_transmitter", BIND_TRANSMITTER},
		{"bind_transceiver", BIND_TRANSCEIVER},
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
			data = append(data, []byte("testSystemId\x00")...)
			data = append(data, []byte("password\x00")...)
			data = append(data, []byte("type\x00")...)
			data = append(data, 0x34, 0x00, 0x00, 0x00)

			got, err := ReadBindRequest(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := BindRequest{
				Header:           Header{CommandLength: 0, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
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
		})
	}
}

func TestReadWriteBindRequest(t *testing.T) {
	cases := []struct {
		name      string
		commandID uint32
	}{
		{"bind_receiver", BIND_RECEIVER},
		{"bind_transmitter", BIND_TRANSMITTER},
		{"bind_transceiver", BIND_TRANSCEIVER},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bindRequest := BindRequest{
				Header:           Header{CommandLength: 0, CommandID: tc.commandID, CommandStatus: 0, SequenceNumber: 1},
				SystemID:         "testSystemId",
				Password:         "password",
				SystemType:       "type",
				InterfaceVersion: 0x34,
				AddrTon:          0x00,
				AddrNpi:          0x00,
				AddressRange:     "",
			}
			data := WriteBindRequest(bindRequest)
			got, err := ReadBindRequest(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != bindRequest {
				t.Fatalf("expected %+v, got %+v", bindRequest, got)
			}
		})
	}
}

func TestReadBindRequest_TooShort(t *testing.T) {
	_, err := ReadBindRequest([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
