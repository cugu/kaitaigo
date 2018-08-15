package expr_1

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpr1(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/str_encodings.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r Expr1
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 10, *r.LenOf1())
	assert.EqualValues(t, 8, *r.LenOf1Mod())
	assert.EqualValues(t, "Some ASC", *r.Str1())
	assert.EqualValues(t, 8, *r.Str1Len())
}
