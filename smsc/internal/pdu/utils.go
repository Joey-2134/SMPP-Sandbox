package pdu

import "errors"

func ReadCOctetString(data []byte) (string, int, error) {
	for i, b := range data {
		if b == 0x00 {
			return string(data[:i]), i + 1, nil
		}
	}
	return "", 0, errors.New("no null terminator found")
}
