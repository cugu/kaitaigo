package kgruntime

import (
	"encoding/binary"
	"io"
)

type KSYDecoder interface {
	KSYDecode(io.ReadSeeker) error
}

type Decoder struct {
	ByteOrder binary.ByteOrder
	Reader    io.ReadSeeker
	Err       error
}

func (d *Decoder) Decode(i interface{}) {
	if d.Err != nil {
		return
	}

	switch v := i.(type) {
	case KSYDecoder:
		d.Err = v.KSYDecode(d.Reader)
	default:
		d.Err = binary.Read(d.Reader, d.ByteOrder, i)
	}
}

func (d *Decoder) DecodePos(i interface{}, pos int64) {
	if d.Err != nil {
		return
	}
	_, d.Err = d.Reader.Seek(pos, io.SeekStart)
	d.Decode(i)
}
