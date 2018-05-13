package spec

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	. "test_formats"
)

func TestPositionInSeq(t *testing.T) {
	f, err := os.Open("../../src/position_in_seq.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)
	var r PositionInSeq
	err = r.Read(s, &r, &r)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, []uint8{1, 2, 3}, r.Numbers)
}
