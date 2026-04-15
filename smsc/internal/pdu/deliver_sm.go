package pdu

import "errors"

type DeliverSM struct {
	Header               Header
	ServiceType          string //max 6
	SourceAddrTon        uint8
	SourceAddrNpi        uint8
	SourceAddr           string // max 21
	DestAddrTon          uint8
	DestAddrNpi          uint8
	DestAddr             string // max 21
	ESMClass             uint8
	ProtocolID           uint8
	PriorityFlag         uint8
	ScheduleDeliveryTime string // must be null for deliver_sm
	ValidityPeriod       string // must be null for deliver_sm
	RegisteredDelivery   uint8
	ReplaceIfPresentFlag uint8 // must be null for deliver_sm
	DataCoding           uint8
	SmDefaultMsgID       uint8 // must be null for deliver_sm
	SMLength             uint8
	Message              []byte // max 254
}

func ReadDeliverSM(data []byte) (DeliverSM, error) {
	if len(data) < 16 {
		return DeliverSM{}, errors.New("data too short to contain a deliver_sm PDU")
	}
	header, err := ReadHeader(data[0:16])
	if err != nil {
		return DeliverSM{}, err
	}

	offset := 16

	serviceType, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return DeliverSM{}, err
	}
	offset += n

	sourceAddrTon := data[offset]
	sourceAddrNpi := data[offset+1]
	offset += 2

	sourceAddr, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return DeliverSM{}, err
	}
	offset += n

	destAddrTon := data[offset]
	destAddrNpi := data[offset+1]
	offset += 2

	destAddr, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return DeliverSM{}, err
	}
	offset += n

	esmClass := data[offset]
	protocolID := data[offset+1]
	priorityFlag := data[offset+2]
	offset += 3

	scheduleDeliveryTime, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return DeliverSM{}, err
	}
	offset += n

	validityPeriod, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return DeliverSM{}, err
	}
	offset += n

	registeredDelivery := data[offset]
	replaceIfPresentFlag := data[offset+1]
	dataCoding := data[offset+2]
	smDefaultMsgID := data[offset+3]
	smLength := data[offset+4]
	offset += 5

	if len(data[offset:]) < int(smLength) {
		return DeliverSM{}, errors.New("data too short to contain the message")
	}
	message := data[offset : offset+int(smLength)]

	return DeliverSM{
		Header:               header,
		ServiceType:          serviceType,
		SourceAddrTon:        sourceAddrTon,
		SourceAddrNpi:        sourceAddrNpi,
		SourceAddr:           sourceAddr,
		DestAddrTon:          destAddrTon,
		DestAddrNpi:          destAddrNpi,
		DestAddr:             destAddr,
		ESMClass:             esmClass,
		ProtocolID:           protocolID,
		PriorityFlag:         priorityFlag,
		ScheduleDeliveryTime: scheduleDeliveryTime,
		ValidityPeriod:       validityPeriod,
		RegisteredDelivery:   registeredDelivery,
		ReplaceIfPresentFlag: replaceIfPresentFlag,
		DataCoding:           dataCoding,
		SmDefaultMsgID:       smDefaultMsgID,
		SMLength:             smLength,
		Message:              message,
	}, nil
}

func WriteDeliverSM(s DeliverSM) []byte {
	data := WriteHeader(s.Header)
	data = append(data, []byte(s.ServiceType)...)
	data = append(data, 0x00)
	data = append(data, s.SourceAddrTon, s.SourceAddrNpi)
	data = append(data, []byte(s.SourceAddr)...)
	data = append(data, 0x00)
	data = append(data, s.DestAddrTon, s.DestAddrNpi)
	data = append(data, []byte(s.DestAddr)...)
	data = append(data, 0x00)
	data = append(data, s.ESMClass, s.ProtocolID, s.PriorityFlag)
	data = append(data, []byte(s.ScheduleDeliveryTime)...)
	data = append(data, 0x00)
	data = append(data, []byte(s.ValidityPeriod)...)
	data = append(data, 0x00)
	data = append(data, s.RegisteredDelivery, s.ReplaceIfPresentFlag, s.DataCoding, s.SmDefaultMsgID)
	data = append(data, s.SMLength)
	data = append(data, s.Message...)
	return data
}
