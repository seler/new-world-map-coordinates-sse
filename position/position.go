package position

import (
	"fmt"
	"image"
	"math"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type OCREngine interface {
	Init()
	GetText(img image.Image) string
	Close()
}

type Grabber interface {
	Grab() image.Image
}

type PositionService struct {
	grabber Grabber
	ocr     OCREngine
}

type Position struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
}

func (from Position) Distance(to Position) float64 {
	return math.Sqrt(
		math.Pow(from.Latitude-to.Latitude, 2) + math.Pow(from.Longitude-to.Longitude, 2),
	)
}

func NewPositionService(grabber Grabber, ocrEngine OCREngine) *PositionService {
	return &PositionService{grabber, ocrEngine}
}

func (p *PositionService) GetPosition() (Position, error) {
	img := p.grabber.Grab()
	text := p.ocr.GetText(img)
	log.Debugf("Parsed text: %v", text)
	return getPositionFromText(text)
}

func fixCommonMistakes(text string) string {
	commonMistakes := [][]string{
		{"l", "1"},
		{"L", "1"},
		{"S", "5"},
	}
	for _, m := range commonMistakes {
		text = strings.Replace(text, m[0], m[1], -1)
	}
	return text
}

func parsePosition(text string) float64 {
	text = strings.Replace(text, ",", ".", -1)
	text = strings.Replace(text, "..", ".", -1)
	text = strings.Replace(text, ". ", ".", -1)
	text = strings.Trim(text, ".")
	n, err := strconv.ParseFloat(text, 64)
	if err != nil {
		log.Errorf("Could not recognize number in position text: %v %v", text, err)
	}
	return n
}

const MIN_LNG, MIN_LAT, MAX_LNG, MAX_LAT = 4000, 0, 15000, 11000

func getPositionFromText(text string) (Position, error) {
	text = fixCommonMistakes(text)
	validID := regexp.MustCompile(`-?\d{2,5}[.,]{1,2}[ ]?\d{3}[.,]`)
	location := validID.FindAllString(text, -1)

	if len(location) == 2 {
		lng, lat := location[0], location[1]

		lngN := parsePosition(lng)
		latN := parsePosition(lat)

		if latN >= MIN_LAT &&
			lngN >= MIN_LNG &&
			latN <= MAX_LAT &&
			lngN <= MAX_LNG {
			return Position{lngN, latN}, nil
		}
		return Position{lngN, latN}, fmt.Errorf("position (%v, %v) out of bounds", lngN, latN)
	} else {
		return Position{0, 0}, fmt.Errorf("not able to find position in %v", text)
	}
}
