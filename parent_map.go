package main

import (
	"github.com/iancoleman/strcase"
)

var parents map[string]string

func addParent(typeName, parent string) {
	if _, ok := parents[typeName]; !ok {
		parents[typeName] = parent
	}
}

func getParent(typeName string) string {
	if parent, ok := parents[typeName]; ok {
		return parent
	}
	return "runtime.KSYDecoder"
}

func prepare(attr Attribute, typeName string) {
	addKaitaiType(strcase.ToCamel(attr.Name()), attr.DataType())
	if attr.Enum != "" {
		addEnumType(attr.Enum, attr.DataType())
	}
	addParent(strcase.ToCamel(attr.Type.Type), strcase.ToCamel(typeName))
	if attr.Type.TypeSwitch.SwitchOn != "" {
		for _, casetype := range attr.Type.TypeSwitch.Cases {
			addParent(strcase.ToCamel(casetype.Type), strcase.ToCamel(typeName))
		}
	}
}

func setupMap(k *Type, typeName string) {
	for _, attr := range k.Seq {
		prepare(attr, typeName)
	}
	for name, attr := range k.Instances {
		attr.ID = name
		prepare(attr, typeName)
	}

	for name, t := range k.Types {
		setupMap(&t, name)
	}
}
