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