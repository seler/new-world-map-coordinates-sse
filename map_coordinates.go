package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	display int
	addr    string
}

type PositionReportCallback func(position Position)

func continuouslyReportPosition(p *PositionService, done <-chan interface{}, callback PositionReportCallback) {
	ticker := time.NewTicker(time.Second / 2)
	gotPosition := make(chan interface{})

	var current, previous Position

	for range ticker.C {
		go func() {
			previous = current
			current = p.GetPosition()
			log.Infof("Position: %+v", current)
			if (current != previous && current != Position{0, 0}) {
				callback(current)
			}
			gotPosition <- nil
		}()
		select {
		case <-done:
			return
		case <-gotPosition:
			continue
		}
	}
}

type Dispatcher struct {
	clients  []chan Position
	position chan Position
}

func (d *Dispatcher) dispatch(done <-chan interface{}) {
	for {
		select {
		case <-done:
			return
		case position := <-d.position:
			for _, client := range d.clients {
				client <- position
			}
		}
	}
}

func mapCoordinates(config Config) {
	grabber := NewScreenGrabber(config.display)
	// grabber := NewFakeGrabber("test.png")
	ocr := NewTesseractClient()
	ocr.Init()

	p := NewPositionService(grabber, ocr)

	dispatcher := &Dispatcher{
		clients:  make([]chan Position, 0),
		position: make(chan Position),
	}

	dispatcherDone := make(chan interface{})
	defer close(dispatcherDone)
	go dispatcher.dispatch(dispatcherDone)

	reportDone := make(chan interface{})
	defer close(reportDone)
	go continuouslyReportPosition(p, reportDone, func(position Position) {
		dispatcher.position <- position
	})

	http.HandleFunc("/events", getSSEHandler(dispatcher))
	httpServer := &http.Server{
		Addr: ":5000",
	}

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-gracefulShutdown
		log.Print("SIGTERM received. Shutdown process initiated\n")
		dispatcherDone <- nil
		reportDone <- nil
		ocr.Close()
		httpServer.Shutdown(context.Background())
	}()

	log.Fatal(httpServer.ListenAndServe())
}