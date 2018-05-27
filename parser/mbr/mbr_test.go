package mbr

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	ks "gitlab.com/cugu/kaitai.go/runtime"
)

func TestMBR(t *testing.T) {
	file, err := os.Open("../../testdata/evidence/filesystem/mbr_fat16.dd")
	defer file.Close()

	if err != nil {
		t.Fatal(err)
	}

	dec := ks.NewDecoder(file)
	mbr := Mbr{}
	dec.Decode(&mbr)
	if dec.Err != nil {
		t.Fatal(dec.Err)
	}

	p0 := mbr.GetPartitions()[0]
	assert.EqualValues(t, 128, p0.GetLbaStart())
	assert.EqualValues(t, 34816, p0.GetNumSectors())
	assert.EqualValues(t, 14, p0.GetPartitionType())
}

func BenchmarkMBR(b *testing.B) {
	for n := 0; n < b.N; n++ {
		file, _ := os.Open("../../testdata/evidence/filesystem/mbr_fat16.dd")
		dec := ks.NewDecoder(file)
		mbr := Mbr{}
		dec.Decode(&mbr)
		file.Close()
	}
}
