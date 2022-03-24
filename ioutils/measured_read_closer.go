package ioutils

import (
	"io"
)

type MeasuredReadCloser struct {
	io.ReadCloser
	ByteLength int
}

func NewMeasuredReadCloser(rc io.ReadCloser) *MeasuredReadCloser {
	return &MeasuredReadCloser{ReadCloser: rc}
}

func (rc *MeasuredReadCloser) Read(p []byte) (int, error) {
	i, err := rc.ReadCloser.Read(p)
	rc.ByteLength += i
	return i, err
}
