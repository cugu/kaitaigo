package gpt

import (
	"os"
	"testing"
	"log"

	//"github.com/stretchr/testify/assert"

	ks "gitlab.com/cugu/kaitai.go/runtime"
)

func TestGPT(t *testing.T) {
	f, err := os.Open("../../testdata/evidence/filesystem/gpt_apfs.dd")
	defer f.Close()

	if err != nil {
		t.Fatal(err)
	}

	d := ks.NewDecoder(f)
	var r Gpt
	r.Init(d, &r, &r)
	d.Decode(&r)
	if d.Err != nil {
		t.Fatal(d.Err)
	}
	log.Printf("%#v\n", r)

	/*assert.EqualValues(t, 128, r.Primary.LbaStart)
	assert.EqualValues(t, 34816, r.Partitions[0].NumSectors)
	assert.EqualValues(t, 14, r.Partitions[0].PartitionType)*/
}


func BenchmarkGPT(b *testing.B) {
	for n := 0; n < b.N; n++ {
		f, _ := os.Open("../../testdata/evidence/filesystem/gpt_apfs.dd")
		d := ks.NewDecoder(f)
		var r Gpt
		r.Init(d, &r, &r)
		d.Decode(&r)
		f.Close()
	}
}
