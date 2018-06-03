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
	"u1": "uint8", "u2": "uint16", "u4": "uint32", "u8": "uint64",
	"u2le": "uint16", "u4le": "uint32", "u8le": "uint64",
	"u2be": "uint16", "u4be": "uint32", "u8be": "uint64",
	"s1": "int8", "s2": "int16", "s4": "int32", "s8": "int64",
	"s2le": "int16", "s4le": "int32", "s8le": "int64",
	"s2be": "int16", "s4be": "int32", "s8be": "int64",
	"f4": "float32", "f8": "float64",
	"f4le": "float32", "f8le": "float64",
	"f4be": "float32", "f8be": "float64",
	"str": "[]byte", "strz": "[]byte",
	"": "[]byte",
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
				if !cast {
					ret += "Get"
				}
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
