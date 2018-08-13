package apfs

import (
    "os"
    "testing"
    "io"
    "log"

    "github.com/stretchr/testify/assert"

    ks "gitlab.com/cugu/kaitai.go/runtime"
)



func TestAPFS(t *testing.T) {
    file, err := os.Open("../../testdata/evidence/filesystem/gpt_apfs.dd")
    defer file.Close()

    if err != nil {
        t.Fatal(err)
    }

    filesystem := io.NewSectionReader(file, 40 * 512, 39024 * 512)

    dec := ks.NewDecoder(filesystem)
    apfs := Apfs{}
    dec.Decode(&apfs)
    if dec.Err != nil {
        t.Fatal(dec.Err)
    }

    p0 := apfs.GetBlock0()
    body := p0.GetBody()
    containerSuperblock := body.(ContainerSuperblock)
    // log.Printf("containerSuperblock: %#v", containerSuperblock)
    blocksize := containerSuperblock.GetBlockSize()
    assert.EqualValues(t, 4096, blocksize)

    assert.EqualValues(t, 0x949, containerSuperblock.GetOmapOid())
    filesystem.Seek(int64(containerSuperblock.GetOmapOid()) * int64(blocksize), io.SeekStart)
    omap := Btree{Root: &apfs}
    dec.Decode(&omap)
    log.Printf("omap: %#v", omap)

    filesystem.Seek(int64(omap.GetTreeRoot()) * int64(blocksize), io.SeekStart)
    omapnode := Node{Root: &apfs}
    dec.Decode(&omapnode)
    log.Printf("omapnode: %#v", omapnode)


}

func BenchmarkAPFS(b *testing.B) {
    for n := 0; n < b.N; n++ {
        file, _ := os.Open("../../testdata/evidence/filesystem/gpt_apfs.dd")
        dec := ks.NewDecoder(file)
        apfs := Apfs{}
        dec.Decode(&apfs)
        file.Close()
    }
}
