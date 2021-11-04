package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// grabber := NewScreenGrabber(0)
	grabber := NewFakeGrabber("test.png")
	ocr := NewOCRClient()
	ocr.Init()
	defer ocr.End()

	p := NewPositionService(grabber, ocr)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ticker := time.NewTicker(time.Second / 10)
		var waitForPosition sync.WaitGroup

		for range ticker.C {
			waitForPosition.Add(1)
			go func() {
				position := p.GetPosition()
				fmt.Printf("position >> %g <<", position)
				waitForPosition.Done()
			}()
			waitForPosition.Wait()
		}
		// never executed. how to close ticker and then politely close ocr client?
		wg.Done()
	}()
	wg.Wait()
}
