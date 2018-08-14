package main

import (
	"fmt"
	"go/format"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/iancoleman/strcase"
)

var typeMapping = map[string]string{
	"u1": "runtime.Uint8", "u2": "runtime.Uint16", "u4": "runtime.Uint32", "u8": "runtime.Uint64",
	"u2le": "runtime.Uint16", "u4le": "runtime.Uint32", "u8le": "runtime.Uint64",
	"u2be": "runtime.Uint16", "u4be": "runtime.Uint32", "u8be": "runtime.Uint64",
	"s1": "runtime.Int8", "s2": "runtime.Int16", "s4": "runtime.Int32", "s8": "runtime.Int64",
	"s2le": "runtime.Int16", "s4le": "runtime.Int32", "s8le": "runtime.Int64",
	"s2be": "runtime.Int16", "s4be": "runtime.Int32", "s8be": "runtime.Int64",
	"f4": "runtime.Float32", "f8": "runtime.Float64",
	"f4le": "runtime.Float32", "f8le": "runtime.Float64",
	"f4be": "runtime.Float32", "f8be": "runtime.Float64",
	"str": "runtime.Bytes", "strz": "runtime.Bytes",
	"": "runtime.Bytes",
}

func bitString(s string) string {
	num, err := strconv.ParseInt(s[2:], 2, 64)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d", num)
}

func isInt(expr string) bool {
	return !strings.Contains(goify(expr, ""), "k.")
}

func goenum(s string, cast string) string {
	parts := strings.SplitN(s, "::", 2)
	s = strcase.ToCamel(parts[0]) + "." + strcase.ToCamel(parts[1])
	if cast != "" {
		return cast + "(" + s + ")"
	}
	return s
}

func goify(expr string, casttype string) string {

	re := regexp.MustCompile("0[bB][0-9]+")
	expr = re.ReplaceAllStringFunc(expr, bitString)

	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Filename = "example"
	startofExpr := true
	cast := false
	//io := false
	ret := ""
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		//if startofExpr {
		//	io = false
		//}
		switch tok {
		case scanner.Ident:
			if startofExpr && s.TokenText() != "_index" {
				if casttype != "" {
					ret += casttype + "("
				}
				ret += "k."
			}
			switch s.TokenText() {
			// case "_io":
			// 	ret += "IO"
			// 	io = true
			case "_parent":
				ret += "Parent"
			case "_root":
				ret += "Root"
			case "_index":
				ret += "index"
			case "to_i":
				ret += "ToI()"
			case "as":
				cast = true
			default:
				/*if !cast {
					ret += "Get"
				}*/
				ret += strcase.ToCamel(s.TokenText())
				// if io {
				if !cast {
					ret += "()"
				}
			}
			if casttype != "" && s.Peek() != 46 && s.Peek() != scanner.Ident {
				ret += ")"
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
			check := goify(parts[0], "")
			cases := strings.SplitN(parts[1], ":", 2)
			ifvalue := goify(cases[0], "")
			elsevalue := goify(cases[1], "")
			return fmt.Sprintf("func()int{if %s{return %s}else{return %s}}()", check, ifvalue, elsevalue)
		default:
			ret += s.TokenText()
		}

	}
	formated, err := format.Source([]byte(ret))
	if err != nil {
		log.Println(ret, err)
		return ret
	}
	return string(formated)
}
