package pdu

import (
	"testing"
)

func TestReadDeliverSM(t *testing.T) {
	data := WriteHeader(Header{
		CommandLength:  0,
		CommandID:      DELIVER_SM,
		CommandStatus:  0,
		SequenceNumber: 1,
	})
	data = append(data, []byte("\x00")...)      // service_type (empty)
	data = append(data, 0x01, 0x01)             // source_addr_ton, source_addr_npi
	data = append(data, []byte("12345\x00")...) // source_addr
	data = append(data, 0x01, 0x01)             // dest_addr_ton, dest_addr_npi
	data = append(data, []byte("67890\x00")...) // dest_addr
	data = append(data, 0x00, 0x00, 0x00)       // esm_class, protocol_id, priority_flag
	data = append(data, []byte("\x00")...)      // schedule_delivery_time (must be null)
	data = append(data, []byte("\x00")...)      // validity_period (must be null)
	data = append(data, 0x01, 0x00, 0x00, 0x00) // registered_delivery, replace_if_present, data_coding, sm_default_msg_id
	data = append(data, 0x05)                   // sm_length
	data = append(data, []byte("Hello")...)     // message

	got, err := ReadDeliverSM(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Header.CommandID != DELIVER_SM {
		t.Fatalf("expected CommandID DELIVER_SM, got 0x%08X", got.Header.CommandID)
	}
	if got.SourceAddr != "12345" {
		t.Fatalf("expected SourceAddr 12345, got %s", got.SourceAddr)
	}
	if got.DestAddr != "67890" {
		t.Fatalf("expected DestAddr 67890, got %s", got.DestAddr)
	}
	if got.SMLength != 5 {
		t.Fatalf("expected SMLength 5, got %d", got.SMLength)
	}
	if string(got.Message) != "Hello" {
		t.Fatalf("expected Message Hello, got %s", got.Message)
	}
}

func TestReadWriteDeliverSM(t *testing.T) {
	deliverSM := DeliverSM{
		Header:               Header{CommandLength: 0, CommandID: DELIVER_SM, CommandStatus: 0, SequenceNumber: 1},
		ServiceType:          "",
		SourceAddrTon:        0x01,
		SourceAddrNpi:        0x01,
		SourceAddr:           "12345",
		DestAddrTon:          0x01,
		DestAddrNpi:          0x01,
		DestAddr:             "67890",
		ESMClass:             0x00,
		ProtocolID:           0x00,
		PriorityFlag:         0x00,
		ScheduleDeliveryTime: "",
		ValidityPeriod:       "",
		RegisteredDelivery:   0x01,
		ReplaceIfPresentFlag: 0x00,
		DataCoding:           0x00,
		SmDefaultMsgID:       0x00,
		SMLength:             5,
		Message:              []byte("Hello"),
	}

	data := WriteDeliverSM(deliverSM)
	got, err := ReadDeliverSM(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Header != deliverSM.Header {
		t.Fatalf("header mismatch: expected %+v, got %+v", deliverSM.Header, got.Header)
	}
	if got.SourceAddr != deliverSM.SourceAddr {
		t.Fatalf("expected SourceAddr %s, got %s", deliverSM.SourceAddr, got.SourceAddr)
	}
	if got.DestAddr != deliverSM.DestAddr {
		t.Fatalf("expected DestAddr %s, got %s", deliverSM.DestAddr, got.DestAddr)
	}
	if got.SMLength != deliverSM.SMLength {
		t.Fatalf("expected SMLength %d, got %d", deliverSM.SMLength, got.SMLength)
	}
	if string(got.Message) != string(deliverSM.Message) {
		t.Fatalf("expected Message %s, got %s", deliverSM.Message, got.Message)
	}
}

func TestReadDeliverSM_TooShort(t *testing.T) {
	_, err := ReadDeliverSM([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
