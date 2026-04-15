package pdu

import "errors"

// see SMPP 3.4 spec section 4.1.3
type BindReceiver struct {
	Header           Header
	SystemID         string
	Password         string
	SystemType       string
	InterfaceVersion uint8
	AddrTon          uint8
	AddrNpi          uint8
	AddressRange     string
}

func ReadBindReceiver(data []byte) (BindReceiver, error) {
	if len(data) < 16 {
		return BindReceiver{}, errors.New("data too short to contain a bind_receiver PDU")
	}
	header, err := ReadHeader(data[0:16])
	if err != nil {
		return BindReceiver{}, err
	}

	offset := 16

	systemID, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return BindReceiver{}, err
	}
	offset += n

	password, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return BindReceiver{}, err
	}
	offset += n

	systemType, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return BindReceiver{}, err
	}
	offset += n

	interfaceVersion := data[offset]
	addrTon := data[offset+1]
	addrNpi := data[offset+2]
	offset += 3

	addressRange, _, err := ReadCOctetString(data[offset:])
	if err != nil {
		return BindReceiver{}, err
	}

	return BindReceiver{
		Header:           header,
		SystemID:         systemID,
		Password:         password,
		SystemType:       systemType,
		InterfaceVersion: interfaceVersion,
		AddrTon:          addrTon,
		AddrNpi:          addrNpi,
		AddressRange:     addressRange,
	}, nil
}

func WriteBindReceiver(bindReceiver BindReceiver) []byte {
	data := WriteHeader(bindReceiver.Header)
	data = append(data, []byte(bindReceiver.SystemID)...)
	data = append(data, 0x00)
	data = append(data, []byte(bindReceiver.Password)...)
	data = append(data, 0x00)
	data = append(data, []byte(bindReceiver.SystemType)...)
	data = append(data, 0x00)
	data = append(data, bindReceiver.InterfaceVersion)
	data = append(data, bindReceiver.AddrTon)
	data = append(data, bindReceiver.AddrNpi)
	data = append(data, []byte(bindReceiver.AddressRange)...)
	data = append(data, 0x00)
	return data
}
