package main

var enumTypes map[string]string

func addEnumType(enumName, enumType string) {
	if _, ok := enumTypes[enumName]; !ok {
		enumTypes[enumName] = enumType
	}
}

func getEnumType(enumName string) string {
	if enumType, ok := enumTypes[enumName]; ok {
		return enumType
	}
	return "int64"
}
