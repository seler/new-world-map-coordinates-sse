package main

import (
	"os"

	"github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var version = "not set"
var logLevel = "not set"
var saveNotRecognized = "not set"

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		DisableQuote:           true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		FullTimestamp:          true,
	})
	log.Infof("new-world-map-coordinates-sse-%s", version)
	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		ll = logrus.InfoLevel
	}
	log.SetLevel(ll)
	log.Debug("running in debug log level")
}

var usage = `New World Map Coordinates SSE: A simple service that exposes New World's player's position as SSE stream.
Usage:
  map-coordinates-sse.exe --version
  map-coordinates-sse.exe [--display=<n>] [--bind=<n>]

Options:
  -h --help		Show this screen.
  --version     Show version.
  --display=<n> Display to use [default: 0].
  --bind=<n>    addr and port to serve on [default: :5000].
`

func main() {
	opts, err := docopt.ParseArgs(usage, os.Args[1:], version)
	if err != nil {
		panic(err)
	}

	display, err := opts.Int("--display")
	if err != nil {
		panic(err)
	}
	addr, err := opts.String("--bind")
	if err != nil {
		panic(err)
	}

	log.Info("Hold CTRL+C to stop\n")

	log.Debugf("Using display=%v", display)
	config := Config{
		display: display,
		addr:    addr,
	}
	mapCoordinates(config)
}
