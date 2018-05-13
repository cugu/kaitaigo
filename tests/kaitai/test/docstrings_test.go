package spec

import (
	"os"
	"testing"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

	. "test_formats"
)

func TestDocstrings(t *testing.T) {
	f, err := os.Open("../../src/fixed_struct.bin")
	if err != nil {
		t.Fatal(err)
	}
	s := kaitai.NewStream(f)

	var r Docstrings
	err = r.Read(s, &r, &r)
	if err != nil {
		t.Fatal(err)
	}

}
