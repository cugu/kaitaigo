package main

import (
	"fmt"
	"go/format"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/Knetic/govaluate"
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

func getType(expr string) (t string) {
	var re = regexp.MustCompile(`\*k.*\(\)`)
	goExpr := re.ReplaceAllString(goify(expr, ""), `"x"`)

	if goExpr == "\"x\"" {
		// return interface if only *k...()
		return "runtime.KSYDecoder"
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(goExpr, map[string]govaluate.ExpressionFunction{
		"len": func(arguments ...interface{}) (interface{}, error) {
			return 0, nil
		},
	})
	if err != nil {
		return "runtime.Int64"
	}

	if expression != nil {
		result, _ := expression.Evaluate(nil)

		if reflect.TypeOf(result) == reflect.TypeOf("string") {
			return "runtime.Bytes"
		} else if reflect.TypeOf(result) == reflect.TypeOf(float64(0)) {
			return "runtime.Int64"
		} else {
			return "runtime.Int64"
		}
	}

	/* var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Filename = "example"

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if tok == scanner.Char || tok == scanner.String || tok == scanner.RawString || tok == scanner.Comment {
			return "runtime.Bytes"
		}
	} */
	return "runtime.Int64"
}

func goenum(s string, cast string) string {
	// cast
	if strings.HasSuffix(s, ".to_i") {
		s = s[:len(s)-5]
		cast = "int64"
	}

	parts := strings.SplitN(s, "::", 2)
	if len(parts) < 2 {
		return s
	}
	s = strcase.ToCamel(parts[0]) + "." + strcase.ToCamel(parts[1])
	if cast != "" {
		return cast + "(" + s + ")"
	}

	return s
}

func isIdentifierPart(tok rune, casting bool) bool {
	// handle greater and lower than
	if tok == '<' || tok == '>' {
		if !casting {
			return false
		} else {
			return true
		}
	}
	return tok == scanner.Ident || tok == '.' || tok == '[' || tok == ']'
}

func goifyIdent(expr, casttype string) string {
	ret := "*k."
	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Filename = "example"
	cast := false
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch s.TokenText() {
		case ".":
			ret += "."
		case "<":
			ret += "("
		case ">":
			ret += ")"
			cast = false
		case "[", "]":
			ret += s.TokenText()
		case "_parent":
			ret += "Parent"
		case "_root":
			ret += "Root"
		case "_index":
			ret += "index"
		case "to_i":
			ret = "int64(" + ret[:len(ret)-1] + ")"
		case "as":
			cast = true
		case "length":
			if expr == "length" {
				ret += strcase.ToCamel(s.TokenText())
				if !cast {
					ret += "()"
				}
			} else {
				ret = "len(" + ret[:len(ret)-1] + ")"
			}
		default:
			ret += strcase.ToCamel(s.TokenText())
			if !cast {
				ret += "()"
			}
		}
	}

	if casttype != "" {
		return casttype + "(" + ret + ")"
	}
	return ret

}

func goify(expr, casttype string) string {

	// replace binary values
	re := regexp.MustCompile("0[bB][0-9]+")
	expr = re.ReplaceAllStringFunc(expr, bitString)

	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Filename = "example"
	identifier := ""
	casting := false
	ret := ""
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {

		// handle identifier chain
		if !isIdentifierPart(tok, casting) && identifier != "" {
			ret += " " + goifyIdent(identifier, casttype)
			identifier = ""
		}

		switch {
		case isIdentifierPart(tok, casting):

			identifierPart := s.TokenText()
			// identify casting start
			if identifierPart == "as" {
				casting = true
			}
			// identify casting end
			if tok == '>' {
				casting = false
			}

			identifier += identifierPart
		case tok == '?':
			parts := strings.SplitN(expr, "?", 2)
			check := goify(parts[0], "")
			cases := strings.SplitN(parts[1], ":", 2)
			ifvalue := goify(cases[0], "")
			elsevalue := goify(cases[1], "")
			return fmt.Sprintf("func()int64{if %s{return %s}else{return %s}}()", check, ifvalue, elsevalue)
		default:
			ret += s.TokenText()
		}
	}

	// handle identifier chain
	if identifier != "" {
		ret += " " + goifyIdent(identifier, casttype)
	}

	// remove whitespace and format code
	ret = strings.TrimSpace(ret)
	formated, err := format.Source([]byte(ret))
	if err != nil {
		log.Println(ret, err)
		return ret
	}
	return string(formated)
}
