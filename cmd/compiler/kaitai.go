package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
	yaml "gopkg.in/yaml.v2"
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

type TypeSwitch struct {
	SwitchOn string            `yaml:"switch-on,omitempty"`
	Cases    map[string]string `yaml:"cases,omitempty"`
}

type Type struct {
	Type       string
	TypeSwitch TypeSwitch
}

func (s *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&s.Type)
	if err != nil {
		err = unmarshal(&s.TypeSwitch)
		return err
	}
	return nil
}

func (t *Type) String() (string, error) {
	if t.Type != "" {
		if val, ok := typeMapping[t.Type]; ok {
			return val, nil
		}

		return strcase.ToCamel(t.Type), nil
	}
	if t.TypeSwitch.SwitchOn != "" {
		return "interface{}", nil
	}
	return "[]byte", nil
}

type Instance struct {
	Value      string `yaml:"value,omitempty"`
	Pos        string `yaml:"pos,omitempty"`
	Whence     string `yaml:"whence,omitempty"`
	Type       Type   `yaml:"type,omitempty"`
	Repeat     string `yaml:"repeat,omitempty"`
	RepeatExpr string `yaml:"repeat-expr,omitempty"`
	Doc        string `yaml:"doc,omitempty"`
}

func (k *Instance) dataType() (string, error) {
	dataType, err := k.Type.String()
	if err != nil {
		return "", err
	}
	if dataType == "[]byte" && k.Value != "" {
		dataType = "int64"
	}

	if k.Repeat != "" {
		if k.RepeatExpr != "" {
			if isInt(k.RepeatExpr) {
				dataType = "[" + goify(k.RepeatExpr, "") + "]" + dataType
			} else {
				dataType = "[]" + dataType
			}
		}
	}

	return dataType, nil
}

func (k *Instance) String() (string, error) {
	doc := ""
	if k.Doc != "" {
		doc = " // " + k.Doc
	}

	dataType, err := k.dataType()
	if err != nil {
		return "", err
	}

	return "%[1]s " + dataType + " `ks:\"%[2]s,instance\"`" + doc + "\n", nil
}

type Contents struct {
	Content []interface{}
}

func (s *Contents) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&s.Content)
}

func (k *Contents) Len() int {
	if len(k.Content) == 0 {
		return 0
	}
	switch v := k.Content[0].(type) {
	case string:
		return len(v)
	default:
		return len(k.Content)
	}

}

type Attribute struct {
	ID         string   `yaml:"id,omitempty"`
	Type       Type     `yaml:"type"`
	Size       string   `yaml:"size,omitempty"`
	Doc        string   `yaml:"doc,omitempty"`
	Repeat     string   `yaml:"repeat,omitempty"`
	RepeatExpr string   `yaml:"repeat-expr,omitempty"`
	Contents   Contents `yaml:"contents,omitempty"`
}

func (k *Attribute) dataType() (string, error) {
	dataType, err := k.Type.String()
	if err != nil {
		return k.ID, err
	}

	if dataType == "[]byte" { // || dataType == "string" {
		if k.Size != "" {
			_, err := strconv.ParseInt(k.Size, 0, 0)
			if err != nil {
				return dataType, nil
			}
			size := strings.Replace(k.Size, "%", "%%", -1)
			dataType = strings.Replace(dataType, "[]", "["+goify(size, "")+"]", 1)
		} else {
			if k.Contents.Len() != 0 {
				dataType = strings.Replace(dataType, "[]", fmt.Sprintf("[%d]", k.Contents.Len()), 1)
			}
		}
	}

	if k.Repeat != "" {
		if k.RepeatExpr != "" {
			if isInt(k.RepeatExpr) {
				dataType = "[" + goify(k.RepeatExpr, "") + "]" + dataType
			} else {
				dataType = "[]" + dataType
			}
		}
	}

	return dataType, nil
}

func (k *Attribute) String() (string, error) {
	doc := ""
	if k.Doc != "" {
		doc = " // " + k.Doc
	}

	dataType, err := k.dataType()
	if err != nil {
		return "", err
	}

	return strcase.ToCamel(k.ID) + " " + dataType + "`ks:\"" + k.ID + ",attribute\"`" + doc + "\n", nil
}

