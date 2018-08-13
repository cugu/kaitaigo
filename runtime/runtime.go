package runtime

import (
	"encoding/binary"
	"io"
	"reflect"
)

func IsNull(value interface{}) bool {
	return reflect.DeepEqual(reflect.Zero(reflect.TypeOf(value)).Interface(), value)
}

type Decoder struct {
	io.ReadSeeker
	ByteOrder binary.ByteOrder
	Err       error
}

func NewDecoder(reader io.ReadSeeker) *Decoder {
	return &Decoder{reader, binary.LittleEndian, nil}
}

func (dec *Decoder) BinaryRead(value interface{}) {
	if dec.Err != nil {
		return
	}
	dec.Err = binary.Read(dec, dec.ByteOrder, value)
}

type KSYDecoder interface {
	Decode(*Decoder) error
	DecodeAncestors(*Decoder, interface{}, interface{}) error
	DecodePos(*Decoder, int64, int, interface{}, interface{}) error
}
