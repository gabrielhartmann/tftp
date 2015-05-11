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
			if err := HandleTftpPackets(s, addr, bytes); err != nil {
				return err
			}
		}
	}
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
