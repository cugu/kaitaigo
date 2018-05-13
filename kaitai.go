package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/kr/pretty"
	yaml "gopkg.in/yaml.v2"
)

var typeMapping = map[string]string{
	"u1":   "uint8",
	"u2":   "uint16",
	"u4":   "uint32",
	"u8":   "uint64",
	"u2le": "uint16",
	"u2be": "uint16",
	"u4le": "uint32",
	"u4be": "uint32",
	"u8le": "uint64",
	"u8be": "uint64",
	"s1":   "int8",
	"s2":   "int16",
	"s4":   "int32",
	"s8":   "int64",
	"s2le": "int16",
	"s2be": "int16",
	"s4le": "int32",
	"s4be": "int32",
	"s8le": "int64",
	"s8be": "int64",
	"f4":   "float32",
	"f8":   "float64",
	"f4be": "float32",
	"f4le": "float32",
	"f8be": "float64",
	"f8le": "float64",
	"str":  "string",
	"strz": "string",
	"":     "[]byte",
}

func goify(s string) string {
	// Create go versions of vars
	re := regexp.MustCompile("[a-z][a-z_]*")
	s = re.ReplaceAllStringFunc(s, strcase.ToCamel)

	re = regexp.MustCompile("_?[A-Z][a-zA-Z0-9_<>\\.]*")
	return re.ReplaceAllStringFunc(s, func(s string) string {
		return "k." + s
	})
}

type Type struct {
	Type string
}

func (s *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&s.Type)
	if err != nil {
		s.Type = "Kaitai"
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

func (k *Instance) String() (string, error) {
	doc := ""
	if k.Doc != "" {
		doc = " // " + k.Doc
	}

	dataType, err := k.Type.String()
	if err != nil {
		return "", err
	}
	return "STRUCT_NAME " + dataType + doc + " //instance\n", nil
}

type Attribute struct {
	ID   string `yaml:"id,omitempty"`
	Type Type   `yaml:"type"`
	Size string `yaml:"size,omitempty"`
	Doc  string `yaml:"doc,omitempty"`
}

func (k *Attribute) String() (string, error) {
	doc := ""
	if k.Doc != "" {
		doc = " // " + k.Doc
	}

	dataType, err := k.Type.String()
	if err != nil {
		return "", err
	}

	if dataType == "[]byte" {
		size, err := strconv.Atoi(k.Size)
		if err != nil {
			return "", err
		}
		dataType = strings.Replace(dataType, "[]", "["+strconv.Itoa(size)+"]", 1)
	}

	return strcase.ToCamel(k.ID) + " " + dataType + doc + "\n", nil
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
	s += "type STRUCT_NAME struct{\n"

	// print attribute
	hasCustomTypes := false
	for _, attribute := range k.Seq {
		attrStr, err := attribute.String()
		if err != nil {
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

		s += "\t" + strings.Replace(attrStr, "STRUCT_NAME", strcase.ToCamel(name), 1)
	}

	// print type end
	s += "}\n\n"

	s += "func (k *STRUCT_NAME) Decode(reader io.ReadSeeker) (err error) {\n"

	if !hasCustomTypes {
		s += "\treturn binary.Read(reader, binary.LittleEndian, k)\n"
	} else {
		s += "\td := decoder{reader: reader, byteOrder: binary.LittleEndian}\n"
		for _, attribute := range k.Seq {
			s += "\td.decode(k." + strcase.ToCamel(attribute.ID) + ")\n"
		}

		for name, instance := range k.Instances {
			if instance.Pos != "" {
				s += "\td.decodePos(k." + strcase.ToCamel(name) + ", " + goify(instance.Pos) + ")\n"
			}
		}

		if !hasValueInstances {
			s += "\treturn d.err\n"
		} else {

			s += "\tif d.err != nil {\n"
			s += "\t\treturn d.err\n"
			s += "\t}\n"

			for name, instance := range k.Instances {
				if instance.Pos == "" {
					s += "\tk." + strcase.ToCamel(name) + " = (" + goify(instance.Value) + ")\n"
				}
			}

			s += "\treturn nil\n"
		}
	}
	s += "}\n\n"

	// print subtypes (flattened)
	for name, t := range k.Types {
		typeStr, err := t.String()
		if err != nil {
			return "", err
		}
		s += strings.Replace(typeStr, "STRUCT_NAME", strcase.ToCamel(name), 2)
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

func createGofile(filepath string, pckg string) {

	outdir := "test_formats"
	filename := path.Base(filepath)

	logfile, err := os.Create(path.Join(outdir, filename+".log"))
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	source, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// print all

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(source), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//
	err = ioutil.WriteFile(
		path.Join(outdir, filename+".generic.unmarshal"),
		[]byte(fmt.Sprintf("%# v\n", pretty.Formatter(m))),
		0644,
	)
	if err != nil {
		log.Printf("error: %v", err)
	}

	// parse kaitai

	kaitai := Kaitai{}

	err = yaml.UnmarshalStrict([]byte(source), &kaitai)
	if err != nil {
		log.Printf("error: %v", err)
	}
	// fmt.Printf("%# v\n", )
	err = ioutil.WriteFile(
		path.Join(outdir, filename+".kaitai.unmarshal"),
		[]byte(fmt.Sprintf("%# v\n", pretty.Formatter(kaitai))),
		0644,
	)
	if err != nil {
		log.Printf("error: %v", err)
	}

	baseStruct := strcase.ToCamel(strings.Replace(filename, ".ksy", "", 1))

	// write go code
	goCode, err := kaitai.String()
	goCode = strings.Replace(goCode, "STRUCT_NAME", baseStruct, 2)
	header := "package " + pckg + "\n"
	header += "\n"
	header += "import (\n"
	header += "\t\"encoding/binary\"\n"
	header += "\t\"fmt\"\n"
	header += "\t\"io\"\n"
	header += "\t\"os\"\n"
	header += "\t\"log\"\n"
	header += ")\n"
	header += "\n"
	main := "func main() {\n"
	main += "\tf, err := os.Open(os.Args[1])\n"
	main += "\tif err != nil {\n"
	main += "\t\tlog.Fatal(err)\n"
	main += "\t}\n"
	main += "\tdefer f.Close()\n"
	main += "\tbaseStruct := " + baseStruct + "{}\n"
	main += "\tx := baseStruct.Decode(f)\n"
	main += "\tfmt.Printf(\"%v\", x)\n"
	main += "}\n"
	if err != nil {
		log.Printf("error: %v", err)
	}
	err = ioutil.WriteFile(
		path.Join(outdir, filename+".go"),
		[]byte(header+goCode+main),
		0644,
	)
	if err != nil {
		log.Printf("error: %v", err)
	}

}

func main() {
	for _, filename := range os.Args[1:] {
		createGofile(filename, "test_formats")
	}
}
