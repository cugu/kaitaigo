package repeat_eos_u4

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatEosU4(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/repeat_eos_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r RepeatEosU4
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, []uint32{0, 66, 66, 2069}, r.Numbers())
}
