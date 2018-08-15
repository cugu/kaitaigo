// Autogenerated from KST: please remove this line if doing any edits by hand!

package spec

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatNStrz(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/repeat_n_strz.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r RepeatNStrz
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 2, r.Qty)
	assert.EqualValues(t, []string{"foo", "bar"}, r.Lines)
}
