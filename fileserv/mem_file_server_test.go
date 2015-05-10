package fileserv

import (
	"bytes"
	"testing"
	"time"
)

func TestFileWrite(t *testing.T) {
	serv := NewMemFileServer()
	file := File{
		Name: "foo",
		Data: []byte{0, 1, 2, 3, 4},
	}

	err := serv.Write(&file)

	if err != nil {
		t.Errorf("Failed to write %v, returned %v", file.Name, err)
	}
}

func TestFileRead(t *testing.T) {
	serv := NewMemFileServer()
	file := File{
		Name: "foo",
		Data: []byte{0, 1, 2, 3, 4},
	}

	err := serv.Write(&file)
	if err != nil {
		t.Errorf("Failed to write %v, returned %v", file.Name, err)
	}

	recvFile, err := serv.Read(file.Name)
	if err != nil {
		t.Errorf("Failed to read %v, returned %v", file.Name, err)
	}

	if recvFile.Name != file.Name {
		t.Errorf("Original file name '%v' != received file name '%v'", file.Name, recvFile.Name)
	}

	if !bytes.Equal(file.Data, recvFile.Data) {
		t.Errorf("Original file data '%v' != received file data '%v'", file.Data, recvFile.Data)
	}
}

func TestFileOverWrite(t *testing.T) {
	serv := NewMemFileServer()
	file := File{
		Name: "foo",
		Data: []byte{0, 1, 2, 3, 4},
	}

	err := serv.Write(&file)

	if err != nil {
		t.Errorf("Failed to write %v, returned %v", file.Name, err)
	}

	err = serv.Write(&file)

	if err == nil {
		t.Errorf("Overwrite of file '%v' should have failed, returned %v", file.Name, err)
	}
}

func TestFileReadFailure(t *testing.T) {
	serv := NewMemFileServer()
	file := File{
		Name: "foo",
		Data: []byte{0, 1, 2, 3, 4},
	}

	_, err := serv.Read(file.Name)
	if err == nil {
		t.Errorf("Should have failed to read %v, returned %v", file.Name, err)
	}
}

func TestFileParallelWriteTest(t *testing.T) {
	serv := NewMemFileServer()
	file := File{
		Name: "foo",
		Data: []byte{0, 1, 2, 3, 4},
	}

	for i := 0; i < 1000; i++ {
		go serv.Write(&file)
	}

	time.Sleep(3 * time.Second)

	if len(serv.fileDir) != 1 {
		t.Errorf("Expected fileDir size to be 1, received: %v", len(serv.fileDir))
	}
}
