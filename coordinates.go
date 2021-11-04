package main

// #cgo windows CXXFLAGS: -std=c++0x -Iinclude
// #cgo windows LDFLAGS: -Ltesseract -llibtesseract-4 -Llept -lliblept-5
// #include <stdlib.h>
// #include <stdbool.h>
// #include "coordinates.h"
import "C"
import (
	"bytes"
	"image"
	"image/png"
	"unsafe"
)

type OCRClient struct{}

func NewOCRClient() *OCRClient {
	return &OCRClient{}
}

func imageToBytes(img image.Image) []byte {
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}

func (c *OCRClient) GetText(img image.Image) string {
	imgBytes := imageToBytes(img)

	text := C.GetText(
		(*C.uchar)(unsafe.Pointer(&imgBytes[0])),
		C.int(len(imgBytes)),
	)
	defer C.free(unsafe.Pointer(text))
	return C.GoString(text)
}
