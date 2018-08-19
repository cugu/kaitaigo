package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Result struct {
	Input  string
	GoCode string
	Type   string
}

func TestGoify(t *testing.T) {

	kaitaiTypes = map[string]string{
		"Itoa": "[]byte",
		"len":  "int64",
	}

	tests := []Result{
		// "_root._io":                                                       "k.Root.IO()",
		// "_io.size - _root.sector_size":                                    "k.IO().Size() - k.Root.SectorSize()",
		Result{
			Input:  "true",
			GoCode: "true",
			Type:   "bool",
		},
		Result{
			Input:  "true && false",
			GoCode: "true && false",
			Type:   "bool",
		},
		Result{
			Input:  "not true",
			GoCode: "!true",
			Type:   "bool",
		},
		Result{
			Input:  "[0x20, 0x30, 0x40]",
			GoCode: "[]byte{0x20, 0x30, 0x40}",
			Type:   "int64", // TODO wrong
		},
		Result{
			Input:  "entries_start",
			GoCode: "k.EntriesStart()",
			Type:   "runtime.KSYDecoder",
		},
		Result{
			Input:  "entries_start.to_s",
			GoCode: "strconv.Itoa(int(k.EntriesStart()))",
			Type:   "[]byte",
		},
		Result{
			Input:  "entries_start * _root.sector_size",
			GoCode: "k.EntriesStart() * k.Root.SectorSize()",
			Type:   "int64",
		},
		Result{
			Input:  "_root.block0.body.as<container_superblock>.block_size",
			GoCode: "k.Root.Block0().Body().(ContainerSuperblock).BlockSize()",
			Type:   "runtime.KSYDecoder",
		},
		Result{
			Input:  "(xp_desc_base + xp_desc_index) * _root.block_size",
			GoCode: "(k.XpDescBase() + k.XpDescIndex()) * k.Root.BlockSize()",
			Type:   "int64",
		},
		Result{
			Input:  "(_parent.node_type & 4) == 0",
			GoCode: "(k.Parent().NodeType() & 4) == 0",
			Type:   "bool",
		},
		Result{
			Input:  "(_parent.level > 0) ? 256 : key_hdr.kind.to_i",
			GoCode: "func()int64{if (k.Parent().Level() > 0){return 256}else{return int64(k.KeyHdr().Kind())}}()",
			Type:   "int64",
		},
		Result{
			Input:  "_root.block_size - data_offset - 40 * (_parent.node_type & 1)",
			GoCode: "k.Root.BlockSize() - k.DataOffset() - 40*(k.Parent().NodeType()&1)",
			Type:   "int64",
		},
		Result{
			Input:  "key_low + ((key_high & 0x0FFFFFFF) << 32)",
			GoCode: "k.KeyLow() + ((k.KeyHigh() & 0x0FFFFFFF) << 32)",
			Type:   "int64",
		},
		Result{
			Input:  "xf_header[_index].length + ((8 - xf_header[_index].length) % 8)",
			GoCode: "len(k.XfHeader()[index]) + ((8 - len(k.XfHeader()[index])) % 8)",
			Type:   "int64",
		},
		Result{
			Input:  "\"test\"",
			GoCode: "\"test\"",
			Type:   "[]byte",
		},
		Result{
			Input:  "-2.7",
			GoCode: "-2.7",
			Type:   "float64",
		},
		Result{
			Input:  "2 < 1",
			GoCode: "2 < 1",
			Type:   "bool",
		},
		Result{
			Input:  "2 + 1",
			GoCode: "2 + 1",
			Type:   "int64",
		},
		Result{
			Input:  "-9837 % 13",
			GoCode: "-9837 % 13",
			Type:   "int64",
		},
		Result{
			Input:  "-9837",
			GoCode: "-9837",
			Type:   "int64",
		},
		Result{
			Input:  "\"test\" + \"test\"",
			GoCode: "\"test\" + \"test\"",
			Type:   "[]byte",
		},
		Result{
			Input:  "str1.length",
			GoCode: "len(k.Str1())",
			Type:   "int64",
		},
	}

	for _, result := range tests {
		assert.EqualValues(t, result.GoCode, goExpr(result.Input, ""))
		ty := getType(result.Input)
		// fmt.Println("cmp", result.Type, ty, goExpr(result.Input, ""))
		assert.EqualValues(t, result.Type, ty)
	}
}
