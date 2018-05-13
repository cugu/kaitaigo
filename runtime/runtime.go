package runtime

import (
	"encoding/binary"
	"io"
)

type Kaitai interface {
	Decode(r io.ReadSeeker) error
}

type decoder struct {
	byteOrder binary.ByteOrder
	reader    io.ReadSeeker
	err       error
}

func (d *decoder) decode(i interface{}) {
	if d.err != nil {
		return
	}

	switch v := i.(type) {
	case Kaitai:
		d.err = v.Decode(d.reader)
	default:
		d.err = binary.Read(d.reader, d.byteOrder, i)
	}
}

func (d *decoder) decodePos(i interface{}, pos int) {
	if d.err != nil {
		return
	}
	_, d.err = reader.Seek(pos, io.SeekStart)
	decode(i)
}

func (d *decoder) decodeValue(i interface{}, value string) {
	if d.err != nil {
		return
	}
	_, d.err = reader.Seek(pos, io.SeekStart)
	decode(i)
}
