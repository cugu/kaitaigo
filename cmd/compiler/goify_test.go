package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoify(t *testing.T) {

	tests := map[string]string{
		"entries_start": "k.GetEntriesStart()",
		// "_root._io":                                                       "k.Root.IO()",
		// "_io.size - _root.sector_size":                                    "k.IO().Size() - k.Root.SectorSize()",
		"entries_start * _root.sector_size":                               "k.GetEntriesStart() * k.Root.GetSectorSize()",
		"_root.block0.body.as<container_superblock>.block_size":           "k.Root.GetBlock0().GetBody().(ContainerSuperblock).GetBlockSize()",
		"(xp_desc_base + xp_desc_index) * _root.block_size":               "(k.GetXpDescBase() + k.GetXpDescIndex()) * k.Root.GetBlockSize()",
		"(_parent.node_type & 4) == 0":                                    "(k.Parent.GetNodeType() & 4) == 0",
		"(_parent.level > 0) ? 256 : key_hdr.kind.to_i":                   "func()int{if (k.Parent.GetLevel() > 0){return 256}else{return k.GetKeyHdr().GetKind().ToI()}}()",
		"_root.block_size - data_offset - 40 * (_parent.node_type & 1)":   "k.Root.GetBlockSize() - k.GetDataOffset() - 40*(k.Parent.GetNodeType()&1)",
		"key_low + ((key_high & 0x0FFFFFFF) << 32)":                       "k.GetKeyLow() + ((k.GetKeyHigh() & 0x0FFFFFFF) << 32)",
		"xf_header[_index].length + ((8 - xf_header[_index].length) % 8)": "k.GetXfHeader()[index].GetLength() + ((8 - k.GetXfHeader()[index].GetLength()) % 8)",
	}
	for input, result := range tests {
		assert.EqualValues(t, result, goify(input, ""))
	}
}
