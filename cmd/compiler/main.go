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
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
	yaml "gopkg.in/yaml.v2"
)

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

	enumTypes = map[string]string{}
	parents = map[string]string{}
	kaitaiTypes = map[string]string{
		"Itoa": "[]byte",
		"len":  "int64",
	}

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
	kaitai := Type{}
	err = YAMLUnmarshal("kaitai", source, &kaitai, path.Join(pckg, filename))
	if err != nil {
		return errors.Wrap(err, "parse kaitai yaml")
	}

	setupMap(&kaitai, baseStruct)
	// fmt.Printf("%#v\n", kaitaiTypes)
	setupMap(&kaitai, baseStruct)
	// fmt.Printf("%#v\n", kaitaiTypes)

	// write go code
	var buffer LineBuffer

	buffer.WriteLine("// file generated at " + time.Now().UTC().Format(time.RFC3339) + "\n")

	parts := strings.Split(pckg, "/")
	lastpart := parts[len(parts)-1]
	buffer.WriteLine("package " + lastpart)
	buffer.WriteLine("import (")
	for _, pkg := range []string{"fmt", "io", "os", "log", "gitlab.com/cugu/kaitai.go/runtime"} {
		buffer.WriteLine("\"" + pkg + "\"")
	}
	buffer.WriteLine(")")
	buffer.WriteLine("var decoder *runtime.Decoder")

	buffer.WriteLine(kaitai.String(baseStruct, baseStruct, baseStruct))

	formated, err := imports.Process("", []byte(buffer.String()), nil)
	if err != nil {
		log.Printf("Format error (%s): %s", filepath, err)
		formated = []byte(buffer.String())
	}
	err = ioutil.WriteFile(path.Join(pckg, filename+".go"), formated, 0644)

	if err != nil {
		return errors.Wrap(err, "create go file")
	}
	return nil

}

func handleFile(filename string) error {
	if strings.HasSuffix(filename, ".ksy") {
		return createGofile(filename, filepath.Dir(filename))
	}
	return nil
}

func main() {
	flag.Parse()
	for _, filename := range flag.Args() {
		var err error
		if strings.HasSuffix(filename, "/...") {
			recPath := strings.Replace(filename, "/...", "", 1)
			err = filepath.Walk(recPath, func(path string, f os.FileInfo, err error) error {
				return handleFile(path)
			})
		} else {
			err = handleFile(filename)
		}
		if err != nil {
			log.Println(err)
		}
	}
}
