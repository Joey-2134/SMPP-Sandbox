package pdu

import (
	"testing"
)

func TestReadCOctetString(t *testing.T) {
	data := []byte("test\x00")
	string, length, err := ReadCOctetString(data)
	if err != nil {
		t.Fatalf("Error reading COctetString: %v", err)
	}
	if string != "test" {
		t.Fatalf("Expected string to be test, got %s", string)
	}
	if length != 5 {
		t.Fatalf("Expected length to be 5, got %d", length)
	}
}
