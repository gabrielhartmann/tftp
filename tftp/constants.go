package tftp

const dataBlockSize = 512
const timeoutCountMax = 3

const (
	INVALID_LOW_OPCODE = iota
	RRQ
	WRQ
	DATA
	ACK
	ERROR
	INVALID_HIGH_OPCODE
)
