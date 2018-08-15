package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

type TypeSwitch struct {
	SwitchOn string          `yaml:"switch-on,omitempty"`
	Cases    map[string]Type `yaml:"cases,omitempty"`
}

type Type struct {
	Type       string
	TypeSwitch TypeSwitch
}

func (y *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&y.Type)
	if err != nil {
		err = unmarshal(&y.TypeSwitch)
		return err
	}
	return nil
}

func (y *Type) String() string {
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
	Content []interface{}
}

func (y *Contents) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&y.Content)
}

func (y *Contents) Len() int {
	if len(y.Content) == 0 {
		return 0
	}
	switch v := y.Content[0].(type) {
	case string:
		return len(v)
	default:
		return len(y.Content)
	}
}

type Attribute struct {
	Category   string   `-`
	ID         string   `yaml:"id,omitempty"`
	Type       Type     `yaml:"type"`
	Size       string   `yaml:"size,omitempty"`
	Doc        string   `yaml:"doc,omitempty"`
	Repeat     string   `yaml:"repeat,omitempty"`
	RepeatExpr string   `yaml:"repeat-expr,omitempty"`
	Contents   Contents `yaml:"contents,omitempty"`
	Value      string   `yaml:"value,omitempty"`
	Pos        string   `yaml:"pos,omitempty"`
	Whence     string   `yaml:"whence,omitempty"`
}

func (k *Attribute) dataType() string {
	dataType := k.Type.String()
	if dataType == "runtime.Bytes" { // || dataType == "string" {
		if k.Value != "" {
			if isNumericExpr(k.Value) {
				dataType = "runtime.Int64"
			} else {
				dataType = "runtime.Bytes"
			}
		} else if k.Size != "" {
			k.Repeat = "yes"
			k.RepeatExpr = strings.Replace(k.Size, "%", "%%", -1)
			return "*runtime.ByteArray"
		} else if k.Contents.Len() != 0 {
			k.Repeat = "yes"
			k.RepeatExpr = fmt.Sprintf("%d", k.Contents.Len())
			return "*runtime.ByteArray"
		}
	}

	if k.Repeat != "" && k.RepeatExpr != "" {
		if isInt(k.RepeatExpr) {
			dataType = "[" + goify(k.RepeatExpr, "") + "]" + dataType
		} else {
			dataType = "[]" + dataType
		}
	} else if dataType != "runtime.KSYDecoder" {
		dataType = "*" + dataType

	}
	return dataType
}

func (k *Attribute) String() string {
	doc := ""
	if k.Doc != "" {
		doc = " // " + k.Doc
	}

	return strcase.ToLowerCamel(k.ID) + " " + k.dataType() + "`ks:\"" + k.ID + "," + k.Category + "\"`" + doc
}

type Kaitai struct {
	Types     map[string]Kaitai         `yaml:"types,omitempty"`
	Seq       []Attribute               `yaml:"seq,omitempty"`
	Enums     map[string]map[int]string `yaml:"enums,omitempty"`
	Doc       string                    `yaml:"doc,omitempty"`
	Instances map[string]Attribute      `yaml:"instances,omitempty"`
}

