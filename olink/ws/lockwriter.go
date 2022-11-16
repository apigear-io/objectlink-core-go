package ws

import (
	"io"
	"sync"
)

// LockWriter is a writer that locks the underlying writer using a mutex.
type LockWriter struct {
	sync.Mutex
	writer io.Writer
}

func NewLockWriter(writer io.Writer) *LockWriter {
	return &LockWriter{
		writer: writer,
	}
}

func (l *LockWriter) Write(p []byte) (n int, err error) {
	l.Lock()
	defer l.Unlock()
	return l.writer.Write(p)
}
