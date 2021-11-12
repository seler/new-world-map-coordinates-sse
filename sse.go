package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/seler/new-world-map-coordinates-sse/position"
)

func getSSEHandler(dispatcher *Dispatcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Connection doesnot support streaming", http.StatusBadRequest)
			return
		}

		sseChannel := make(chan position.Position)
		dispatcher.clients = append(dispatcher.clients, sseChannel)

		d := make(chan interface{})
		defer close(d)

		for {
			select {
			case <-d:
				close(sseChannel)
				return
			case position := <-sseChannel:
				positionJSON, _ := json.Marshal(position)
				fmt.Fprintf(w, "data: %v \n\n", string(positionJSON))
				flusher.Flush()
			}
		}

	}

}
