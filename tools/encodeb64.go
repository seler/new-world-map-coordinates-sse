//go:build ignore

package main

import (
	"encoding/base64"
	"io"
	"os"
)

// Encodes any file as base64 string and prints to stdout
func main() {
	if len(os.Args[1:]) != 1 {
		panic("invalid number of arguments")
	}
	filename := os.Args[1]
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	io.Copy(encoder, f)
	encoder.Close()
}
