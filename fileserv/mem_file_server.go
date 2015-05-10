package fileserv

import (
	"errors"
	"fmt"
	. "sync"
)

type InMemFileServer struct {
	fileDir map[string]*File
	mutex   *Mutex
}

func NewMemFileServer() *InMemFileServer {
	dir := make(map[string]*File)
	return &InMemFileServer{
		fileDir: dir,
		mutex:   &Mutex{},
	}
}

func (s *InMemFileServer) Write(file *File) error {
	defer s.mutex.Unlock()
	s.mutex.Lock()

	if s.FileExists(file.Name) {
		return errors.New(fmt.Sprintf("File '%v' already exists", file.Name))
	}

	s.fileDir[file.Name] = file
	return nil
}

func (s *InMemFileServer) Read(file string) (*File, error) {
	if f, ok := s.fileDir[file]; ok {
		return f, nil
	}

	return &File{}, errors.New(fmt.Sprintf("File '%v' doesn't exist", file))
}

func (s *InMemFileServer) FileExists(file string) bool {
	if _, ok := s.fileDir[file]; ok {
		return true
	}

	return false
}
