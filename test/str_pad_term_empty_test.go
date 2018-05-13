package spec

import (
	"os"
	"testing"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"github.com/stretchr/testify/assert"

	. "test_formats"
)

func TestStrPadTermEmpty(t *testing.T) {
	f, err := os.Open("../../src/str_pad_term_empty.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)

	var r StrPadTermEmpty
	err = r.Read(s, &r, &r)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, "", r.StrPad)
	assert.EqualValues(t, "", r.StrTerm)
	assert.EqualValues(t, "", r.StrTermAndPad)
	assert.EqualValues(t, "@", r.StrTermInclude)
}
