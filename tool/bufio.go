package main

import (
	"fmt"
	"io"
	"os"
)

func main() {

	path := "./a.txt"
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	buf := NewBuf(f)

	for {
		bytes, err := buf.ReadLine()
		if err != nil {
			break
		}
		fmt.Println(string(bytes))
	}

}

type buf struct {
	fd   io.Reader
	buf  []byte
	r, w int
	eof  bool
}

func NewBuf(reader io.Reader) *buf {
	return &buf{
		fd:  reader,
		buf: make([]byte, 4096),
		r:   0,
		w:   0,
	}
}

func (b *buf) ReadLine() ([]byte, error) {
	sysRead := func() error {
		cnt, err := b.fd.Read(b.buf)
		b.r = 0
		b.w = cnt
		return err
	}

	bytes := make([]byte, 512)
	index := 0
	for {
		if b.r == b.w {
			if err := sysRead(); err != nil {
				if err == io.EOF {
					b.eof = true
				} else {
					return nil, err
				}
			}
		}
		for i := b.r; i < b.w; i++ {
			if b.buf[i] != '\n' {
				if index < len(bytes) {
					bytes[index] = b.buf[i]
				} else {
					tmp := make([]byte, len(bytes))
					bytes = append(bytes, tmp...)
					bytes[index] = b.buf[i]
				}
				index++
				b.r++
			} else {
				b.r++
				return bytes[:index], nil
			}
		}
		if b.eof {
			return bytes[:index], io.EOF
		}
	}

	return nil, nil
}
