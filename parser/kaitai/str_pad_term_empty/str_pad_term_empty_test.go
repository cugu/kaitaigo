package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrPadTermEmpty(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/str_pad_term_empty.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r StrPadTermEmpty
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, "", r.StrPad)
	assert.EqualValues(t, "", r.StrTerm)
	assert.EqualValues(t, "", r.StrTermAndPad)
	assert.EqualValues(t, "@", r.StrTermInclude)
}
