package tftp

import (
	"errors"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	. "github.com/gabrielhartmann/tftp/fileserv"
)

type WriteSession struct {
	writer       *TftpReaderWriter
	fileServ     FileServer
	block        uint16
	fileName     string
	dataBuffer   []byte
	fileComplete bool
	timeoutCount int
}

var logPrefix string

func StartNewWriteSession(remoteAddr *net.UDPAddr, file string, fileServ FileServer) error {
	// Create TftpReaderWriter
	writer, err := NewTftpReaderWriter(remoteAddr, true)
	if err != nil {
		return err
	}

	if fileServ.FileExists(file) {
		return HandleError(writer, FileExists, fmt.Sprintf("File '%v' already exists", file))
	}

	writeSession := &WriteSession{
		writer:       writer,
		fileServ:     fileServ,
		block:        0,
		fileName:     file,
		dataBuffer:   []byte{},
		fileComplete: false,
		timeoutCount: 0,
	}

	logrus.Infof("[Write Session %v]: Start for file '%v'", remoteAddr.Port, file)

	// Main work floop with bounded timeouts
	for writeSession.timeoutCount < timeoutCountMax {
		if err = writeSession.Start(); err != nil {
			if isTimeout(err) {
				logrus.Infof("[Write Session %v]: timeout %d", remoteAddr.Port, writeSession.timeoutCount)
				writeSession.timeoutCount++
			} else {
				return err
			}
		} else {
			return err
		}
	}

	return err
}

func (s *WriteSession) Start() error {
	// Write initial ack packet
	if err := s.writeAck(); err != nil {
		return err
	}

	// Read packets continuously until the file is complete
	// The file is complete when a data packet is received
	// with fewer than 512 bytes.  See the Data() method below
	for {
		if s.fileComplete {
			logrus.Infof("[Write Session]: completed file: '%v'", s.fileName)
			return nil
		}

		if bytes, _, err := s.writer.Read(); err != nil {
			return err
		} else {
			if err := s.handleTftpPackets(bytes); err != nil {
				return err
			}
		}
	}
}

func isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}

// This function determines the type of a packet and routes it to the
// appropriate handling method
func (s *WriteSession) handleTftpPackets(input []byte) error {
	code, err := getOpcode(input)
	if err != nil {
		return err
	}

	switch code {
	case RRQ:
		if file, mode, err := parseRequest(input[2:]); err == nil {
			return s.ReadReq(file, mode)
		} else {
			return err
		}
	case WRQ:
		if file, mode, err := parseRequest(input[2:]); err == nil {
			return s.WriteReq(file, mode)
		} else {
			return err
		}
	case DATA:
		if block, data, err := parseData(input[2:]); err == nil {
			return s.Data(block, data)
		} else {
			return err
		}
	case ACK:
		if block, err := parseAck(input[2:]); err == nil {
			return s.Ack(block)
		} else {
			return err
		}
	case ERROR:
		if code, msg, err := parseError(input[2:]); err == nil {
			return s.Err(code, msg)
		} else {
			return err
		}
	}

	return errors.New("We should never reach the end of ParseTftpPackets")
}

// Generate the next ACK packet
func (s *WriteSession) getAckPacket() (*AckPacket, error) {
	if bytes, err := convertIntToBytes(s.block); err != nil {
		return nil, err
	} else {
		blockArr := [2]byte{bytes[0], bytes[1]}
		return NewAckPacket(blockArr), nil
	}
}

// Write the next ACK packet
func (s *WriteSession) writeAck() error {
	if ack, err := s.getAckPacket(); err != nil {
		return err
	} else {
		_, err = s.writer.Write(ack.bytes)
		return err
	}
}

func (s *WriteSession) ReadReq(file string, mode string) error {
	return errors.New("ReadReq operations are not supported on read handlers")
}

func (s *WriteSession) WriteReq(file string, mode string) error {
	return errors.New("WriteReq operations are not supported on read handlers")
}

func (s *WriteSession) Data(block uint16, data []byte) error {
	if block == s.block+1 {
		s.dataBuffer = append(s.dataBuffer, data...)
		s.block++
	} else {
		return errors.New(fmt.Sprintf("Expected block %v, received %v", s.block+1, block))
	}

	s.writeAck()

	if len(data) < dataBlockSize {
		s.fileComplete = true
		file := File{
			Name: s.fileName,
			Data: s.dataBuffer,
		}

		if err := s.fileServ.Write(&file); err != nil {
			return err
		}

		logrus.Infof("[Write Session]: Wrote %v to file server %v bytes", file.Name, len(file.Data))
	}

	return nil
}

func (s *WriteSession) Ack(block uint16) error {
	return errors.New("Ack operations are not supported on this handlers")
	return nil
}

func (s *WriteSession) Err(code uint16, msg string) error {
	logrus.Infof("Received Error with code %v and message %v", code, msg)
	return errors.New(fmt.Sprintf("Received Error with code %v and message %v", code, msg))
}
