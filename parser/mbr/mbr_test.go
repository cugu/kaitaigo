package mbr

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMBR(t *testing.T) {
	file, err := os.Open("../../testdata/evidence/filesystem/mbr_fat16.dd")
	defer file.Close()

	if err != nil {
		t.Fatal(err)
	}

	mbr := Mbr{}
	err = mbr.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	p0 := mbr.Partitions()[0]
	assert.EqualValues(t, 128, p0.LbaStart())
	assert.EqualValues(t, 34816, p0.NumSectors())
	assert.EqualValues(t, 14, p0.PartitionType())
}

func BenchmarkMBR(b *testing.B) {
	for n := 0; n < b.N; n++ {
		file, _ := os.Open("../../testdata/evidence/filesystem/mbr_fat16.dd")
		mbr := Mbr{}
		err := mbr.Decode(file)
		if err != nil {
			b.Fatal(err)
		}
		file.Close()
	}
}
