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
	"str": "string", "strz": "string",
	"": "[]byte",
}

type Type struct {
	Type string
}

func (s *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&s.Type)
	if err != nil {
		s.Type = "runtime.KSYDecoder"
		log.Printf("Type unmarshal error: %s", err)
		return nil
	}
	return nil
}

func (t *Type) String() (string, error) {
	if val, ok := typeMapping[t.Type]; ok {
		return val, nil
	}
	return strcase.ToCamel(t.Type), nil
}

type Instance struct {
	Value string `yaml:"value,omitempty"`
	Pos   string `yaml:"pos,omitempty"`
	Type  Type   `yaml:"type,omitempty"`
	Doc   string `yaml:"doc,omitempty"`
}

func (k *Instance) dataType() (string, error) {
	dataType, err := k.Type.String()
	if err != nil {
		return "", err
	}
	if dataType == "[]byte" && k.Value != "" {
		dataType = "int64"
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

type Attribute struct {
	ID         string `yaml:"id,omitempty"`
	Type       Type   `yaml:"type"`
	Size       string `yaml:"size,omitempty"`
	Doc        string `yaml:"doc,omitempty"`
	Repeat     string `yaml:"repeat,omitempty"`
	RepeatExpr string `yaml:"repeat-expr,omitempty"`
}

func (k *Attribute) String() (string, error) {
	doc := ""
	if k.Doc != "" {
		doc = " // " + k.Doc
	}

	dataType, err := k.Type.String()
	if err != nil {
		return k.ID, err
	}

	if dataType == "[]byte" {
		if k.Size != "" {
			_, err := strconv.ParseInt(k.Size, 0, 0)
			if err != nil {
				return k.ID, err
			}
			dataType = strings.Replace(dataType, "[]", "["+k.Size+"]", 1)
		} else {
			dataType = strings.Replace(dataType, "[]", "[2]", 1) // TODO
		}
	}

	if k.Repeat != "" {
		if k.RepeatExpr != "" {
			dataType = "[" + goify(k.RepeatExpr) + "]" + dataType
		}
	}

	return strcase.ToCamel(k.ID) + " " + dataType + "`ks:\"" + k.ID + ",attribute\"`" + doc + "\n", nil
}

type Kaitai struct {
	Types     map[string]Kaitai         `yaml:"types,omitempty"`
	Seq       []Attribute               `yaml:"seq,omitempty"`
	Enums     map[string]map[int]string `yaml:"enums,omitempty"`
	Doc       string                    `yaml:"doc,omitempty"`
	Instances map[string]Instance       `yaml:"instances,omitempty"`
}

func (k *Kaitai) String() (string, error) {
	s := ""

	// print doc string
	if k.Doc != "" {
		s += "// " + strings.Replace(strings.TrimSpace(k.Doc), "\n", "\n// ", -1) + "\n"
	}

	// print type start
	s += "type %[1]s struct{\n"
	s += "\truntime.KaitaiHeader\n"
	// s += "\tIo *runtime.Stream\n"
	// s += "\tParent interface{} \n"
	// s += "\tRoot *%[3]s\n\n"

	// print attribute
	hasCustomTypes := false
	for _, attribute := range k.Seq {
		attrStr, err := attribute.String()
		if err != nil {
			log.Printf("Error in %s\n", attrStr)
			return "", err
		}
		if _, ok := typeMapping[attribute.Type.Type]; !ok {
			hasCustomTypes = true
		}
		s += "\t" + attrStr
	}

	hasValueInstances := false
	for name, instance := range k.Instances {
		hasCustomTypes = true
		attrStr, err := instance.String()
		if err != nil {
			return "", err
		}
		if instance.Value != "" {
			hasValueInstances = true
		}

		s += "\t" + fmt.Sprintf(attrStr, strcase.ToCamel(name), name)
	}

	// print type end
	s += "}\n\n"

	if hasCustomTypes && hasValueInstances {
		s += "func (k *%[1]s) KSYDecode(d runtime.Stream) (err error) {\n"

		// s += "\td := runtime.NewDecoder(reader)\n"
		for _, attribute := range k.Seq {
			reference := "&"
			s += "\td.Decode(" + reference + "k." + strcase.ToCamel(attribute.ID) + ")\n"
		}

		for name, instance := range k.Instances {
			if instance.Pos != "" {
				s += "\td.DecodePos(&k." + strcase.ToCamel(name) + ", " + goify(instance.Pos) + ")\n"
			}
		}

		if !hasValueInstances {
			s += "\treturn d.Err\n"
		} else {
			s += "\tif d.Err != nil {\n"
			s += "\t\treturn d.Err\n"
			s += "\t}\n"

			for name, instance := range k.Instances {
				if instance.Pos == "" {
					dataType, err := instance.dataType()
					if err != nil {
						return "", err
					}
					s += "\tk." + strcase.ToCamel(name) + " = " + dataType + "(" + goify(instance.Value) + ")\n"
				}
			}
			s += "\treturn nil\n"
		}
		s += "}\n\n"
	}

	// print subtypes (flattened)
	for name, t := range k.Types {
		typeStr, err := t.String()
		if err != nil {
			return "", err
		}
		s += fmt.Sprintf(typeStr, strcase.ToCamel(name), name, "%[1]s")
	}

	for enum, values := range k.Enums {
		s += "var " + strcase.ToCamel(enum) + " = struct {\n"
		for _, value := range values {
			s += "\t" + strcase.ToCamel(value) + " int\n"
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

	// write go code
	goCode, err := kaitai.String()
	if err != nil {
		return errors.Wrap(err, "kaitai code gen")
	}
	goCode = fmt.Sprintf(goCode, baseStruct, baseStructSnake, baseStruct)
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
