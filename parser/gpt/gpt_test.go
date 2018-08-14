package gpt

import (
	"encoding/binary"
	// "fmt"
	"os"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/stretchr/testify/assert"

		"gitlab.com/cugu/kaitai.go/runtime"
)

func TestGPT(t *testing.T) {
	file, err := os.Open("../../testdata/evidence/filesystem/gpt_apfs.dd")
	defer file.Close()

	if err != nil {
		t.Fatal(err)
	}

	gpt := Gpt{}
	err = gpt.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	primary := gpt.Primary()
	assert.EqualValues(t, [8]runtime.Byte{0x45, 0x46, 0x49, 0x20, 0x50, 0x41, 0x52, 0x54}, primary.Signature())
	partitions := primary.Entries()[0]
	name := partitions.Name()
	bytes := []byte{}
	for i := 0; i < len(name); i += 1 {
		bytes = append(bytes, byte(name[i]))
	}
	u16 := []uint16{}
	for i := 0; i < len(bytes); i += 2 {
		u16 = append(u16, binary.LittleEndian.Uint16(bytes[i:i+2]))
	}
	assert.EqualValues(t, "disk image", strings.Trim(string(utf16.Decode(u16)), "\x00"))
	assert.EqualValues(t, 40, partitions.FirstLba())
}

func BenchmarkGPT(b *testing.B) {
	for n := 0; n < b.N; n++ {
		file, _ := os.Open("../../testdata/evidence/filesystem/gpt_apfs.dd")
		gpt := Gpt{}
	err := gpt.Decode(file)
	if err != nil {
		b.Fatal(err)
	}
		file.Close()
	}
}
