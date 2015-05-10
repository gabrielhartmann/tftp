package tftp

import (
	"net"
	"time"

	"github.com/Sirupsen/logrus"
)

type TftpReaderWriter struct {
	buf       []byte
	conn      *net.UDPConn
	localAddr *net.UDPAddr
	timeout   bool
}

const timeoutSec = 3

func NewTftpReaderWriter(remoteAddr *net.UDPAddr, timeout bool) (*TftpReaderWriter, error) {
	// Resolve UDP address
	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		logrus.Infof("Failed to resolve UDP address: %v", err)
		return nil, err
	}

	var conn *net.UDPConn = nil

	if remoteAddr == nil {
		// Listen on UDP connection
		conn, err = net.ListenUDP("udp", localAddr)
		if err != nil {
			logrus.Infof("Failed to UDP listen: %v", err)
			return nil, err
		} else {
			logrus.Infof("UDP local address: %v", conn.LocalAddr())
		}
	} else {
		// Connect to remote address
		conn, err = net.DialUDP("udp", localAddr, remoteAddr)
		if err != nil {
			return nil, err
		}
	}

	writer := &TftpReaderWriter{
		conn:      conn,
		buf:       make([]byte, 1024),
		localAddr: localAddr,
		timeout:   timeout,
	}

	return writer, nil
}

func (rw *TftpReaderWriter) Write(bytes []byte) (int, error) {
	rw.setDeadline()
	return rw.conn.Write(bytes)
}

func (rw *TftpReaderWriter) Read() ([]byte, *net.UDPAddr, error) {
	rw.setDeadline()

	// Read bytes into buffer
	length, addr, err := rw.conn.ReadFromUDP(rw.buf)
	if err != nil {
		logrus.Infof("Read error: ", err)
		return []byte{}, nil, err
	}

	// Copy bytes into new buffer
	copyBuf := make([]byte, length)
	copy(copyBuf, rw.buf[0:length])

	return copyBuf, addr, nil
}

func (rw *TftpReaderWriter) setDeadline() {
	if rw.timeout {
		rw.conn.SetDeadline(time.Now().Add(3 * time.Second))
	}
}
