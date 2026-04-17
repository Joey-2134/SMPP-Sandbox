package pdu

import "errors"

const (
	NULL_TERMINATOR = 0x00
)

func ReadCOctetString(data []byte) (string, int, error) {
	for i, b := range data {
		if b == NULL_TERMINATOR {
			return string(data[:i]), i + 1, nil
		}
	}
	return "", 0, errors.New("no null terminator found")
}
