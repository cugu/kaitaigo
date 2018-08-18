package main

var kaitaiTypes map[string]string

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
