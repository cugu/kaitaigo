package runtime

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"reflect"
	"strings"
	"fmt"
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
	s := &Decoder{reader, binary.LittleEndian, nil}
	return s
}

type KSYDecoder interface {
	KSYDecode(*Decoder) error
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

func (dec *Decoder) Decode(element interface{}) (value reflect.Value) {
	return dec.DecodeAncestors(element, reflect.Value{}, reflect.Value{})
}

func (dec *Decoder) DecodeAncestors(element interface{}, parent reflect.Value, root reflect.Value) (value reflect.Value) {
	// skip if previous errors
	if dec.Err != nil {
		return
	}

	// check if pointer
	pointer := reflect.ValueOf(element)
	if pointer.Type().Kind() != reflect.Ptr {
		dec.Err = errors.New("attempt to decode into a non-pointer")
		return
	}

	// get stored value
	value = decAlloc(pointer)

	// run KSYDecode if exists
	decoderType := reflect.TypeOf((*KSYDecoder)(nil)).Elem()
	if pointer.Type().Implements(decoderType) {
		dec.Err = element.(KSYDecoder).KSYDecode(dec)
		return value
	}

	// check if value can be set
	if !value.CanSet() {
		dec.Err = errors.New("Value cannot be set!!" + value.String())
		return
	}

	switch value.Kind() {
	case reflect.Array:
		switch value.Type().Elem().Kind() {
		case reflect.Bool,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			// array of builtin types
			dec.Err = binary.Read(dec, dec.ByteOrder, element)
			if dec.Err != nil {
				return
			}
		default:
			// other array
			for i := 0; i < value.Len(); i++ {
				item := dec.DecodeAncestors(value.Index(i).Addr().Interface(), parent, root)
				value.Index(i).Set(item)
			}
		}
	case reflect.Struct:
		// struct

		// set dec, parent and root
		decField := value.FieldByName("Dec")
		decField.Set(reflect.ValueOf(dec))
		startField := value.FieldByName("Start")
		pos, err := dec.Seek(0, io.SeekCurrent)
		if err != nil {
			dec.Err = err
			return
		}
		startField.Set(reflect.ValueOf(pos))
		parentField := value.FieldByName("Parent")
		if !parent.IsValid() {
			parent = pointer
		}
		if parentField.IsNil() {
			parentField.Set(parent)
		}
		rootField := value.FieldByName("Root")
		if !root.IsValid() {
			root = pointer
		}
		if rootField.IsNil() {
			rootField.Set(root)
		}


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
						unsupportedError := fmt.Sprintf("Unsupported flag %q in tag %q of type %s", flag, tag, value)
						dec.Err = errors.New(unsupportedError)
						return
					}
				}
			}

			if attribute {
				substruct := dec.DecodeAncestors(field.Addr().Interface(), value.Addr(), root)
				field.Set(substruct)
			}

		}
		// value = value.Addr()
	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		// builtin types
		dec.Err = binary.Read(dec, dec.ByteOrder, element)
	default:
		log.Printf("Type unknown %+v\n", value)
	}

	return
}

func (dec *Decoder) DecodePos(element interface{}, offset int64, whence int, parent interface{}, root interface{}) {
	if dec.Err != nil {
		return
	}
	_, dec.Err = dec.Seek(offset, whence)
	dec.DecodeAncestors(element, reflect.ValueOf(parent), reflect.ValueOf(root))
}

