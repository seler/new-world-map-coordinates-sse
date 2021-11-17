package position

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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
	grabber             Grabber
	ocr                 OCREngine
	lastSampleCollected time.Time
	collectSamples      bool
}

type Position struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
	Timestamp time.Time
}

const COLLECT_SAMPLE_DATA = true
const SAMPLE_DATA_TIMEOUT = 1 * time.Minute

func (from Position) Distance(to Position) float64 {
	return math.Sqrt(
		math.Pow(from.Latitude-to.Latitude, 2) + math.Pow(from.Longitude-to.Longitude, 2),
	)
}

func NewPositionService(grabber Grabber, ocrEngine OCREngine, collectSamples bool) *PositionService {
	return &PositionService{grabber, ocrEngine, time.Now(), collectSamples}
}

func saveSample(position Position, img image.Image, text string, err error) {
	// discard position 0, 0
	if err != nil && position.Longitude == 0 && position.Latitude == 0 {
		return
	}

	// generate image filename
	var valid string
	if err != nil {
		valid = "invalid"
	} else {
		valid = "valid"
	}
	imageFilename := fmt.Sprintf("samples/sample-%.3f-%.3f-%d-%s.png",
		position.Longitude, position.Latitude, position.Timestamp.Unix(), valid)

	// save image file
	file, err := os.Create(imageFilename)
	if err != nil {
		log.Errorf("could not create image file: %v", err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Errorf("could not save image file: %v", err)
		os.Remove(imageFilename)
		return
	}

	// append position, text, !!err, image filename to csv
	SAMPLES_FILENAME := "samples/samples.csv"
	samplesFile, err := os.OpenFile(SAMPLES_FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("could not save image file: %v", err)
		os.Remove(imageFilename)
		return
	}
	defer samplesFile.Close()
	id := uuid.NewString()
	data := fmt.Sprintf("%s, %.3f, %.3f, %d, %s, %s\n",
		id, position.Longitude, position.Latitude, position.Timestamp.Unix(), valid, imageFilename)
	samplesFile.WriteString(data)
}

func (p *PositionService) GetPosition() (Position, error) {

	img := p.grabber.Grab()
	text := p.ocr.GetText(img)
	log.Debugf("Parsed text: %v", text)
	position, err := getPositionFromText(text)

	if COLLECT_SAMPLE_DATA && (position.Timestamp.Sub(p.lastSampleCollected) > SAMPLE_DATA_TIMEOUT) {
		p.lastSampleCollected = position.Timestamp
		if p.collectSamples {
			go saveSample(position, img, text, err)
		}
	}

	return position, err
}

func fixCommonMistakes(text string) string {
	commonMistakes := [][]string{
		{"l", "1"},
		{"L", "1"},
		{"S", "5"},
		{"B", "8"},
		{"&", "8"},
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
	validID := regexp.MustCompile(`-?\d{2,5}[.,]{1,2}[ ]?\d{3}`)
	location := validID.FindAllString(text, -1)

	now := time.Now()

	if len(location) == 2 {
		lng, lat := location[0], location[1]

		lngN := parsePosition(lng)
		latN := parsePosition(lat)

		if latN >= MIN_LAT &&
			lngN >= MIN_LNG &&
			latN <= MAX_LAT &&
			lngN <= MAX_LNG {
			return Position{lngN, latN, now}, nil
		}
		return Position{lngN, latN, now}, fmt.Errorf("position (%v, %v) out of bounds", lngN, latN)
	} else {
		return Position{0, 0, now}, fmt.Errorf("not able to find position in %v", text)
	}
}
