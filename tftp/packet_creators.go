package tftp

type DataPacket struct {
	block uint16
	data  []byte
	bytes []byte
}

func NewDataPacket(block [2]byte, data []byte) *DataPacket {
	block64, _ := convertBytesToInt([]byte{block[0], block[1]})
	block16 := uint16(block64)
	bytes := append([]byte{0x0, 0x3, block[0], block[1]}, data...)

	return &DataPacket{
		block: block16,
		data:  data,
		bytes: bytes,
	}
}

type AckPacket struct {
	block uint16
	bytes []byte
}

func NewAckPacket(block [2]byte) *AckPacket {
	block64, _ := convertBytesToInt([]byte{block[0], block[1]})
	block16 := uint16(block64)
	bytes := []byte{0x0, 0x4, block[0], block[1]}

	return &AckPacket{
		block: block16,
		bytes: bytes,
	}
}

type ErrorPacket struct {
	code  uint16
	msg   string
	bytes []byte
}

func NewErrorPacket(code [2]byte, msg string) *ErrorPacket {
	code64, _ := convertBytesToInt([]byte{code[0], code[1]})
	code16 := uint16(code64)
	bytes := []byte{0x0, 0x5, code[0], code[1]}
	bytes = append(bytes, []byte(msg)...)
	bytes = append(bytes, 0)

	return &ErrorPacket{
		code:  code16,
		msg:   msg,
		bytes: bytes,
	}
}
