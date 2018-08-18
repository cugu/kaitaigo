package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

type TypeSwitch struct {
	SwitchOn string             `yaml:"switch-on,omitempty"`
	Cases    map[string]TypeKey `yaml:"cases,omitempty"`
}

type TypeKey struct {
	Type       string
	TypeSwitch TypeSwitch
	CustomType bool
}

func (y *TypeKey) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&y.Type)
	if err != nil {
		err = unmarshal(&y.TypeSwitch)
		return err
	}
	if _, ok := typeMapping[y.Type]; !ok {
		y.CustomType = true
	}
	return nil
}

func (y *TypeKey) String() string {
	if y.Type != "" {
		if val, ok := typeMapping[y.Type]; ok {
			return val
		}
		return strcase.ToCamel(y.Type)
	} else if y.TypeSwitch.SwitchOn != "" {
		return "runtime.KSYDecoder"
	}
	return "runtime.Bytes"
}

type Contents struct {
	ContentString string
	ContentArray  []interface{}
	TypeSwitch    TypeSwitch
}

func (y *Contents) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&y.ContentString)
	if err != nil {
		err := unmarshal(&y.ContentArray)
		if err != nil {
			err = unmarshal(&y.TypeSwitch)
			return err
		}
		return err
	}
	return nil
}

func (y *Contents) Len() int {
	if len(y.ContentString) != 0 {
		return len(y.ContentString)
	}
	if len(y.ContentArray) == 0 {
		return 0
	}
	switch v := y.ContentArray[0].(type) {
	case string:
		return len(v)
	default:
		return len(y.ContentArray)
	}
}

type Attribute struct {
	Category   string   `-`
	ID         string   `yaml:"id,omitempty"`
	Type       TypeKey  `yaml:"type"`
	Size       string   `yaml:"size,omitempty"`
	Doc        string   `yaml:"doc,omitempty"`
	Repeat     string   `yaml:"repeat,omitempty"`
	RepeatExpr string   `yaml:"repeat-expr,omitempty"`
	Contents   Contents `yaml:"contents,omitempty"`
	Value      string   `yaml:"value,omitempty"`
	Pos        string   `yaml:"pos,omitempty"`
	Whence     string   `yaml:"whence,omitempty"`
	Enum       string   `yaml:"enum,omitempty"`
}

func (k *Attribute) Name() string {
	return strcase.ToLowerCamel(k.ID)
}

func (k *Attribute) DataType() string {
	dataType := k.Type.String()
	if dataType == "runtime.Bytes" { // || dataType == "runtime.String" {
		if k.Value != "" {
			dataType = getType(k.Value)
		} else if k.Size != "" {
			k.Repeat = "yes"
			k.RepeatExpr = strings.Replace(k.Size, "%", "%%", -1)
			return "runtime.ByteSlice"
		} else if k.Contents.Len() != 0 {
			k.Repeat = "yes"
			k.RepeatExpr = fmt.Sprintf("%d", k.Contents.Len())
			return "runtime.ByteSlice"
		}
	}

	if k.Repeat != "" {
		if !isInt(k.RepeatExpr) || k.Repeat == "eos" {
			dataType = "[]" + dataType
		} else {
			dataType = "[" + goExpr(k.RepeatExpr, "") + "]" + dataType
		}
	} else if k.Type.CustomType {
		dataType = "*" + dataType
	}
	return dataType
}

func (k *Attribute) String() string {
	doc := ""
	if k.Doc != "" {
		doc = " // " + k.Doc
	}

	return k.Name() + " " + k.DataType() + "`ks:\"" + k.ID + "," + k.Category + "\"`" + doc
}

type Type struct {
	Types     map[string]Type           `yaml:"types,omitempty"`
	Seq       []Attribute               `yaml:"seq,omitempty"`
	Enums     map[string]map[int]string `yaml:"enums,omitempty"`
	Doc       string                    `yaml:"doc,omitempty"`
	Instances map[string]Attribute      `yaml:"instances,omitempty"`
}

