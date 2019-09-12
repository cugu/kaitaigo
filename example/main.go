package main

import (
	"bytes"
	"log"
)

func main() {
	f := bytes.NewReader([]byte("\x05Hello world!"))
	var r MyFormat
	err := r.Decode(f) // Decode takes any io.ReadSeeker
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(r.Data())) // Prints "Hello"
}
