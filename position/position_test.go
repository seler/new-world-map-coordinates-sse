package position

import (
	"fmt"
	"image"
	"strings"
	"testing"
)

type FakeOCR struct {
	text string
}

func (o *FakeOCR) Init() {}

func (o *FakeOCR) GetText(img image.Image) string {
	return o.text
}

func (o *FakeOCR) Close() {}

type FakeGrabber struct{}

func (g FakeGrabber) Grab() image.Image {
	return image.NewNRGBA(image.Rect(0, 0, 1, 1))
}

func TestGetPositionReturnsParsedPosition(t *testing.T) {
	expected := Position{4123.456, 789.123}
	grabber := FakeGrabber{}
	ocr := &FakeOCR{fmt.Sprintf("%.3f, %.3f", expected.Longitude, expected.Latitude)}
	p := NewPositionService(grabber, ocr)

	got, err := p.GetPosition()
	if err != nil {
		t.Errorf("Got error while trying to get position: %v", err)
	}
	if got != expected {
		t.Errorf("Got %v; expected %v", got, expected)
	}
}

func TestGetPositionReturnsParsesNegativePosition(t *testing.T) {
	expected := Position{-4123.456, -789.123}
	grabber := FakeGrabber{}
	ocr := &FakeOCR{fmt.Sprintf("%.3f, %.3f", expected.Longitude, expected.Latitude)}
	p := NewPositionService(grabber, ocr)

	got, _ := p.GetPosition()
	if got != expected {
		t.Errorf("Got %v; expected %v", got, expected)
	}
}

func TestGetPositionRaisesOutOfBounds(t *testing.T) {
	lngOutOfMinBound := Position{MIN_LNG - 0.001, MIN_LAT + 123.456}
	lngOutOfMaxBound := Position{MAX_LNG + 0.001, MIN_LAT + 123.456}
	latOutOfMinBound := Position{MIN_LNG + 1, MIN_LAT - 123.456}
	latOutOfMaxBound := Position{MIN_LNG + 1, MAX_LAT + 123.456}
	positions := []Position{
		lngOutOfMinBound,
		lngOutOfMaxBound,
		latOutOfMinBound,
		latOutOfMaxBound,
	}
	grabber := FakeGrabber{}

	for _, pos := range positions {
		ocr := &FakeOCR{fmt.Sprintf("%.3f, %.3f", pos.Longitude, pos.Latitude)}
		p := NewPositionService(grabber, ocr)
		got, err := p.GetPosition()
		if err == nil {
			t.Errorf("Expected out of bounds error but got position: %v", got)
		}
		if err != nil && !strings.Contains(err.Error(), "out of bounds") {
			t.Errorf("Expected out of bounds error but got error: %v", err)
		}
	}
}
