package pdu

import (
	"testing"
)

func TestReadHeader(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x10, // Command Length = 16
		0x00, 0x00, 0x00, 0x01, // Command ID = 1, bind_receiver
		0x00, 0x00, 0x00, 0x00, // Command Status = 0, no error
		0x00, 0x00, 0x00, 0x01, // Sequence Number = 1
	}
	header, err := ReadHeader(data)
	if err != nil {
		t.Fatalf("Error reading header: %v", err)
	}
	if header.CommandLength != 16 {
		t.Fatalf("Expected command length to be 16, got %d", header.CommandLength)
	}
	if header.CommandID != 1 {
		t.Fatalf("Expected command ID to be 1, got %d", header.CommandID)
	}
	if header.CommandStatus != 0 {
		t.Fatalf("Expected command status to be 0, got %d", header.CommandStatus)
	}
	if header.SequenceNumber != 1 {
		t.Fatalf("Expected sequence number to be 1, got %d", header.SequenceNumber)
	}
}

func TestWriteHeader(t *testing.T) {
	header := Header{
		CommandLength:  16,
		CommandID:      1,
		CommandStatus:  0,
		SequenceNumber: 1,
	}
	data := WriteHeader(header)
	expected := []byte{
		0x00, 0x00, 0x00, 0x10,
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01,
	}
	for i, b := range expected {
		if data[i] != b {
			t.Fatalf("byte %d: expected 0x%02X, got 0x%02X", i, b, data[i])
		}
	}
}

func TestReadHeader_TooShort(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x00, // Command Length = 0
	}
	_, err := ReadHeader(data)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestReadWriteHeader(t *testing.T) {
	header := Header{
		CommandLength:  16,
		CommandID:      4, // submit_sm
		CommandStatus:  0,
		SequenceNumber: 1,
	}
	data := WriteHeader(header)
	readHeader, err := ReadHeader(data)
	if err != nil {
		t.Fatalf("Error reading header: %v", err)
	}
	if readHeader != header {
		t.Fatalf("expected %+v, got %+v", header, readHeader)
	}
}
