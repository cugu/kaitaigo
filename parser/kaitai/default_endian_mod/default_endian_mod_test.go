package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEndianMod(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/fixed_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r DefaultEndianMod
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, 0x4b434150, r.Main.One)
	assert.EqualValues(t, -52947, r.Main.Nest.Two)
	assert.EqualValues(t, 0x5041434b, r.Main.NestBe.Two)
}
