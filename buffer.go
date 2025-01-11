package logstream

import (
	"bytes"
	"io"
	"sync"
)

type bufferedWriter struct {
	writer io.Writer
	buffer bytes.Buffer
	mutex  sync.Mutex
}

func newBufferedWriter(w io.Writer) *bufferedWriter {
	return &bufferedWriter{writer: w}
}

func (w *bufferedWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.buffer.Write(p)
}

func (w *bufferedWriter) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	_, err := w.writer.Write(w.buffer.Bytes())
	w.buffer.Reset()
	return err
}
