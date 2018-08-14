package apfs

import (
    "os"
    "testing"
    "io"
    "log"

    "github.com/stretchr/testify/assert"
)



func TestAPFS(t *testing.T) {
    file, err := os.Open("../../testdata/evidence/filesystem/gpt_apfs.dd")
    defer file.Close()

    if err != nil {
        t.Fatal(err)
    }

    filesystem := io.NewSectionReader(file, 40 * 512, 39024 * 512)

    apfs := Apfs{}
    err = apfs.Decode(filesystem)
    if err != nil {
        t.Fatal(err)
    }

    p0 := apfs.Block0()
    body := p0.Body()
    containerSuperblock := body.(*ContainerSuperblock)
    // log.Printf("containerSuperblock: %#v", containerSuperblock)
    blocksize := containerSuperblock.BlockSize()
    assert.EqualValues(t, 4096, blocksize)

    assert.EqualValues(t, 0x949, containerSuperblock.OmapOid())
    filesystem.Seek(int64(containerSuperblock.OmapOid()) * int64(blocksize), io.SeekStart)
    omap := Btree{}
    omap.DecodeAncestors(apfs.Dec, &apfs, &apfs)
    log.Printf("omap: %#v", omap)

    filesystem.Seek(int64(omap.TreeRoot()) * int64(blocksize), io.SeekStart)
    omapnode := Node{}
    omapnode.DecodeAncestors(apfs.Dec, &apfs, &apfs)
    log.Printf("omapnode: %#v", omapnode)


}

func BenchmarkAPFS(b *testing.B) {
    for n := 0; n < b.N; n++ {
        file, _ := os.Open("../../testdata/evidence/filesystem/gpt_apfs.dd")
        apfs := Apfs{}
        err := apfs.Decode(file)
        if err != nil {
            b.Fatal(err)
        }
        file.Close()
    }
}
