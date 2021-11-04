package main

import (
	"fmt"
)

func main() {
	// grabber := NewScreenGrabber(0)
	grabber := NewFakeGrabber("test.png")
	ocr := NewOCRClient()
	ocr.Init()
	defer ocr.End()

	p := NewPositionService(grabber, ocr)

	position := p.GetPosition()
	fmt.Printf("position >> %g <<", position)
	position = p.GetPosition()
	fmt.Printf("position >> %g <<", position)
}
