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
	if dad, ok := parents[typeName]; ok {
		return dad
	}
	return "runtime.KSYDecoder"
}

func setupMap(k *Kaitai, typeName string) {
	for _, attribute := range k.Seq {
		addParent(strcase.ToCamel(attribute.Type.Type), strcase.ToCamel(typeName))
		if attribute.Type.TypeSwitch.SwitchOn != "" {
			for _, casetype := range attribute.Type.TypeSwitch.Cases {
				addParent(strcase.ToCamel(casetype.Type), strcase.ToCamel(typeName))
			}
		}
	}
	for _, instance := range k.Instances {
		addParent(strcase.ToCamel(instance.Type.Type), strcase.ToCamel(typeName))
		if instance.Type.TypeSwitch.SwitchOn != "" {
			for _, casetype := range instance.Type.TypeSwitch.Cases {
				addParent(strcase.ToCamel(casetype.Type), strcase.ToCamel(typeName))
			}
		}
	}

	for name, t := range k.Types {
		setupMap(&t, name)
	}
}
