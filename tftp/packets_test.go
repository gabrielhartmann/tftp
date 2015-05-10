package tftp

import (
	"bytes"
	"testing"
)

func TestGetOpCodeNegative(t *testing.T) {
	opcodeNegativeTestHelper(t, []byte{1})
	opcodeNegativeTestHelper(t, []byte{0, 0})
	opcodeNegativeTestHelper(t, []byte{0, 0, 1})
	opcodeNegativeTestHelper(t, []byte{0, 6})
}

func TestGetOpCodePositive(t *testing.T) {
	opcodePositiveTestHelper(t, []byte{0, 1}, 1)
	opcodePositiveTestHelper(t, []byte{0, 2}, 2)
	opcodePositiveTestHelper(t, []byte{0, 3}, 3)
	opcodePositiveTestHelper(t, []byte{0, 4}, 4)
	opcodePositiveTestHelper(t, []byte{0, 5}, 5)
	opcodePositiveTestHelper(t, []byte{0, 1, 5}, 1)
	opcodePositiveTestHelper(t, []byte{0, 2, 5}, 2)
	opcodePositiveTestHelper(t, []byte{0, 3, 5}, 3)
	opcodePositiveTestHelper(t, []byte{0, 4, 5}, 4)
	opcodePositiveTestHelper(t, []byte{0, 5, 5}, 5)

}

func TestParseRequestPositive(t *testing.T) {
	expectedFile := "foo"
	expectedMode := "octet"
	input := []byte{'f', 'o', 'o', 0, 'o', 'c', 't', 'e', 't', 0}
	parseRequestHelperPositive(t, input, expectedFile, expectedMode)
}

func TestParseRequestNegative(t *testing.T) {
	parseRequestHelperNegative(t, []byte{'f', 'o', 'o', 0})
	parseRequestHelperNegative(t, []byte{'f', 'o', 'o', 0, 'o', 'c', 't', 'e', 't'})
	parseRequestHelperNegative(t, []byte{0, 'f', 'o', 'o', 0, 'o', 'c', 't', 'e', 't', 0})
	parseRequestHelperNegative(t, []byte{})
	parseRequestHelperNegative(t, []byte{0, 0})
}

func TestParseDataPositive(t *testing.T) {
	var expectedBlock uint16 = 4
	expectedData := []byte{'d', 'a', 't', 'a'}
	parseDataHelperPositive(t, []byte{0, 4, 'd', 'a', 't', 'a'}, expectedBlock, expectedData)

	expectedBlock = 9
	expectedData = []byte{}
	parseDataHelperPositive(t, []byte{0, 9}, expectedBlock, expectedData)

}

func TestParseDataNegative(t *testing.T) {
	_, _, err := parseData([]byte{4})

	if err == nil {
		t.Errorf("Expected parseData to fail, but err: %v", err)
	}
}

func TestParseAck(t *testing.T) {
	var expectedBlock uint16 = 7
	block, err := parseAck([]byte{0, 7})

	if err != nil {
		t.Errorf("Expected successful parseAck, returned err: %v", err)
	}

	if block != expectedBlock {
		t.Errorf("Expected block: %v, returned %v", expectedBlock, block)
	}
}

func TestParseError(t *testing.T) {
	var expectedErrorCode uint16 = 2
	expectedMessage := "msg"

	code, msg, err := parseError([]byte{0, 2, 'm', 's', 'g', 0})

	if err != nil {
		t.Errorf("Expected parseError to succeed, returned error: %v", err)
	}

	if code != expectedErrorCode {
		t.Errorf("Expected error code: %v, returned %v", expectedErrorCode, code)
	}

	if msg != "msg" {
		t.Errorf("Expected error message: %v, returned %v", expectedMessage, msg)
	}

}

func TestCreateDataPacket(t *testing.T) {
	data := []byte{'a', 'b', 'c'}
	expectedBytes := []byte{0, 3, 0, 1, 'a', 'b', 'c'}
	dataPacket := NewDataPacket([2]byte{0, 1}, data)

	if dataPacket.block != 1 {
		t.Errorf("Expected block value of 1, received: %v", dataPacket.block)
	}

	if !bytes.Equal(dataPacket.data, data) {
		t.Errorf("Expected data: %v, received: %v", dataPacket.data, data)
	}

	if !bytes.Equal(dataPacket.bytes, expectedBytes) {
		t.Errorf("Expected bytes: %v, received: %v", dataPacket.bytes, expectedBytes)
	}
}

func TestCreateAckPacket(t *testing.T) {
	expectedBytes := []byte{0, 4, 0, 7}
	ackPacket := NewAckPacket([2]byte{0, 7})

	if ackPacket.block != 7 {
		t.Errorf("Expected block value of 7, received: %v", ackPacket.block)
	}

	if !bytes.Equal(ackPacket.bytes, expectedBytes) {
		t.Errorf("Expected bytes: %v, received: %v", expectedBytes, ackPacket.bytes)
	}
}

func TestCreateErrorPacket(t *testing.T) {
	expectedBytes := []byte{0, 5, 0, 2, 'e', 'r', 'r', 0}
	msg := "err"
	errPacket := NewErrorPacket([2]byte{0, 2}, msg)

	if errPacket.msg != msg {
		t.Errorf("Expected msg: %v, received: %v", msg, errPacket.msg)
	}

	if errPacket.code != 2 {
		t.Errorf("Expected error code: 2, received: %v", errPacket.code)
	}

	if !bytes.Equal(errPacket.bytes, expectedBytes) {
		t.Errorf("Expected bytes: %v, received: %v", expectedBytes, errPacket.bytes)
	}
}
