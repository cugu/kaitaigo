package position_in_seq

import (
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/cugu/kaitai.go/runtime"
)

func TestPositionInSeq(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/position_in_seq.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r PositionInSeq
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 3, r.Header().QtyNumbers())
	assert.EqualValues(t, []runtime.Uint8{1, 2, 3}, r.Numbers())
}
