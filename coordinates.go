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

func imageToBytes(img image.Image) []byte {
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}

type OCRClient struct {
	api C.TessBaseAPI
}

func NewOCRClient() *OCRClient {
	api := C.TessNew()
	return &OCRClient{api}
}

func (c *OCRClient) Init() {
	C.TessInit(c.api)
}

func (c *OCRClient) GetText(img image.Image) string {
	imgBytes := imageToBytes(img)

	text := C.TessGetText(
		c.api,
		(*C.uchar)(unsafe.Pointer(&imgBytes[0])),
		C.int(len(imgBytes)),
	)
	defer C.free(unsafe.Pointer(text))
	return C.GoString(text)
}

func (c *OCRClient) End() {
	C.TessEnd(c.api)
}
