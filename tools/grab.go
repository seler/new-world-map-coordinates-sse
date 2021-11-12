//go:build ignore

package main

import (
	"image/png"
	"os"

	"github.com/seler/new-world-map-coordinates-sse/grabber"
	log "github.com/sirupsen/logrus"
)

func main() {
	grabber := grabber.NewScreenGrabber(0)
	img := grabber.Grab()
	f, err := os.Create("grab.png")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
