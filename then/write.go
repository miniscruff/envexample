package then

import (
	"bytes"
	"errors"
	"testing"
)

// ErrWriter is a simple struct that will return an error when trying to Write
type ErrWriter struct {
	err error
}

func NewErrWriter() *ErrWriter {
	return &ErrWriter{
		err: errors.New("error from ErrWriter"),
	}
}

func (w *ErrWriter) Write(data []byte) (int, error) {
	return 0, w.err
}

// Raised will assert the error value is the one we would of returned if
// Write was called.
func (w *ErrWriter) Raised(t *testing.T, err error) {
	t.Helper()
	Err(t, w.err, err)
}

// CountWriter is a simple struct that will successfully write to a byte buffer
// a specified number of times, and then it will error
type CountWriter struct {
	err    error
	writer bytes.Buffer
	count  int
}

func NewCountWriter(count int) *CountWriter {
	return &CountWriter{
		err:   errors.New("error from CountWriter"),
		count: count,
	}
}

func (w *CountWriter) Write(value []byte) (int, error) {
	if w.count <= 0 {
		return 0, w.err
	}

	w.count--

	return w.writer.Write(value)
}

// Raised will assert the error value is the one we would of returned if
// Write was called too many times.
func (w *CountWriter) Raised(t *testing.T, err error) {
	t.Helper()
	Err(t, w.err, err)
}