func (k *Type) InitAttr(attr Attribute) string {
	var buffer LineBuffer

	if attr.Value != "" {
		// value instance
		if attr.DataType() == "runtime.KSYDecoder" || strings.HasPrefix(attr.DataType(), "*") {
			buffer.WriteLine("k." + attr.Name() + " = " + goExpr(attr.Value, ""))
		} else {
			buffer.WriteLine("k." + attr.Name() + " = " + attr.DataType() + "(" + goExpr(attr.Value, "") + ")")
		}
		return buffer.String()
	}

	if attr.Pos != "" {
		buffer.WriteLine("_, decoder.Err = decoder.Seek(k.Start, io.SeekStart)")
		whence := "io.SeekCurrent"
		switch attr.Whence {
		case "seek_set":
			whence = "io.SeekStart"
		case "seek_end":
			whence = "io.SeekEnd"
		}
		buffer.WriteLine("_, decoder.Err = decoder.Seek(" + goExpr(attr.Pos, "int64") + ", " + whence + ")")
	}

	switch {
	case attr.DataType() == "runtime.ByteSlice":
		// byteslice
		buffer.WriteLine("k." + attr.Name() + " = make(runtime.ByteSlice, " + goExpr(attr.RepeatExpr, "int64") + ")")
	case attr.Repeat != "":
		switch attr.Repeat {
		case "expr":
			if attr.RepeatExpr == "" {
				panic("RepeatExpr is missing")
			}
			// array
			if strings.HasPrefix(attr.DataType(), "[]") {
				buffer.WriteLine("k." + attr.Name() + " = make(" + attr.DataType() + ", " + goExpr(attr.RepeatExpr, "") + ")")
			}
			buffer.WriteLine("for i := 0; i < int(" + goExpr(attr.RepeatExpr, "") + "); i += 1 {")
			buffer.WriteLine("k." + attr.Name() + "[i].DecodeAncestors(k, k.Root)")
			buffer.WriteLine("}")
			return buffer.String()
		case "eos":
			buffer.WriteLine("k." + attr.Name() + " = " + attr.DataType() + "{}")
			buffer.WriteLine("for i := 0; true; i++ {")
			buffer.WriteLine("var elem " + attr.DataType()[2:])
			buffer.WriteLine("elem.DecodeAncestors(k, k.Root)")
			buffer.WriteLine("if decoder.Err != nil {")
			buffer.WriteLine("decoder.Err = nil")
			buffer.WriteLine("break")
			buffer.WriteLine("}")
			buffer.WriteLine("k." + attr.Name() + " = append(k." + attr.Name() + ", elem)")
			buffer.WriteLine("}")
			return buffer.String()
		}
	case attr.Type.CustomType:
		// custom struct
		// init variable
		buffer.WriteLine("k." + attr.Name() + " = &" + attr.DataType()[1:] + "{}")
	case attr.Type.TypeSwitch.SwitchOn != "":
		buffer.WriteLine("switch " + goExpr(attr.Type.TypeSwitch.SwitchOn, "int64") + " {")
		for casevalue, casetype := range attr.Type.TypeSwitch.Cases {
			buffer.WriteLine("case " + goenum(casevalue, "int64") + ":")
			buffer.WriteLine("k." + attr.Name() + " = &" + casetype.String() + "{}")
		}
		buffer.WriteLine("}")
	}

	buffer.WriteLine("k." + attr.Name() + ".DecodeAncestors(k, k.Root)")
	return buffer.String()

}

func (k *Type) String(typeName string, parent string, root string) string {
	var buffer LineBuffer

	// print doc string
	if k.Doc != "" {
		buffer.WriteLine("// " + strings.Replace(strings.TrimSpace(k.Doc), "\n", "\n// ", -1))
	}

	// print type start
	buffer.WriteLine("type " + typeName + " struct{")
	buffer.WriteLine("Start int64")
	buffer.WriteLine("parent interface{}")
	buffer.WriteLine("Root *" + root)

	// print attrs and insts
	for _, attr := range k.Seq {
		attr.Category = "attribute"
		buffer.WriteLine(attr.String())
	}

	for name, inst := range k.Instances {
		inst.Category = "instance"
		inst.ID = name
		buffer.WriteLine(inst.String())
		buffer.WriteLine(strcase.ToLowerCamel(inst.ID) + "Set bool")
	}

	// print type end
	buffer.WriteLine("}")

	// decode function
	buffer.WriteLine("func (k *" + typeName + ") Decode(reader io.ReadSeeker) (err error) {")
	buffer.WriteLine("if decoder != nil && decoder.Err != nil { return decoder.Err }")
	buffer.WriteLine("if decoder == nil { decoder = &runtime.Decoder{reader, binary.LittleEndian, nil}; runtime.RTDecoder = decoder }")
	buffer.WriteLine("k.DecodeAncestors(k, k)")
	buffer.WriteLine("return decoder.Err")
	buffer.WriteLine("}")

	// parent function
	buffer.WriteLine("func (k *" + typeName + ") Parent() (*" + parent + ") {")
	buffer.WriteLine("return k.parent.(*" + parent + ")")
	buffer.WriteLine("}")

	// decode ancestors function
	buffer.WriteLine("func (k *" + typeName + ") DecodeAncestors(parent interface{}, root interface{}) () {")
	buffer.WriteLine("if decoder.Err != nil { return }")
	buffer.WriteLine("k.parent = parent")
	buffer.WriteLine("k.Root = root.(*" + root + ")")
	for _, attr := range k.Seq {
		buffer.WriteString(k.InitAttr(attr))
	}
	buffer.WriteLine("return")
	buffer.WriteLine("}")

	// create getter
	for _, attr := range k.Seq {
		buffer.WriteLine("func (k *" + typeName + ") " + strcase.ToCamel(attr.Name()) + "() (" + attr.DataType() + ") {")
		buffer.WriteLine("return " + "" + "k." + attr.Name())
		buffer.WriteLine("}")
	}

	// create inst getter
	for name, inst := range k.Instances {
		inst.ID = name
		buffer.WriteLine("func (k *" + typeName + ") " + strcase.ToCamel(inst.Name()) + "() (" + inst.DataType() + ") {")
		buffer.WriteLine("if !k." + inst.Name() + "Set {")
		buffer.WriteString(k.InitAttr(inst))
		buffer.WriteLine("k." + inst.Name() + "Set = true")
		buffer.WriteLine("}")
		buffer.WriteLine("return k." + inst.Name())
		buffer.WriteLine("}")
	}

	// print subtypes (flattened)
	for name, t := range k.Types {
		typeStr := t.String(strcase.ToCamel(name), getParent(strcase.ToCamel(name)), root)
		buffer.WriteLine(typeStr)
	}

	// print enums
	for enum, values := range k.Enums {
		buffer.WriteLine("var " + strcase.ToCamel(enum) + " = struct {")
		for _, value := range values {
			buffer.WriteLine(strcase.ToCamel(value) + " " + getEnumType(enum))
		}
		buffer.WriteLine("}{")
		for x, value := range values {
			buffer.WriteLine(strcase.ToCamel(value) + ": " + strconv.Itoa(x) + ",")
		}
		buffer.WriteLine("}")
	}

	return buffer.String()
}
