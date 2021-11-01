package main

// #cgo windows CXXFLAGS: -std=c++0x -Iinclude
// #cgo windows LDFLAGS: -Ltesseract -llibtesseract-4 -Llept -lliblept-5
// #include <stdlib.h>
// #include <stdbool.h>
// #include "coordinates.h"
import "C"

func Hello() string {
	return C.GoString(C.Hello())
}
