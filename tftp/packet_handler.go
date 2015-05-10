package tftp

import (
	"errors"
)

type PacketHandler interface {
	ReadReq(file string, mode string) error
	WriteReq(file string, mode string) error
	Data(block uint16, data []byte) error
	Ack(block uint16) error
	Err(code uint16, msg string) error
}

// This function determines the type of a packet and routes it to the
// appropriate handling method
func HandleTftpPackets(handler PacketHandler, input []byte) error {
	code, err := getOpcode(input)
	if err != nil {
		return err
	}

	switch code {
	case RRQ:
		if file, mode, err := parseRequest(input[2:]); err == nil {
			return handler.ReadReq(file, mode)
		} else {
			return err
		}
	case WRQ:
		if file, mode, err := parseRequest(input[2:]); err == nil {
			return handler.WriteReq(file, mode)
		} else {
			return err
		}
	case DATA:
		if block, data, err := parseData(input[2:]); err == nil {
			return handler.Data(block, data)
		} else {
			return err
		}
	case ACK:
		if block, err := parseAck(input[2:]); err == nil {
			return handler.Ack(block)
		} else {
			return err
		}
	case ERROR:
		if code, msg, err := parseError(input[2:]); err == nil {
			return handler.Err(code, msg)
		} else {
			return err
		}
	default:
		return errors.New("We should never reach the end of HandleTftpPackets")
	}
}
