package spec

import (
	"os"
	"testing"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"github.com/stretchr/testify/assert"

	. "test_formats"
)

func TestBytesPadTerm(t *testing.T) {
	f, err := os.Open("../../src/str_pad_term.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)

	var r BytesPadTerm
	err = r.Read(s, &r, &r)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, "str1", r.StrPad)
	assert.EqualValues(t, "str2foo", r.StrTerm)
	assert.EqualValues(t, "str+++3bar+++", r.StrTermAndPad)
	assert.EqualValues(t, "str4baz@", r.StrTermInclude)
}
