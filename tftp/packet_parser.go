package tftp

import (
	"errors"
	"fmt"
)

func parseRequest(input []byte) (file string, mode string, err error) {
	file = ""
	mode = ""

	// Find file name
	for i, e := range input {
		if e == 0 {
			file = string(input[:i])

			if len(input) > i+1 {
				input = input[i+1:]
			} else {
				return "", "", errors.New("Input is not long enough include the required 'mode' string")
			}

			break
		}
	}

	// Find mode
	for i, e := range input {
		if e == 0 {
			mode = string(input[:i])
			break
		}
	}

	if err := validateMode(mode); err != nil {
		return "", "", errors.New(fmt.Sprintf("Invalid mode: %v", mode))
	}

	return file, mode, nil
}

func parseData(input []byte) (block uint16, data []byte, err error) {
	block, err = getTwoByteInt(input)
	if err != nil {
		return 0, []byte{}, err
	}

	if len(input) > 2 {
		data = input[2:]
	} else {
		data = []byte{}
	}

	return block, data, err
}

func parseAck(input []byte) (uint16, error) {
	return getTwoByteInt(input)
}

func parseError(input []byte) (uint16, string, error) {
	code, err := getTwoByteInt(input)
	if err != nil {
		return 0, "", err
	}

	if len(input) <= 2 {
		return 0, "", errors.New("Input is too short to extract error message")
	}

	if input[len(input)-1] != 0 {
		return 0, "", errors.New("Error packet must be terminated by a 0")
	}

	if err = validateErrorCode(code); err != nil {
		return 0, "", err
	}

	return code, string(input[2 : len(input)-1]), nil
}
