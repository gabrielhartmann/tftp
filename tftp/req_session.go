package tftp

import (
	"errors"
	"net"

	"github.com/Sirupsen/logrus"
	. "github.com/gabrielhartmann/tftp/fileserv"
)

type ReqSession struct {
	rw       *TftpReaderWriter
	fileServ FileServer
}

func StartNewReqSession() error {
	rw, err := NewTftpReaderWriter(nil, false)
	if err != nil {
		logrus.Errorf("Failed to start request session with err: %v", err)
		return err
	}

	reqSession := NewReqSession(rw)
	logrus.Infof("[Request Session]: Starting")
	return reqSession.Start()
}

func NewReqSession(rw *TftpReaderWriter) *ReqSession {
	return &ReqSession{
		rw:       rw,
		fileServ: NewMemFileServer(),
	}
}

func (s *ReqSession) Start() error {
	// Main work loop that reads read/write requests
	// and spawns read or write sessions as apporpriate to handle them
	for {
		if bytes, addr, err := s.rw.Read(); err != nil {
			return err
		} else {
			if err := s.handleTftpPackets(addr, bytes); err != nil {
				return err
			}
		}
	}
}

// Routes tftp packets to the appropriate handlers below
func (s *ReqSession) handleTftpPackets(addr *net.UDPAddr, input []byte) error {
	code, err := getOpcode(input)
	if err != nil {
		return err
	}

	switch code {
	case RRQ:
		if file, mode, err := parseRequest(input[2:]); err == nil {
			return s.ReadReq(addr, file, mode)
		} else {
			return err
		}
	case WRQ:
		if file, mode, err := parseRequest(input[2:]); err == nil {
			return s.WriteReq(addr, file, mode)
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

func (s *ReqSession) ReadReq(addr *net.UDPAddr, file string, mode string) error {
	logrus.Infof("[Request Session]: Received ReadReq for file: %v, in mode %v", file, mode)
	go StartNewReadSession(addr, file, s.fileServ)
	return nil
}

func (s *ReqSession) WriteReq(addr *net.UDPAddr, file string, mode string) error {
	logrus.Infof("[Request Session]: Received WriteReq for file: %v, in mode %v", file, mode)
	go StartNewWriteSession(addr, file, s.fileServ)
	return nil
}

func (s *ReqSession) Data(block uint16, data []byte) error {
	logrus.Infof("[Request Session]: Data operations are not supported in req session")
	return errors.New("Data operations are not supported on this handler")
}

func (s *ReqSession) Ack(block uint16) error {
	logrus.Infof("[Request Session]: Ack operations are not supported in req session")
	return errors.New("Ack operations are not supported on this handler")
}

func (s *ReqSession) Err(code uint16, msg string) error {
	logrus.Infof("[Request Session]: Error operations are not supported in req session")
	return errors.New("Error operations are not supported on this handler")
}