func (k *Kaitai) String(typeName string, parent string, root string) string {
	var buffer LineBuffer

	// print doc string
	if k.Doc != "" {
		buffer.WriteLine("// " + strings.Replace(strings.TrimSpace(k.Doc), "\n", "\n// ", -1))
	}

	// print type start
	buffer.WriteLine("type " + typeName + " struct{")
	buffer.WriteLine("Dec *runtime.Decoder")
	buffer.WriteLine("Start int64")
	buffer.WriteLine("Parent *" + parent)
	buffer.WriteLine("Root *" + root)
	buffer.WriteLine("")

	// print attributes and instances
	for _, attribute := range k.Seq {
		attribute.Category = "attribute"
		buffer.WriteLine(attribute.String())
	}

	for name, instance := range k.Instances {
		instance.Category = "instance"
		instance.ID = name
		buffer.WriteLine(instance.String())
	}

	// print type end
	buffer.WriteLine("}")

	// decode function
	buffer.WriteLine("func (k *" + typeName + ") Decode(reader io.ReadSeeker) (err error) {")
	buffer.WriteLine("return k.DecodeAncestors(&runtime.Decoder{reader, binary.LittleEndian, nil}, k, k)")
	buffer.WriteLine("}")

	// decode pos function
	buffer.WriteLine("func (k *" + typeName + ") DecodePos(dec *runtime.Decoder, offset int64, whence int, parent interface{}, root interface{}) (err error) {")
	buffer.WriteLine("if dec.Err != nil { return dec.Err }")
	buffer.WriteLine("_, dec.Err = dec.Seek(offset, whence)")
	buffer.WriteLine("return k.DecodeAncestors(dec, parent, root)")
	buffer.WriteLine("}")

	// decode ancestors function
	buffer.WriteLine("func (k *" + typeName + ") DecodeAncestors(dec *runtime.Decoder, parent interface{}, root interface{}) (err error) {")
	buffer.WriteLine("if dec.Err != nil { return dec.Err }")
	buffer.WriteLine("k.Parent = parent.(*" + parent + ")")
	buffer.WriteLine("k.Root = root.(*" + root + ")")
	buffer.WriteLine("k.Dec = dec")
	for _, attribute := range k.Seq {
		laName := strcase.ToLowerCamel(attribute.ID)
		aName := strcase.ToCamel(attribute.ID)
		dataType := attribute.dataType()
		if attribute.Repeat != "" && attribute.RepeatExpr != "" && dataType != "*runtime.ByteArray" {
			if strings.HasPrefix(dataType, "[]") {
				buffer.WriteLine("k." + laName + " = make(" + dataType + ", " + goify(attribute.RepeatExpr, "") + ")")
			}
			if isNumericExpr(attribute.RepeatExpr) {
				buffer.WriteLine("for i := 0; i < " + goify(attribute.RepeatExpr, "") + "; i += 1 {")
			} else {
				buffer.WriteLine("for i := 0; i < int(" + goify(attribute.RepeatExpr, "") + "); i += 1 {")
			}
			buffer.WriteLine("k." + laName + "[i].DecodeAncestors(k.Dec, k, k.Root)")
			buffer.WriteLine("}")
		} else if dataType == "runtime.KSYDecoder" {
			buffer.WriteLine("switch " + goify(attribute.Type.TypeSwitch.SwitchOn, "int64") + " {")
			for casevalue, casetype := range attribute.Type.TypeSwitch.Cases {
				buffer.WriteLine("case " + goenum(casevalue, "int64") + ":")
				buffer.WriteLine("so := " + casetype.String() + "{}")
				buffer.WriteLine("so.DecodeAncestors(k.Dec, k, k.Root)")
				buffer.WriteLine("k." + laName + " = &so")
			}
			buffer.WriteLine("}")
		} else {
			if dataType == "*runtime.ByteArray" {
				buffer.WriteLine("k." + laName + " = &runtime.ByteArray{Size: " + goify(attribute.RepeatExpr, "int64") + "}")
			} else {
				buffer.WriteLine("var tmp" + aName + " " + dataType[1:])
				buffer.WriteLine("k." + laName + " = &tmp" + aName + "")

			}
			buffer.WriteLine("k." + laName + ".DecodeAncestors(k.Dec, k, k.Root)")
		}

	}
	buffer.WriteLine("return dec.Err")
	buffer.WriteLine("}")

	// create getter
	for _, attribute := range k.Seq {
		aName := strcase.ToCamel(attribute.ID)
		laName := strcase.ToLowerCamel(attribute.ID)
		dataType := attribute.dataType()

		buffer.WriteLine("func (k *" + typeName + ") " + aName + "() (" + dataType + ") {")

		buffer.WriteLine("return " + "" + "k." + laName)
		buffer.WriteLine("}")
	}

	// create instance getter
	for name, instance := range k.Instances {
		iName := strcase.ToCamel(name)
		liName := strcase.ToLowerCamel(name)
		dataType := instance.dataType()

		buffer.WriteLine("func (k *" + typeName + ") " + iName + "() (" + dataType + ") {")

		if dataType == "runtime.KSYDecoder" {
			buffer.WriteLine("switch " + goify(instance.Type.TypeSwitch.SwitchOn, "int64") + " {")
			for casevalue, casetype := range instance.Type.TypeSwitch.Cases {
				buffer.WriteLine("case " + goenum(casevalue, "int64") + ":")
				buffer.WriteLine("so := " + casetype.String() + "{}")
				buffer.WriteLine("so.DecodeAncestors(k.Dec, k, k.Root)")
				buffer.WriteLine("k." + liName + " = &so")
			}
			buffer.WriteLine("}")
		}

		buffer.WriteLine("if k." + liName + " == nil{")

		if instance.Pos == "" {
			buffer.WriteLine("tmp := " + dataType[1:] + "(" + goify(instance.Value, "") + ")")
			buffer.WriteLine("k." + liName + " = &tmp")
		} else {

			if strings.HasPrefix(dataType, "*") {
				buffer.WriteLine("var tmp" + iName + " " + dataType[1:])
				buffer.WriteLine("k." + liName + " = &tmp" + iName)
			} else {
				buffer.WriteLine("var tmp" + iName + " " + dataType)
				buffer.WriteLine("k." + liName + " = tmp" + iName)
			}

			buffer.WriteLine("_, k.Dec.Err = k.Dec.Seek(k.Start, io.SeekStart)")
			whence := "io.SeekCurrent"
			switch instance.Whence {
			case "seek_set":
				whence = "io.SeekStart"
			case "seek_end":
				whence = "io.SeekEnd"
			}

			if instance.Repeat != "" && instance.RepeatExpr != "" {
				dataType := instance.dataType()
				buffer.WriteLine("k." + liName + " = make(" + dataType + ", " + goify(instance.RepeatExpr, "") + ")") // TODO: needed?
				buffer.WriteLine("for i := 0; i < " + goify(instance.RepeatExpr, "int") + "; i++ {")
				buffer.WriteLine("k." + liName + "[i].DecodePos(k.Dec, " + goify(instance.Pos, "int64") + ", " + whence + ", k, k.Root)")
				buffer.WriteLine("}")
			} else {
				buffer.WriteLine("k." + liName + ".DecodePos(k.Dec, " + goify(instance.Pos, "int64") + ", " + whence + ", k, k.Root)")
			}

		}
		buffer.WriteLine("}")
		buffer.WriteLine("return k." + liName)
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
			buffer.WriteLine(strcase.ToCamel(value) + " int64")
		}
		buffer.WriteLine("}{")
		for x, value := range values {
			buffer.WriteLine(strcase.ToCamel(value) + ": " + strconv.Itoa(x) + ",")
		}
		buffer.WriteLine("}")
	}

	return buffer.String()
}
