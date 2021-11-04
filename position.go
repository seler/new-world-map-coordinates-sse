package main

import (
	"fmt"
	"image"
)

type OCREngine interface {
	Init()
	GetText(img image.Image) string
	End()
}

type Grabber interface {
	grab() image.Image
}

type PositionService struct {
	grabber Grabber
	ocr     OCREngine
}

type Position struct {
	latitude  float32
	longitude float32
	height    float32
}

func NewPositionService(grabber Grabber, ocrEngine OCREngine) *PositionService {
	return &PositionService{grabber, ocrEngine}
}

func (p *PositionService) GetPosition() Position {
	img := p.grabber.grab()
	text := p.ocr.GetText(img)
	fmt.Printf("text2 >> %s <<", text)
	return Position{0, 0, 0}
}
