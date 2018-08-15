package runtime

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

var RTDecoder *Decoder

type Byte byte

func (v *Byte) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Byte) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
	return
}

type Bytes []byte

func (v *Bytes) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Bytes) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	lv, err := ioutil.ReadAll(RTDecoder)
	if err != nil {
		RTDecoder.Err = err
		return err
	}
	*v = lv
	return
}

type ByteSlice []byte

func (v *ByteSlice) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *ByteSlice) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	_, err = RTDecoder.Read(*v)
	if err != nil {
		RTDecoder.Err = err
	}
	return
}

type Uint8 uint8

func (v *Uint8) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Uint8) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
	return
}

type Uint16 uint16

func (v *Uint16) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Uint16) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
	return
}

type Uint32 uint32

func (v *Uint32) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Uint32) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
	return
}

type Uint64 uint64

func (v *Uint64) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Uint64) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
	return
}

type Int8 int8

func (v *Int8) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Int8) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
	return
}

type Int16 int16

func (v *Int16) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Int16) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
	return
}

type Int32 int32

func (v *Int32) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Int32) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
	return
}

type Int64 int64

func (v *Int64) Decode(reader io.ReadSeeker) (err error) {
	return v.DecodeAncestors(v, v)
}

func (v *Int64) DecodeAncestors(parent interface{}, root interface{}) (err error) {
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
	return
}
