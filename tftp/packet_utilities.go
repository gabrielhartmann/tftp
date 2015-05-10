package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

func validateMode(mode string) error {
	if !strings.EqualFold(mode, "octet") {
		return errors.New("Only 'octet' mode is supported")
	}

	return nil
}

func validateErrorCode(code uint16) error {
	if code > 7 {
		return errors.New(fmt.Sprintf("Invalid error code: %v", code))
	}

	return nil
}

func getOpcode(input []byte) (uint16, error) {
	opCode, err := getTwoByteInt(input)
	if err != nil {
		return 0, err
	}

	if opCode <= INVALID_LOW_OPCODE || opCode >= INVALID_HIGH_OPCODE {
		return 1, errors.New(fmt.Sprintf("Invalid opCode '%v' out of range", opCode))
	}

	return opCode, nil
}

func getTwoByteInt(input []byte) (uint16, error) {
	if len(input) < 2 {
		return 0, errors.New("Require at least 2 bytes to parse uint16")
	}

	output, err := convertBytesToInt(input[:2])
	if err != nil {
		return 0, err
	}

	return uint16(output), nil
}

func convertBytesToInt(input []byte) (uint64, error) {
	// Pad a slice as golang doesn't provide a straight conversion to 16-bit integers
	pad := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	paddedInput := append(pad, input...)

	var result uint64
	err := binary.Read(bytes.NewBuffer(paddedInput[:]), binary.BigEndian, &result)

	if err != nil {
		return 0, errors.New(fmt.Sprintf("Failed to convert bytes '%v' to integer '%v'", paddedInput, result))
	}

	return result, nil
}

func convertIntToBytes(input uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, input)

	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Failed to convert uint16 '%v' to bytes", input))
	}

	return buf.Bytes(), nil
}
