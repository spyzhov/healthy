package internal

import (
	"bytes"
	"io"

	"github.com/spyzhov/safe"
)

type ReaderCloserCallback struct {
	ReadCloser io.ReadCloser
	Callback   func([]byte, int) ([]byte, error)
	buffer     bytes.Buffer
}

func NewReaderCloserCallback(reader io.ReadCloser, callback func([]byte, int) ([]byte, error)) *ReaderCloserCallback {
	return &ReaderCloserCallback{
		ReadCloser: reader,
		Callback:   callback,
	}
}

func (r *ReaderCloserCallback) Read(p []byte) (n int, err error) {
	if r == nil || safe.IsNil(r.ReadCloser) {
		return 0, io.EOF
	}
	var (
		m    int
		tmp  []byte
		read error
	)
	// read previous data:
	m, err = r.buffer.Read(p)
	if err == io.EOF {
		err = nil
		r.buffer.Reset()
	}
	if err != nil || m == cap(p) {
		return m, err
	}
	// read new data:
	n, read = r.ReadCloser.Read(p[m:])
	if read != nil {
		if read != io.EOF {
			return n + m, read
		}
	}
	// extend data:
	if r.Callback != nil {
		tmp, err = r.Callback(p[m:], n)
		if err != nil {
			return n + m, err
		}
		n = copy(p[m:], tmp)
		if n < len(tmp) {
			r.buffer.Write(tmp[n:])
			read = nil
		}
	}
	return n + m, read
}

func (r *ReaderCloserCallback) Close() error {
	if r == nil || safe.IsNil(r.ReadCloser) {
		return nil
	}
	return r.ReadCloser.Close()
}
