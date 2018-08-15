package runtime

import (
	"encoding/binary"
	"io/ioutil"
)

var RTDecoder *Decoder

type Byte byte

func (v *Byte) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Bytes []byte

func (v *Bytes) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	lv, err := ioutil.ReadAll(RTDecoder)
	if err != nil {
		RTDecoder.Err = err
		return
	}
	*v = lv
}

type ByteSlice []byte

func (v *ByteSlice) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	_, err := RTDecoder.Read(*v)
	if err != nil {
		RTDecoder.Err = err
	}
}

type Uint8 uint8

func (v *Uint8) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Uint16 uint16

func (v *Uint16) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Uint32 uint32

func (v *Uint32) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Uint64 uint64

func (v *Uint64) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Int8 int8

func (v *Int8) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Int16 int16

func (v *Int16) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Int32 int32

func (v *Int32) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}

type Int64 int64

func (v *Int64) DecodeAncestors(parent interface{}, root interface{}) {
	if RTDecoder.Err != nil {
		return
	}
	RTDecoder.Err = binary.Read(RTDecoder, RTDecoder.ByteOrder, v)
}
