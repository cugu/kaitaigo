package main

var kaitaiTypes map[string]string

func init() {
	kaitaiTypes = map[string]string{
		"Itoa":    "runtime.Bytes",
		"len":     "runtime.Int64",
		"bool":    "bool",
		"string":  "runtime.Bytes",
		"int":     "runtime.Int32",
		"int8":    "runtime.Int8",
		"int16":   "runtime.Int16",
		"int32":   "runtime.Int32",
		"int64":   "runtime.Int64",
		"uint":    "runtime.Uint",
		"uint8":   "runtime.Uint8",
		"uint16":  "runtime.Uint16",
		"uint32":  "runtime.Uint32",
		"uint64":  "runtime.Uint64",
		"byte":    "runtime.Byte",
		"rune":    "runtime.Byte",
		"float32": "runtime.Float32",
		"float64": "runtime.Float64",
	}
}

func addKaitaiType(kaitaiName, kaitaiType string) {
	if _, ok := kaitaiTypes[kaitaiName]; !ok {
		kaitaiTypes[kaitaiName] = kaitaiType
	}
}

func getKaitaiType(kaitaiName string) string {
	if kaitaiType, ok := kaitaiTypes[kaitaiName]; ok {
		return kaitaiType
	}
	return "runtime.KSYDecoder"
}
