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

type KSYDecoder interface {
	Decode(io.ReadSeeker) error
	DecodeAncestors(*Decoder, interface{}, interface{}) error
	DecodePos(*Decoder, int64, int, interface{}, interface{}) error
}
