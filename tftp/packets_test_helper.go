package tftp

import (
	"bytes"
	"testing"
)

func opcodeNegativeTestHelper(t *testing.T, bytes []byte) {
	code, err := getOpcode(bytes)

	if err == nil {
		t.Errorf("Expected invalid opcode error from %v, received: %v with error %v", bytes, code, err)
	}
}

func opcodePositiveTestHelper(t *testing.T, bytes []byte, expectedCode uint16) {
	code, err := getOpcode(bytes)

	if err != nil {
		t.Errorf("Expected valid opcode from %v, received: %v with error %v", bytes, code, err)
	}
}

func parseRequestHelperPositive(t *testing.T, input []byte, expectedFile string, expectedMode string) {
	file, mode, err := parseRequest(input)
	if err != nil {
		t.Errorf("parseRequest failed with error: %v", err)
	}

	if file != expectedFile {
		t.Errorf("Expected file: %v, returned: %v", expectedFile, file)
	}

	if mode != expectedMode {
		t.Errorf("Expected mode: %v, returned: %v", expectedMode, mode)
	}

}

func parseRequestHelperNegative(t *testing.T, input []byte) {
	file, mode, err := parseRequest(input)
	if err == nil {
		t.Errorf("Expected parse failure, returned: file:%v mode:%v err:%v", file, mode, err)
	}
}

func parseDataHelperPositive(t *testing.T, input []byte, expectedBlock uint16, expectedData []byte) {
	block, data, err := parseData(input)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if block != expectedBlock {
		t.Errorf("Expected block: %v, returned: %v", expectedBlock, block)
	}

	if !bytes.Equal(data, expectedData) {
		t.Errorf("Expected data: %v, returned: %v", expectedData, data)
	}
}
