package main

import (
	"fmt"
	"image"

	"github.com/kbinani/screenshot"
)

type Grabber struct {
	display int
	bounds  image.Rectangle
}

func NewGrabber(display int) Grabber {
	bounds := screenshot.GetDisplayBounds(display)
	// seems that New World position info is in same position and scale in all resolutions
	positionInfoBounds := image.Rect(
		bounds.Max.X-286,
		bounds.Min.Y+19,
		bounds.Max.X-5,
		bounds.Min.Y+35)
	return Grabber{display, positionInfoBounds}
}

func (g *Grabber) grab() image.Image {
	img, err := screenshot.CaptureRect(g.bounds)
	if err != nil {
		panic(err)
	}
	return img
}

func main() {
	display := 0

	grabber := NewGrabber(display)
	img := grabber.grab()
	ocr := NewOCRClient()
	text := ocr.GetText(img)
	fmt.Printf("text >> %s <<", text)
}
