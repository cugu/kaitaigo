package spec

import (
	"os"
	"testing"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	. "test_formats"
	"github.com/stretchr/testify/assert"
)

func TestRepeatUntilS4(t *testing.T) {
	f, err := os.Open("../../src/repeat_until_s4.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)
	var r RepeatUntilS4
	err = r.Read(s, &r, &r)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, []int32{66, 4919, -251658241, -1}, r.Entries)
	assert.EqualValues(t, "foobar", r.Afterall)
}
