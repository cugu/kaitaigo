package main

var kaitaiTypes map[string]string

func isNative(dataType string) bool {
	nativeTypes := map[string]bool{
		"bool":    true,
		"string":  true,
		"int":     true,
		"int8":    true,
		"int16":   true,
		"int32":   true,
		"int64":   true,
		"uint":    true,
		"uint8":   true,
		"uint16":  true,
		"uint32":  true,
		"uint64":  true,
		"byte":    true,
		"rune":    true,
		"float32": true,
		"float64": true,
		"[]byte":  true,
	}
	if _, ok := nativeTypes[dataType]; ok {
		return true
	}
	return false
}

func addKaitaiType(kaitaiName, kaitaiType string) {
	if val, ok := kaitaiTypes[kaitaiName]; !ok || val == "runtime.KSYDecoder" {
		kaitaiTypes[kaitaiName] = kaitaiType
	}
}

func getKaitaiType(kaitaiName string) string {
	if kaitaiType, ok := kaitaiTypes[kaitaiName]; ok {
		return kaitaiType
	}
	return "runtime.KSYDecoder"
}
