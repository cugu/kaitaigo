package spec

import (
	"os"
	"testing"
)

func TestDocstrings(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/fixed_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r Docstrings
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

}
