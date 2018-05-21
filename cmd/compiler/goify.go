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

func bitString(s string) string {
	num, err := strconv.ParseInt(s[2:], 2, 64)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d", num)
}

func goify(expr string) string {
	re := regexp.MustCompile("0[bB][0-9]+")
	expr = re.ReplaceAllStringFunc(expr, bitString)

	log.Println(expr)
	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Filename = "example"
	startofExpr := true
	cast := false
	io := false
	ret := ""
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		log.Println(s.TokenText(), tok)
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
