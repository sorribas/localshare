package localsharelib

import "bytes"
import "io"
import "os"

type File interface {
	Name() string
	Open() (io.Reader, error)
}

type FsFile struct {
	path string
	name string
}

func NewFsFile(path, name string) *FsFile {
	return &FsFile{path, name}
}

func (fsFile *FsFile) Name() string {
	return fsFile.name
}

func (fsFile *FsFile) Open() (io.Reader, error) {
	return os.Open(fsFile.path)
}

func (instance *LocalshareInstance) AddFile(f File) {
	instance.files[f.Name()] = f
}

func (instance *LocalshareInstance) SharedFiles() map[string]File {
	return instance.files
}

type InMemoryFile struct {
	name string
	data []byte
}

func NewInMemoryFile(name string, data []byte) *InMemoryFile {
	return &InMemoryFile{name, data}
}

func (imf *InMemoryFile) Name() string {
	return imf.name
}

func (imf *InMemoryFile) Open() (io.Reader, error) {
	return bytes.NewReader(imf.data), nil
}
