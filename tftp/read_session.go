package tftp

import (
	"errors"
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	. "github.com/gabrielhartmann/tftp/fileserv"
)

type ReadSession struct {
	rw           *TftpReaderWriter
	file         *File
	currBlock    uint16
	lastBlock    int
	fileComplete bool
	timeoutCount int
}

func StartNewReadSession(remoteAddr *net.UDPAddr, fileName string, fileServ FileServer) error {
	// Create TftpReaderWriter
	rw, err := NewTftpReaderWriter(remoteAddr, true)
	if err != nil {
		return err
	}

	if !fileServ.FileExists(fileName) {
		return HandleError(rw, FileNotFound, fmt.Sprintf("File '%v' doesn't exist", fileName))
	}

	file, err := fileServ.Read(fileName)
	if err != nil {
		return err
	}

	// Set the last block we expect to receive an ACK for.
	// In the case of empty files we still need to transmit one block
	// hence the lastBlock == 0, case
	lastBlock := (len(file.Data) / dataBlockSize)
	if lastBlock*dataBlockSize < len(file.Data) || lastBlock == 0 {
		lastBlock++
	}

	readSession := &ReadSession{
		rw:           rw,
		file:         file,
		currBlock:    1,
		lastBlock:    lastBlock,
		fileComplete: false,
		timeoutCount: 0,
	}

	logrus.Infof("[Read Session %v]: Start for file '%v'", remoteAddr.Port, file.Name)

	// Main work loop with bounded timeouts
	for readSession.timeoutCount < timeoutCountMax {
		if err = readSession.Start(); err != nil {
			if isTimeout(err) {
				logrus.Infof("[Read Session %v]: timeout %d", remoteAddr.Port, readSession.timeoutCount)
				readSession.timeoutCount++
			} else {
				return err
			}
		} else {
			return err
		}
	}

	return err
}

func (s *ReadSession) Start() error {
	// Write initial data packet
	s.writeData()

	// Read packets continuously until the last packet is ACKed.
	// See the ACK() method below
	for {
		if s.fileComplete {
			logrus.Infof("[Read Session]: completed file '%v' with %v bytes", s.file.Name, len(s.file.Data))
			return nil
		}

		if bytes, _, err := s.rw.Read(); err != nil {
			return err
		} else {
			if err := HandleTftpPackets(s, s.rw.remoteAddr, bytes); err != nil {
				return err
			}
		}
	}
}

// Get the next block of data from the file depending
// on the current block ID
func (s *ReadSession) getData() []byte {
	// Blocks are 1 indexed
	currBlock32 := int(s.currBlock)
	start := (currBlock32 - 1) * dataBlockSize
	if start >= len(s.file.Data) {
		return []byte{}
	}

	end := start + dataBlockSize
	if len(s.file.Data) <= end {
		return s.file.Data[start:]
	}

	return s.file.Data[start:end]
}

// Build the next Data packet for transmission
func (s *ReadSession) getDataPacket() (*DataPacket, error) {
	if bytes, err := convertIntToBytes(s.currBlock); err != nil {
		return nil, err
	} else {
		data := s.getData()
		blockArr := [2]byte{bytes[0], bytes[1]}
		return NewDataPacket(blockArr, data), nil
	}
}

// Send the next data packet to the requestor
func (s *ReadSession) writeData() error {
	if data, err := s.getDataPacket(); err != nil {
		return err
	} else {
		_, err = s.rw.Write(data.bytes)
		return err
	}
}

func (s *ReadSession) ReadReq(addr *net.UDPAddr, file string, mode string) error {
	return errors.New("ReadReq operations are not supported on read handlers")
}

func (s *ReadSession) WriteReq(addr *net.UDPAddr, file string, mode string) error {
	return errors.New("WriteReq operations are not supported on read handlers")
}

func (s *ReadSession) Data(block uint16, data []byte) error {
	return errors.New("Data operations are not supported on read handlers")
}

func (s *ReadSession) Ack(block uint16) error {
	if int(block) == s.lastBlock {
		s.fileComplete = true
		return nil
	}

	s.currBlock++
	return s.writeData()
}

func (s *ReadSession) Err(code uint16, msg string) error {
	logrus.Infof("Received Error with code %v and message %v", code, msg)
	return errors.New(fmt.Sprintf("Received Error with code %v and message %v", code, msg))
}
