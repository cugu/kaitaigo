package repeat_until_s4

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/cugu/kaitai.go/runtime"
)

func TestRepeatUntilS4(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/repeat_until_s4.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r RepeatUntilS4
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, []runtime.Int32{66, 4919, -251658241, -1}, r.Entries())
	assert.EqualValues(t, "foobar", r.Afterall())
}
