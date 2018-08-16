package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoify(t *testing.T) {

	tests := map[string]string{
		// "_root._io":                                                       "k.Root.IO()",
		// "_io.size - _root.sector_size":                                    "k.IO().Size() - k.Root.SectorSize()",
		"entries_start":                                                   "k.EntriesStart()",
		"entries_start * _root.sector_size":                               "k.EntriesStart() * k.Root.SectorSize()",
		"_root.block0.body.as<container_superblock>.block_size":           "k.Root.Block0().Body().(ContainerSuperblock).BlockSize()",
		"(xp_desc_base + xp_desc_index) * _root.block_size":               "(k.XpDescBase() + k.XpDescIndex()) * k.Root.BlockSize()",
		"(_parent.node_type & 4) == 0":                                    "(k.Parent().NodeType() & 4) == 0",
		"(_parent.level > 0) ? 256 : key_hdr.kind.to_i":                   "func()int64{if (k.Parent().Level() > 0){return 256}else{return int64(k.KeyHdr().Kind())}}()",
		"_root.block_size - data_offset - 40 * (_parent.node_type & 1)":   "k.Root.BlockSize() - k.DataOffset() - 40*(k.Parent().NodeType()&1)",
		"key_low + ((key_high & 0x0FFFFFFF) << 32)":                       "k.KeyLow() + ((k.KeyHigh() & 0x0FFFFFFF) << 32)",
		"xf_header[_index].length + ((8 - xf_header[_index].length) % 8)": "len(k.XfHeader()[index]) + ((8 - len(k.XfHeader()[index])) % 8)",
	}

	for input, result := range tests {
		assert.EqualValues(t, result, goify(input, ""))
	}
}
