package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/iancoleman/strcase"
)

var typeMapping = map[string]string{
	"u1":   "runtime.Uint8",
	"u2":   "runtime.Uint16",
	"u4":   "runtime.Uint32",
	"u8":   "runtime.Uint64",
	"u2le": "runtime.Uint16Le",
	"u4le": "runtime.Uint32Le",
	"u8le": "runtime.Uint64Le",
	"u2be": "runtime.Uint16Be",
	"u4be": "runtime.Uint32Be",
	"u8be": "runtime.Uint64Be",
	"s1":   "runtime.Int8",
	"s2":   "runtime.Int16",
	"s4":   "runtime.Int32",
	"s8":   "runtime.Int64",
	"s2le": "runtime.Int16Le",
	"s4le": "runtime.Int32Le",
	"s8le": "runtime.Int64Le",
	"s2be": "runtime.Int16Be",
	"s4be": "runtime.Int32Be",
	"s8be": "runtime.Int64Be",
	"f4":   "runtime.Float32",
	"f8":   "runtime.Float64",
	"f4le": "runtime.Float32Le",
	"f8le": "runtime.Float64Le",
	"f4be": "runtime.Float32Be",
	"f8be": "runtime.Float64Be",
	"str":  "runtime.Bytes",
	"strz": "runtime.BytesZ",
	"":     "runtime.Bytes",
}

func bitString(s string) string {
	num, err := strconv.ParseInt(s[2:], 2, 64)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d", num)
}

func isInt(expr string) bool {
	return !strings.Contains(goExpr(expr, ""), "k.")
}

func getExprType(expr ast.Expr) (s string, r bool) {
	// fmt.Printf("%#v\n", expr)
	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BasicLit:
			s = expressionTypes[x.Kind]
			return false
		case *ast.Ident:
			s = getKaitaiType(x.Name)
			return false
		case *ast.UnaryExpr:
			s = "runtime.Int64"
			return false
		case *ast.BinaryExpr:
			s = expressionTypes[x.Op]
			return x.Op == token.ADD || x.Op == token.ADD_ASSIGN
		case *ast.CallExpr:
			s, r = getExprType(x.Fun)
			return r
		case *ast.SelectorExpr:
			s, r = getExprType(x.Sel)
			return r
		case *ast.FuncType:
			s, r = getExprType(x.Results.List[0].Type)
			return r
		default:
			return true
		}
	})
	return
}

func getType(expr string) (t string) {
	var re = regexp.MustCompile(`\*k.*\(\)`)
	goExpr := re.ReplaceAllString(goExpr(expr, ""), `"x"`)

	// fmt.Println()
	// fmt.Println(goExpr)

	exprx, _ := parser.ParseExpr(goExpr)
	var s string
	if exprx != nil {
		s, _ = getExprType(exprx)
	}
	switch s {
	case "int":
		return "runtime.Int64"
	case "string":
		return "runtime.Bytes"
	case "bool":
		return "bool"
	case "":
		return "runtime.Int64"
	default:
		return s
	}
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
	return tok == scanner.Ident || tok == '.' || tok == '[' || tok == ']' || tok == '"' || tok == '_'
}

func goExprIdent(expr, casttype, current_attr string) string {
	ret := "k."
	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Filename = "example"
	cast := false
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch s.TokenText() {
		case "_":
			ret = ret[:len(ret)-1] + "." + current_attr
		case "\"":
			// fmt.Println("....")
			ret += "\""
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
			ret += "Parent()"
		case "_root":
			ret += "Root"
		case "_index":
			ret += "index"
		case "to_i":
			ret = "int64(" + ret[:len(ret)-1] + ")"
		case "to_s":
			ret = "strconv.Itoa(int(" + ret[:len(ret)-1] + "))"
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

func goExpr(expr, casttype string) string {
	return goExprAttr(expr, casttype, "")
}

func goExprAttr(expr, casttype, current_attr string) string {
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
			ret += " " + goExprIdent(identifier, casttype, current_attr)
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
			check := goExpr(parts[0], "")
			cases := strings.SplitN(parts[1], ":", 2)
			ifvalue := goExpr(cases[0], "")
			elsevalue := goExpr(cases[1], "")
			return fmt.Sprintf("func()int64{if %s{return %s}else{return %s}}()", check, ifvalue, elsevalue)
		default:
			ret += s.TokenText()
		}
	}

	// handle identifier chain
	if identifier != "" {
		ret += " " + goExprIdent(identifier, casttype, current_attr)
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
