package main

import (
	"fmt"
	"go/format"
	"strings"
	"testing"
	"text/scanner"

	"github.com/iancoleman/strcase"
	"github.com/stretchr/testify/assert"
)

func goify(expr string) string {
	// Create go versions of vars
	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Filename = "example"
	startofExpr := true
	cast := false
	io := false
	ret := ""
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Println(s.TokenText(), tok)
		if startofExpr {
			io = false
		}
		switch tok {
		case scanner.Ident:
			if startofExpr && s.TokenText() != "_index" {
				ret += "k."
			}
			switch s.TokenText() {
			case "_io":
				ret += "IO()"
				io = true
			case "_parent":
				ret += "Parent()"
			case "_root":
				ret += "Root()"
			case "_index":
				ret += "index"
			case "to_i":
				ret += "ToI()"
			case "as":
				cast = true
			default:
				ret += strcase.ToCamel(s.TokenText())
				if io {
					ret += "()"
				}
			}
			startofExpr = s.Peek() != 46
		case 60, 62:
			if s.TokenText() == "<" && cast {
				ret += "("
				startofExpr = false
			} else if s.TokenText() == ">" && cast {
				ret += ")"
				cast = false
				startofExpr = false
			} else {
				ret += s.TokenText()
			}
		case 91, 93:
			startofExpr = false
			ret += s.TokenText()
		case 63:
			parts := strings.SplitN(expr, "?", 2)
			check := goify(parts[0])
			cases := strings.SplitN(parts[1], ":", 2)
			ifvalue := goify(cases[0])
			elsevalue := goify(cases[1])
			return "func()int{if " + check + "{return " + ifvalue + "}else{return " + elsevalue + "}}()"
		default:
			ret += s.TokenText()
		}

	}
	formated, err := format.Source([]byte(ret))
	if err != nil {
		fmt.Println(ret, err)
	}
	return string(formated)
}

func TestGoify(t *testing.T) {

	fmt.Println("EOF", scanner.EOF)
	fmt.Println("Ident", scanner.Ident)
	fmt.Println("Int", scanner.Int)
	fmt.Println("Float", scanner.Float)
	fmt.Println("Char", scanner.Char)
	fmt.Println("String", scanner.String)
	fmt.Println("RawString", scanner.RawString)
	fmt.Println("Comment", scanner.Comment)

	tests := map[string]string{
		"entries_start":                                                   "k.EntriesStart",
		"_root._io":                                                       "k.Root().IO()",
		"_io.size - _root.sector_size":                                    "k.IO().Size() - k.Root().SectorSize",
		"entries_start * _root.sector_size":                               "k.EntriesStart * k.Root().SectorSize",
		"_root.block0.body.as<container_superblock>.block_size":           "k.Root().Block0.Body.(ContainerSuperblock).BlockSize",
		"(xp_desc_base + xp_desc_index) * _root.block_size":               "(k.XpDescBase + k.XpDescIndex) * k.Root().BlockSize",
		"(_parent.node_type & 4) == 0":                                    "(k.Parent().NodeType & 4) == 0",
		"(_parent.level > 0) ? 256 : key_hdr.kind.to_i":                   "func()int{if (k.Parent().Level > 0){return 256}else{return k.KeyHdr.Kind.ToI()}}()",
		"_root.block_size - data_offset - 40 * (_parent.node_type & 1)":   "k.Root().BlockSize - k.DataOffset - 40*(k.Parent().NodeType&1)",
		"key_low + ((key_high & 0x0FFFFFFF) << 32)":                       "k.KeyLow + ((k.KeyHigh & 0x0FFFFFFF) << 32)",
		"xf_header[_index].length + ((8 - xf_header[_index].length) % 8)": "k.XfHeader[index].Length + ((8 - k.XfHeader[index].Length) % 8)",
	}
	for input, result := range tests {
		assert.EqualValues(t, result, goify(input))
	}
}
