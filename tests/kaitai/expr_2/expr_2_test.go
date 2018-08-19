package expr_2

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpr2(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/str_encodings.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r Expr2
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 10, r.Str1().LenOrig())
	assert.EqualValues(t, 7, r.Str1().LenMod())
	assert.EqualValues(t, "Some AS", r.Str1().Str())

	assert.EqualValues(t, 7, r.Str1Len())
	assert.EqualValues(t, 7, r.Str1LenMod())

	assert.EqualValues(t, 0x49, r.Str1Byte1())
	assert.EqualValues(t, 0x49, r.Str1Avg())
	assert.EqualValues(t, "e", r.Str1Char5())

	str1Tuple5 := r.Str1Tuple5()
	assert.EqualValues(t, 0x65, str1Tuple5.Byte0())
	assert.EqualValues(t, 0x20, str1Tuple5.Byte1())
	assert.EqualValues(t, 0x41, str1Tuple5.Byte2())
	assert.EqualValues(t, 0x30, str1Tuple5.Avg())

	str2Tuple5 := r.Str2Tuple5()
	assert.EqualValues(t, 0x65, str2Tuple5.Byte0())
	assert.EqualValues(t, 0x20, str2Tuple5.Byte1())
	assert.EqualValues(t, 0x41, str2Tuple5.Byte2())
	assert.EqualValues(t, 0x30, str2Tuple5.Avg())
}
