package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

	. "test_formats"
)

func TestExpr1(t *testing.T) {
	f, err := os.Open("../../src/str_encodings.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)

	var r Expr1
	err = r.Read(s, &r, &r)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 10, r.LenOf1)

	lenOf1Mod, err := r.LenOf1Mod()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, 8, lenOf1Mod)

	assert.EqualValues(t, "Some ASC", r.Str1)
	
	str1Len, err := r.Str1Len()
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, 8, str1Len)
}