var allTypes map[string]Kaitai

type Kaitai struct {
	Types     map[string]Kaitai         `yaml:"types,omitempty"`
	Seq       []Attribute               `yaml:"seq,omitempty"`
	Enums     map[string]map[int]string `yaml:"enums,omitempty"`
	Doc       string                    `yaml:"doc,omitempty"`
	Instances map[string]Instance       `yaml:"instances,omitempty"`
}

func (k *Kaitai) getParent(sname string) string {
	for ksName, ks := range allTypes {
		for _, attribute := range ks.Seq {
			if attribute.Type.Type == sname { // TODO: add TypeSwitch support
				return strcase.ToCamel(ksName)
			}
		}
		for _, instance := range ks.Instances {
			if instance.Type.Type == sname { // TODO: add TypeSwitch support
				return strcase.ToCamel(ksName)
			}
		}
	}
	return "interface{}"
}

func (k *Kaitai) setupMap(sname string) {
	allTypes[sname] = *k
	for name, t := range k.Types {
		t.setupMap(name)
	}
}

func (k *Kaitai) String(sname string, parent string, root string) (string, error) {
	s := ""

	// print doc string
	if k.Doc != "" {
		s += "// " + strings.Replace(strings.TrimSpace(k.Doc), "\n", "\n// ", -1) + "\n"
	}

	// print type start
	s += "type " + sname + " struct{\n"
	s += "\tDec *runtime.Decoder \n"
	s += "\tStart int64 \n"
	s += "\tParent *" + parent + "\n"
	s += "\tRoot *" + root + "\n\n"

	// print attribute
	for _, attribute := range k.Seq {
		attrStr, err := attribute.String()
		if err != nil {
			log.Printf("Error in %s\n", attrStr)
			return "", err
		}
		s += "\t" + attrStr
	}

	for name, instance := range k.Instances {
		attrStr, err := instance.String()
		if err != nil {
			return "", err
		}
		s += "\t" + fmt.Sprintf(attrStr, strcase.ToCamel(name), name)
	}

	// print type end
	s += "}\n\n"

	// create getter
	for _, attribute := range k.Seq {
		dataType, err := attribute.dataType()
		if err != nil {
			return "", err
		}
		s += "func (k *" + sname + ") Get" + strcase.ToCamel(attribute.ID) + "() (" + dataType + ") {\n"
		s += "\treturn k." + strcase.ToCamel(attribute.ID)
		s += "}\n\n"
	}
	for name, instance := range k.Instances {
		dataType, err := instance.dataType()
		if err != nil {
			return "", err
		}
		s += "func (k *" + sname + ") Get" + strcase.ToCamel(name) + "() (" + dataType + ") {\n"

		s += "\tif runtime.IsNull(k." + strcase.ToCamel(name) + "){\n"
		if instance.Pos == "" {
			s += "\t\tk." + strcase.ToCamel(name) + "=" + dataType + "(" + goify(instance.Value, "") + ")\n"
		} else {
			s += "\t\t_, k.Dec.Err = k.Dec.Seek(k.Start, io.SeekStart)\n"
			whence := "io.SeekCurrent"
			switch instance.Whence {
			case "seek_set":
				whence = "io.SeekStart"
			case "seek_end":
				whence = "io.SeekEnd"
			}

			repeat := (instance.Repeat != "" && instance.RepeatExpr != "")
			if repeat {
				dataType, err := instance.dataType()
				if err != nil {
					return "", err
				}
				s += "\t\tk." + strcase.ToCamel(name) + " = make(" + dataType + ", " + goify(instance.RepeatExpr, "") + ")\n"
				s += "\t\tfor i := 0; i < " + goify(instance.RepeatExpr, "int") + "; i++ {\n"
				s += "\t\tk.Dec.DecodePos(&k." + strcase.ToCamel(name) + "[i], " + goify(instance.Pos, "int64") + ", " + whence + ", k, k.Root)\n"
				s += "\t\t}\n"
			} else {
				s += "\t\tk.Dec.DecodePos(&k." + strcase.ToCamel(name) + ", " + goify(instance.Pos, "int64") + ", " + whence + ", k, k.Root)\n"
			}

		}
		s += "\t}\n"
		s += "\treturn k." + strcase.ToCamel(name)
		s += "}\n\n"
	}

	// print subtypes (flattened)
	for name, t := range k.Types {
		typeStr, err := t.String(strcase.ToCamel(name), k.getParent(name), root)
		if err != nil {
			return "", err
		}
		s += typeStr
	}

	for enum, values := range k.Enums {
		s += "var " + strcase.ToCamel(enum) + " = struct {\n"
		for _, value := range values {
			s += "\t" + strcase.ToCamel(value) + " int64\n"
		}
		s += "}{\n"
		for x, value := range values {
			s += "\t" + strcase.ToCamel(value) + ": " + strconv.Itoa(x) + ",\n"
		}
		s += "}\n\n"
	}

	return s, nil
}

