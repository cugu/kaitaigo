package docstrings_docref

import (
	"os"
	"testing"
)

func TestDocstringsDocref(t *testing.T) {
	f, err := os.Open("../../../testdata/kaitai/fixed_struct.bin")
	if err != nil {
		t.Fatal(err)
	}

	var r DocstringsDocref
	err = r.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

}
