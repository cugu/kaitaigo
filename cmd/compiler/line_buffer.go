package main

import "strings"

type LineBuffer struct {
	strings.Builder
}

func (lb *LineBuffer) WriteLine(s string) {
	lb.WriteString(s + "\n")
}
