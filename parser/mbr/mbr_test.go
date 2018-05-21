package mbr

import (
	"os"
	"testing"
	"log"

	"github.com/stretchr/testify/assert"

	ks "gitlab.com/cugu/kaitai.go/runtime"
)

func TestMBR(t *testing.T) {
	f, err := os.Open("../../testdata/evidence/filesystem/mbr_fat16.dd")
	defer f.Close()

	if err != nil {
		t.Fatal(err)
	}

	d := ks.NewDecoder(f)
	r := Mbr{}
	r.Init(d, nil, nil)
	d.Decode(&r)
	if d.Err != nil {
		t.Fatal(d.Err)
	}
	log.Printf("%#v\n", r)

	assert.EqualValues(t, 128, r.Partitions[0].LbaStart)
	assert.EqualValues(t, 34816, r.Partitions[0].NumSectors)
	assert.EqualValues(t, 14, r.Partitions[0].PartitionType)
}


func BenchmarkMBR(b *testing.B) {
	for n := 0; n < b.N; n++ {
		f, _ := os.Open("../../testdata/evidence/filesystem/mbr_fat16.dd")
		d := ks.NewDecoder(f)
		var r Mbr
		d.Decode(&r)
		f.Close()
	}
}
