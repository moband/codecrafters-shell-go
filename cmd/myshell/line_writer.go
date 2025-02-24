package main

import "io"

type lineWriter struct {
	w   io.Writer
	buf []byte
}

func (lw *lineWriter) Write(p []byte) (n int, err error) {
	if lw.buf == nil {
		lw.buf = make([]byte, 0, 4096)
	}

	for _, b := range p {
		if b == '\n' {
			lw.buf = append(lw.buf, '\n', '\r')
		} else {
			lw.buf = append(lw.buf, b)
		}
	}

	lw.w.Write(lw.buf)
	lw.buf = lw.buf[:0]
	return len(p), err
}
