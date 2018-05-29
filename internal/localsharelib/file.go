package localsharelib

import "bytes"
import "io"
import "os"

type File interface {
	Name() string
	Open() (io.Reader, error)
	Size() int64
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

func (fsFile *FsFile) Size() int64 {
	stat, err := os.Stat(fsFile.path)
	if err != nil {
		return 0
	}
	return stat.Size()
}

func (instance *LocalshareInstance) AddFile(f File) {
	instance.files[f.Name()] = f
	// announce the new file list hash through mdns
	instance.mdnsServer.SetText([]string{hash(instance.files)})
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

func (imf *InMemoryFile) Size() int64 {
	return int64(len(imf.data))
}
