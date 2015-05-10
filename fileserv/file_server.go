package fileserv

type FileServer interface {
	Write(file *File) error
	Read(file string) (*File, error)
	FileExists(file string) bool
}

type File struct {
	Name string
	Data []byte
}