func YAMLUnmarshal(name string, source []byte, m interface{}, path string) error {
	err := yaml.Unmarshal(source, m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(
		path+"."+name+".unmarshal",
		[]byte(fmt.Sprintf("%s%# v\n", "// file generated at "+time.Now().UTC().Format(time.RFC3339)+"\n", pretty.Formatter(m))),
		0644,
	)
}

func createGofile(filepath string, pckg string) error {
	filename := path.Base(filepath)
	baseStructSnake := strings.Replace(filename, ".ksy", "", 1)
	baseStruct := strcase.ToCamel(baseStructSnake)

	// setup logging
	logfile, err := os.Create(path.Join(pckg, filename+".log"))
	if err != nil {
		return errors.Wrap(err, "create logfile")
	}
	defer func() {
		logfile.Sync()
		logfile.Close()
	}()

	log.SetOutput(io.MultiWriter(os.Stderr, logfile))

	log.Println("generate", filepath)

	// read source
	source, err := ioutil.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "read source")
	}

	// parse generic
	m := make(map[interface{}]interface{})
	err = YAMLUnmarshal("generic", source, &m, path.Join(pckg, filename))
	if err != nil {
		return errors.Wrap(err, "parse generic yaml")
	}

	// parse kaitai
	kaitai := Kaitai{}
	err = YAMLUnmarshal("kaitai", source, &kaitai, path.Join(pckg, filename))
	if err != nil {
		return errors.Wrap(err, "parse kaitai yaml")
	}

	allTypes = map[string]Kaitai{}
	kaitai.setupMap(baseStruct)

	// write go code
	goCode, err := kaitai.String(baseStruct, baseStruct, baseStruct)
	if err != nil {
		return errors.Wrap(err, "kaitai code gen")
	}
	parts := strings.Split(pckg, "/")
	lastpart := parts[len(parts)-1]
	header := "// file generated at " + time.Now().UTC().Format(time.RFC3339) + "\n"
	header += "package " + lastpart + "\n"
	header += "import (\n"
	for _, pkg := range []string{"fmt", "io", "os", "log", "gitlab.com/cugu/kaitai.go/runtime"} {
		header += "\"" + pkg + "\"\n"
	}
	header += ")\n"

	formated, err := imports.Process("", []byte(header+goCode), nil)
	if err != nil {
		log.Print("Format error", err)
		formated = []byte(header + goCode)
	}
	err = ioutil.WriteFile(path.Join(pckg, filename+".go"), formated, 0644)
	if err != nil {
		return errors.Wrap(err, "create go file")
	}
	return nil

}

func printErrors(filename string) {
	err := createGofile(filename, filepath.Dir(filename))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()
	for _, filename := range flag.Args() {
		if strings.HasSuffix(filename, ".ksy") {
			printErrors(filename)
		} else if strings.HasSuffix(filename, "/...") {
			recPath := strings.Replace(filename, "/...", "", 1)
			filepath.Walk(recPath, func(path string, f os.FileInfo, err error) error {
				if strings.HasSuffix(path, ".ksy") {
					printErrors(path)
				}
				return nil
			})
		}
	}
}
