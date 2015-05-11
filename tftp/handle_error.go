package tftp

import (
	"errors"
	"net"

	"github.com/Sirupsen/logrus"
)

const (
	UndefinedError = iota
	FileNotFound
	AccessViolation
	DiskFull
	IllegalOperation
	UnknownTid
	FileExists
	NoSuchUser
)

func getErrorPacket(code uint16, msg string) *ErrorPacket {
	codeBytes, _ := convertIntToBytes(code)
	errorPacket := NewErrorPacket([2]byte{codeBytes[0], codeBytes[1]}, msg)
	return errorPacket
}

func HandleError(writer *TftpReaderWriter, code uint16, msg string) error {
	errorPacket := getErrorPacket(code, msg)

	logrus.Infof("Sending error packet: code: %v, msg: %v", errorPacket.code, errorPacket.msg)
	writer.Write(errorPacket.bytes)
	return errors.New(msg)
}

func isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}
