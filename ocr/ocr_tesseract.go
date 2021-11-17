package ocr

// #cgo windows CXXFLAGS: -std=c++0x -Iinclude
// #cgo windows LDFLAGS: -Ltesseract -llibtesseract-4 -Llept -lliblept-5
// #include <stdlib.h>
// #include <stdbool.h>
// #include "ocr_tesseract.h"
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

type TesseractClient struct {
	api        C.TessBaseAPI
	saveImages bool
}

func NewTesseractClient(saveImages bool) *TesseractClient {
	api := C.TessNew()
	return &TesseractClient{api, saveImages}
}

func (t *TesseractClient) Init() {
	C.TessInit(t.api)
}

func (t *TesseractClient) GetText(img image.Image) string {
	imgBytes := imageToBytes(img)

	saveImages := 0
	if t.saveImages {
		saveImages = 1
	}

	text := C.TessGetText(
		t.api,
		(*C.uchar)(unsafe.Pointer(&imgBytes[0])),
		C.int(len(imgBytes)),
		C.int(saveImages),
	)
	defer C.free(unsafe.Pointer(text))
	return C.GoString(text)
}

func (t *TesseractClient) Close() {
	C.TessEnd(t.api)
}
