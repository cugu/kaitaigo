package ks

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"reflect"
	"strings"
	"fmt"
)


func NewDecoder(reader io.ReadSeeker) *Stream {
	s := &Stream{reader: reader, ByteOrder: binary.LittleEndian}
	s.Pos, _ = s.reader.Seek(0, io.SeekCurrent)
	s.Size, _ = s.GetSize()
	return s
}

type Stream struct {
	reader io.ReadSeeker
	ByteOrder binary.ByteOrder
	Size int64
	Pos int64
	Err       error
}


func (k *Stream) Seek(offset int64, whence int) (n int64, err error) {
	n, err = k.reader.Seek(offset, whence)
	k.Pos = n
	k.Size, _ = k.GetSize()
	return
}

func (k *Stream) Read(data []byte) (n int, err error) {
	n, err = k.reader.Read(data)
	k.Pos, _ = k.reader.Seek(0, io.SeekCurrent)
	k.Size, _ = k.GetSize()
	return
}

func (k *Stream) GetSize() (int64, error) {
	// Go has no internal ReadSeeker function to get current ReadSeeker size,
	// thus we use the following trick.
	// Remember our current position
	curPos := k.Pos
	// Seek to the end of the File object
	_, err := k.reader.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	// Remember position, which is equal to the full length
	fullSize := k.Pos
	// Seek back to the current position
	_, err = k.reader.Seek(curPos, io.SeekStart)
	return fullSize, err
}

type KSYDecoder interface {
	KSYDecode(*Stream) error
}

// decAlloc takes a value and returns a settable value that can
// be assigned to. If the value is a pointer, decAlloc guarantees it points to storage.
// The callers to the individual decoders are expected to have used decAlloc.
// The individual decoders don't need to it.
func decAlloc(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

func dbg(s string, v reflect.Value) {

}

func (dec *Stream) Decode(element interface{}) (value reflect.Value) {
	if dec.Err != nil {
		return
	}

	pointer := reflect.ValueOf(element)
	if pointer.Type().Kind() != reflect.Ptr {
		dec.Err = errors.New("attempt to decode into a non-pointer")
		return
	}

	value = decAlloc(pointer)


	decoderType := reflect.TypeOf((*KSYDecoder)(nil)).Elem()
	if pointer.Type().Implements(decoderType) {
		dec.Err = element.(KSYDecoder).KSYDecode(dec)
		return value
	}


	if !value.CanSet() {
		panic("Value cannot be set!!" + value.String())
	}

	switch value.Kind() {
	case reflect.Array:
		switch value.Type().Elem().Kind() {
		case reflect.Bool,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			// array of builtin types
			dec.Err = binary.Read(dec.reader, dec.ByteOrder, element)
		default:
			// other array
			for i := 0; i < value.Len(); i++ {
				item := dec.Decode(value.Index(i).Addr().Interface())
				value.Index(i).Set(item)
			}
		}
	case reflect.Struct:
		// struct
		for i := 0; i < value.NumField(); i++ {
			attribute := false
			field := value.Field(i)
			tag := value.Type().Field(i).Tag.Get("ks")
			fields := strings.Split(tag, ",")
			if len(fields) > 1 {
				for _, flag := range fields[1:] {
					switch flag {
					case "attribute":
						attribute = true
					case "instance":
					default:
						dec.Err = errors.New(fmt.Sprintf("Unsupported flag %q in tag %q of type %s", flag, tag, value))
						return
					}
				}
				tag = fields[0]
			}

			if attribute {
				substruct := dec.Decode(field.Addr().Interface())
				value.Field(i).Set(substruct)
			}

		}
	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		// builtin types
		dec.Err = binary.Read(dec.reader, dec.ByteOrder, element)
	default:
		log.Printf("Type unknown %+v\n", value)
	}

	return
}

/*
func SimpleDecode() {
	switch v := i.(type) {
	case KSYDecoder:
		log.Println("kaitai type")
		d.Err = v.KSYDecode(d.Reader)
	case float32, float64, int16, int32, int64, int8, string, uint16, uint32, uint64, uint8:
		log.Println("builtin type")
		d.Err = binary.Read(d.Reader, d.ByteOrder, i)
	default:
}
*/
func (d *Stream) DecodePos(element interface{}, pos int64) {
	if d.Err != nil {
		return
	}
	_, d.Err = d.reader.Seek(pos, io.SeekStart)
	d.Decode(element)
}

