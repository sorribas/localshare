package writercounter

import "io"

type WriterCounter struct {
	Count  int64
	writer io.Writer
}

func NewWriterCounter(writer io.Writer) *WriterCounter {
	return &WriterCounter{writer: writer}
}

func (wc *WriterCounter) Write(bts []byte) (int, error) {
	n, err := wc.writer.Write(bts)
	wc.Count += int64(n)
	return n, err
}
