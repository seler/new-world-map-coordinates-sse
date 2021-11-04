package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type PositionReportCallback func(position Position)

func continuouslyReportPosition(p *PositionService, wg sync.WaitGroup, callback PositionReportCallback) {
	ticker := time.NewTicker(time.Second * 2)
	var waitForPosition sync.WaitGroup

	for range ticker.C {
		waitForPosition.Add(1)
		go func() {
			position := p.GetPosition()
			callback(position)
			waitForPosition.Done()
		}()
		waitForPosition.Wait()
	}
	// never executed. how to close ticker and then politely close ocr client?
	// - tested, just break out of a loop - use this when detecting interrupt signal - todo
	wg.Done()
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

func main() {
	// grabber := NewScreenGrabber(0)
	grabber := NewFakeGrabber("test.png")
	ocr := NewOCRClient()
	ocr.Init()
	defer ocr.End()

	p := NewPositionService(grabber, ocr)

	dispatcher := &Dispatcher{
		clients:  make([]chan Position, 0),
		position: make(chan Position),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go continuouslyReportPosition(p, wg, func(position Position) {
		dispatcher.position <- position
	})

	done := make(chan interface{})
	defer close(done)
	go dispatcher.dispatch(done)

	http.HandleFunc("/events", getSSEHandler(dispatcher))
	log.Fatal(http.ListenAndServe(":5000", nil))

	wg.Wait()
}
