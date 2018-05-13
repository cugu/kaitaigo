package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

	. "test_formats"
)

func assertInstanceEqualInt(t *testing.T, expected int, instCall func() (int, error)) {
	actual, err := instCall()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, expected, actual)
}

func assertInstanceEqualUint8(t *testing.T, expected uint8, instCall func() (uint8, error)) {
	actual, err := instCall()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, expected, actual)
}

func assertInstanceEqualString(t *testing.T, expected string, instCall func() (string, error)) {
	actual, err := instCall()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, expected, actual)
}

func TestExpr2(t *testing.T) {
	f, err := os.Open("../../src/str_encodings.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)

	var r Expr2
	err = r.Read(s, &r, &r)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 10, r.Str1.LenOrig)
	assertInstanceEqualInt(t, 7, r.Str1.LenMod)
	assert.EqualValues(t, "Some AS", r.Str1.Str)

	assertInstanceEqualInt(t, 7, r.Str1Len)
	assertInstanceEqualInt(t, 7, r.Str1LenMod)

	assertInstanceEqualUint8(t, 0x49, r.Str1Byte1)
	assertInstanceEqualInt(t, 0x49, r.Str1Avg)
	assertInstanceEqualString(t, "e", r.Str1Char5)

	str1Tuple5, err := r.Str1Tuple5()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, 0x65, str1Tuple5.Byte0)
	assert.EqualValues(t, 0x20, str1Tuple5.Byte1)
	assert.EqualValues(t, 0x41, str1Tuple5.Byte2)
	assertInstanceEqualInt(t, 0x30, str1Tuple5.Avg)

	str2Tuple5, err := r.Str2Tuple5()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, 0x65, str2Tuple5.Byte0)
	assert.EqualValues(t, 0x20, str2Tuple5.Byte1)
	assert.EqualValues(t, 0x41, str2Tuple5.Byte2)
	assertInstanceEqualInt(t, 0x30, str2Tuple5.Avg)
}
