package pdu

import "errors"

// see SMPP 3.4 spec section 4.1
type BindRequest struct {
	Header           Header
	SystemID         string
	Password         string
	SystemType       string
	InterfaceVersion uint8
	AddrTon          uint8
	AddrNpi          uint8
	AddressRange     string
}

func ReadBindRequest(data []byte) (BindRequest, error) {
	if len(data) < 16 {
		return BindRequest{}, errors.New("data too short to contain a bind_receiver PDU")
	}
	header, err := ReadHeader(data[0:16])
	if err != nil {
		return BindRequest{}, err
	}

	offset := 16

	systemID, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return BindRequest{}, err
	}
	offset += n

	password, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return BindRequest{}, err
	}
	offset += n

	systemType, n, err := ReadCOctetString(data[offset:])
	if err != nil {
		return BindRequest{}, err
	}
	offset += n

	interfaceVersion := data[offset]
	addrTon := data[offset+1]
	addrNpi := data[offset+2]
	offset += 3

	addressRange, _, err := ReadCOctetString(data[offset:])
	if err != nil {
		return BindRequest{}, err
	}

	return BindRequest{
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

func WriteBindRequest(bindRequest BindRequest) []byte {
	data := WriteHeader(bindRequest.Header)
	data = append(data, []byte(bindRequest.SystemID)...)
	data = append(data, 0x00)
	data = append(data, []byte(bindRequest.Password)...)
	data = append(data, 0x00)
	data = append(data, []byte(bindRequest.SystemType)...)
	data = append(data, 0x00)
	data = append(data, bindRequest.InterfaceVersion)
	data = append(data, bindRequest.AddrTon)
	data = append(data, bindRequest.AddrNpi)
	data = append(data, []byte(bindRequest.AddressRange)...)
	data = append(data, 0x00)
	return data
}
