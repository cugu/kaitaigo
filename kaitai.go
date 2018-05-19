package main

import (
	"flag"
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

func goify(s string, t string) string {
	// Create go versions of vars


	re := regexp.MustCompile("[a-z][a-z_]*")
	s = re.ReplaceAllStringFunc(s, strcase.ToCamel)

	s = strings.Replace(s, "_", "", -1)

	re = regexp.MustCompile("_?[a-zA-Z0-9_\\.]+")
	return re.ReplaceAllStringFunc(s, func(s string) string {
		_, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return s
		}
		if strings.HasPrefix(s, "0") {
			return strings.Replace(strings.ToLower(s), "0b", "0", 1)
		}
		return t + "(k." + s + ")"
	})
}

type Type struct {
	Type string
}

func (s *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := unmarshal(&s.Type)
	if err != nil {
		s.Type = "ks.KSYDecoder"
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
	if dataType == "[]byte" && k.Value != "" {
		dataType = "int64"
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
			dataType = "[" + goify(k.RepeatExpr, "") + "]" + dataType
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
	s += "\tIo *ks.Stream\n"
	s += "\tParent interface{} \n"
	s += "\tRoot *%[3]s\n\n"

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
		s += "func (k *%[1]s) KSYDecode(d ks.Stream) (err error) {\n"

		// s += "\td := ks.NewDecoder(reader)\n"
		for _, attribute := range k.Seq {
			reference := "&"
			s += "\td.Decode(" + reference + "k." + strcase.ToCamel(attribute.ID) + ")\n"
		}

		for name, instance := range k.Instances {
			if instance.Pos != "" {
				s += "\td.DecodePos(&k." + strcase.ToCamel(name) + ", " + goify(instance.Pos, "") + ")\n"
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
					s += "\tk." + strcase.ToCamel(name) + " = " + goify(instance.Value, "int64") + "\n"
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

func createGofile(filepath string, pckg string) {

	filename := path.Base(filepath)

	logfile, err := os.Create(path.Join(pckg, filename+".log"))
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
		path.Join(pckg, filename+".generic.unmarshal"),
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
		path.Join(pckg, filename+".kaitai.unmarshal"),
		[]byte(fmt.Sprintf("%# v\n", pretty.Formatter(kaitai))),
		0644,
	)
	if err != nil {
		log.Printf("error: %v", err)
	}

	baseStructSnake := strings.Replace(filename, ".ksy", "", 1)
	baseStruct := strcase.ToCamel(baseStructSnake)

	// write go code
	goCode, err := kaitai.String()
	goCode = fmt.Sprintf(goCode, baseStruct, baseStructSnake, baseStruct)
	header := "package " + pckg + "\n"
	header += "\n"
	header += "import (\n"
	for _, pkg := range []string{"fmt", "io", "os", "log", "ks"} {
		header += "\t\"" + pkg + "\"\n"
	}
	header += ")\n"
	header += "\n"
	//main := "func main() {\n"
	//main += "\tf, err := os.Open(os.Args[1])\n"
	//main += "\tif err != nil {\n"
	//main += "\t\tlog.Fatal(err)\n"
	//main += "\t}\n"
	//main += "\tdefer f.Close()\n"
	//main += "\tbaseStruct := " + baseStruct + "{}\n"
	//main += "\tx := baseStruct.Decode(f)\n"
	//main += "\tfmt.Printf(\"%v\", x)\n"
	//main += "}\n"
	if err != nil {
		log.Printf("error: %v", err)
	}
	err = ioutil.WriteFile(
		path.Join(pckg, filename+".go"),
		[]byte(header+goCode), //+main),
		0644,
	)
	if err != nil {
		log.Printf("error: %v", err)
	}

}

func main() {
	var outdir = flag.String("d", "out", "the species we are studying")
	flag.Parse()
	os.MkdirAll(*outdir, 0755)
	for _, filename := range flag.Args() {
		createGofile(filename, *outdir)
	}
}
