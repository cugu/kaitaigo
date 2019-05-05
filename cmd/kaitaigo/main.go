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

func YAMLUnmarshal(name string, source []byte, m interface{}, path string, debug bool) error {
	err := yaml.Unmarshal(source, m)
	if err != nil || !debug {
		return err
	}
	return ioutil.WriteFile(
		path+"."+name+".unmarshal",
		[]byte(fmt.Sprintf("%s%# v\n", "// file generated at "+time.Now().UTC().Format(time.RFC3339)+"\n", pretty.Formatter(m))),
		0644,
	)
}

func createGofile(ksyPath, pkg string, debug bool) error {
	filename := path.Base(ksyPath)
	dir := filepath.Dir(ksyPath)

	// setup logging
	if debug {
		logfile, err := os.Create(path.Join(dir, filename+".log"))
		if err != nil {
			return errors.Wrap(err, "create logfile")
		}
		defer func() {
			logfile.Sync()
			logfile.Close()
		}()
		log.SetOutput(io.MultiWriter(os.Stderr, logfile))
	}

	// start generation
	log.Println("generate", ksyPath)

	// read source
	source, err := ioutil.ReadFile(ksyPath)
	if err != nil {
		return errors.Wrap(err, "read source")
	}

	// parse generic
	m := make(map[interface{}]interface{})
	err = YAMLUnmarshal("generic", source, &m, path.Join(dir, filename), debug)
	if err != nil {
		return errors.Wrap(err, "parse generic yaml")
	}

	// parse kaitai
	kaitai := Type{}
	enumTypes = map[string]string{}
	parents = map[string]string{}
	kaitaiTypes = map[string]string{
		"Itoa": "[]byte",
		"len":  "int64",
	}
	err = YAMLUnmarshal("kaitai", source, &kaitai, path.Join(dir, filename), debug)
	if err != nil {
		return errors.Wrap(err, "parse kaitai yaml")
	}
	baseStruct := strcase.ToCamel(kaitai.Meta.ID)

	setupMap(&kaitai, baseStruct)
	setupMap(&kaitai, baseStruct)

	// write go code
	var buffer LineBuffer
	buffer.WriteLine("// file generated at " + time.Now().UTC().Format(time.RFC3339) + "\n")
	buffer.WriteLine("package " + pkg)
	buffer.WriteLine("import (\"gitlab.com/dfir/binary/kaitaigo/runtime\")")
	buffer.WriteLine("var decoder io.ReadSeeker")
	buffer.WriteLine(kaitai.String(baseStruct, baseStruct, baseStruct))

	// format and add imports
	formated, err := imports.Process("", []byte(buffer.String()), nil)
	if err != nil {
		log.Printf("Format error (%s): %s", ksyPath, err)
		formated = []byte(buffer.String())
	}
	err = ioutil.WriteFile(path.Join(dir, filename+".go"), formated, 0644)

	return errors.Wrap(err, "create go file")
}

func handleFile(filename, pkg string, debug bool) error {
	if strings.HasSuffix(filename, ".ksy") {
		return createGofile(filename, pkg, debug)
	}
	return nil
}

func main() {
	debug := flag.Bool("debug", false, "debug output")
	flag.Parse()
	for _, filename := range flag.Args() {
		var err error
		if strings.HasSuffix(filename, "/...") {
			recPath := strings.Replace(filename, "/...", "", 1)
			err = filepath.Walk(recPath, func(path string, f os.FileInfo, err error) error {
				abspath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				return handleFile(path, filepath.Base(filepath.Dir(abspath)), *debug)
			})
		} else {
			abspath, err := filepath.Abs(filename)
			if err == nil {
				err = handleFile(filename, filepath.Base(filepath.Dir(abspath)), *debug)
			}
		}
		if err != nil {
			log.Println(err)
		}
	}
}
