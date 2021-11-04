package main

import (
	"image"
	"image/png"
	"os"

	"github.com/kbinani/screenshot"
)

type ScreenGrabber struct {
	display int
	bounds  image.Rectangle
}

func NewScreenGrabber(display int) *ScreenGrabber {
	bounds := screenshot.GetDisplayBounds(display)
	// seems that New World position info is in same position and scale in all resolutions
	positionInfoBounds := image.Rect(
		bounds.Max.X-286,
		bounds.Min.Y+19,
		bounds.Max.X-5,
		bounds.Min.Y+35)
	return &ScreenGrabber{display, positionInfoBounds}
}

func (g ScreenGrabber) grab() image.Image {
	img, err := screenshot.CaptureRect(g.bounds)
	if err != nil {
		panic(err)
	}
	return img
}

type FakeGrabber struct {
	img image.Image
}

func NewFakeGrabber(fileName string) *FakeGrabber {
	file, _ := os.Open("test.png")
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	return &FakeGrabber{img}
}

func (g FakeGrabber) grab() image.Image {
	return g.img
}
