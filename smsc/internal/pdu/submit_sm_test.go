package pdu

import (
	"testing"
)

func TestReadSubmitSM(t *testing.T) {
	data := WriteHeader(Header{
		CommandLength:  0,
		CommandID:      SUBMIT_SM,
		CommandStatus:  0,
		SequenceNumber: 1,
	})
	data = append(data, []byte("\x00")...)      // service_type (empty)
	data = append(data, 0x01, 0x01)             // source_addr_ton, source_addr_npi
	data = append(data, []byte("12345\x00")...) // source_addr
	data = append(data, 0x01, 0x01)             // dest_addr_ton, dest_addr_npi
	data = append(data, []byte("67890\x00")...) // dest_addr
	data = append(data, 0x00, 0x00, 0x00)       // esm_class, protocol_id, priority_flag
	data = append(data, []byte("\x00")...)      // schedule_delivery_time (empty)
	data = append(data, []byte("\x00")...)      // validity_period (empty)
	data = append(data, 0x01, 0x00, 0x00, 0x00) // registered_delivery, replace_if_present, data_coding, sm_default_msg_id
	data = append(data, 0x05)                   // sm_length
	data = append(data, []byte("Hello")...)     // message

	got, err := ReadSubmitSM(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := SubmitSM{
		Header:               Header{CommandLength: 0, CommandID: SUBMIT_SM, CommandStatus: 0, SequenceNumber: 1},
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
	if got.Header != want.Header {
		t.Fatalf("header mismatch: expected %+v, got %+v", want.Header, got.Header)
	}
	if got.SourceAddr != want.SourceAddr {
		t.Fatalf("expected SourceAddr %s, got %s", want.SourceAddr, got.SourceAddr)
	}
	if got.DestAddr != want.DestAddr {
		t.Fatalf("expected DestAddr %s, got %s", want.DestAddr, got.DestAddr)
	}
	if got.SMLength != want.SMLength {
		t.Fatalf("expected SMLength %d, got %d", want.SMLength, got.SMLength)
	}
	if string(got.Message) != string(want.Message) {
		t.Fatalf("expected Message %s, got %s", want.Message, got.Message)
	}
}

func TestReadWriteSubmitSM(t *testing.T) {
	submitSM := SubmitSM{
		Header:               Header{CommandLength: 0, CommandID: SUBMIT_SM, CommandStatus: 0, SequenceNumber: 1},
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

	data := WriteSubmitSM(submitSM)
	got, err := ReadSubmitSM(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Header != submitSM.Header {
		t.Fatalf("header mismatch: expected %+v, got %+v", submitSM.Header, got.Header)
	}
	if got.SourceAddr != submitSM.SourceAddr {
		t.Fatalf("expected SourceAddr %s, got %s", submitSM.SourceAddr, got.SourceAddr)
	}
	if got.DestAddr != submitSM.DestAddr {
		t.Fatalf("expected DestAddr %s, got %s", submitSM.DestAddr, got.DestAddr)
	}
	if got.SMLength != submitSM.SMLength {
		t.Fatalf("expected SMLength %d, got %d", submitSM.SMLength, got.SMLength)
	}
	if string(got.Message) != string(submitSM.Message) {
		t.Fatalf("expected Message %s, got %s", submitSM.Message, got.Message)
	}
}

func TestReadSubmitSM_TooShort(t *testing.T) {
	_, err := ReadSubmitSM([]byte{0x00, 0x01, 0x02})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
