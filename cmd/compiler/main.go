package main

import (
	"bytes"
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

	parents = map[string]string{}
	setupMap(&kaitai, baseStruct)

	// write go code
	var buffer bytes.Buffer

	parts := strings.Split(pckg, "/")
	lastpart := parts[len(parts)-1]
	buffer.WriteString("// file generated at " + time.Now().UTC().Format(time.RFC3339) + "\n")
	buffer.WriteString("package " + lastpart + "\n")
	buffer.WriteString("import (\n")
	for _, pkg := range []string{"fmt", "io", "os", "log", "gitlab.com/cugu/kaitai.go/runtime"} {
		buffer.WriteString("\"" + pkg + "\"\n")
	}
	buffer.WriteString(")\n")
	buffer.WriteString(kaitai.String(baseStruct, baseStruct, baseStruct))

	formated, err := imports.Process("", buffer.Bytes(), nil)
	if err != nil {
		log.Print("Format error", err)
		formated = buffer.Bytes()
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
