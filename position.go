package main

import (
	"image"
	"regexp"
	"strconv"
	"strings"
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
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

func NewPositionService(grabber Grabber, ocrEngine OCREngine) *PositionService {
	return &PositionService{grabber, ocrEngine}
}

func (p *PositionService) GetPosition() Position {
	img := p.grabber.grab()
	text := p.ocr.GetText(img)
	return getPositionFromText(text)
}

func parsePosition(text string) float64 {
	lng := strings.Replace(text, ",", ".", -1)
	lng = strings.Replace(lng, "..", ".", -1)
	lngN, _ := strconv.ParseFloat(lng, 64)
	return lngN
}

func getPositionFromText(text string) Position {
	validID := regexp.MustCompile(`\d{2,5}[.,]{1,2}\d{3}`)
	location := validID.FindAllString(text, -1)

	if len(location) == 3 {
		lng, lat := location[0], location[1]

		lngN := parsePosition(lng)
		latN := parsePosition(lat)

		return Position{latN, lngN}
	} else {
		return Position{0, 0}
	}
}
