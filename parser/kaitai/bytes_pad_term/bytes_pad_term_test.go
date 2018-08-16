package bytes_pad_term

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesPadTerm(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/str_pad_term.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r BytesPadTerm
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, "str1", r.StrPad())
	assert.EqualValues(t, "str2foo", r.StrTerm())
	assert.EqualValues(t, "str+++3bar+++", r.StrTermAndPad())
	assert.EqualValues(t, "str4baz@", r.StrTermInclude())
}
