package sse

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Broker struct {
	event   string
	queue   int
	clients *sync.Map
}

func (b *Broker) Write(p []byte) (n int, err error) {
	b.clients.Range(func(key, value interface{}) bool {
		ch := key.(chan string)
		select {
		case ch <- string(p):
		default:
		}
		return true
	})
	return len(p), nil
}

func NewBroker(event string, queue int) *Broker {
	return &Broker{event: event, clients: new(sync.Map), queue: queue}
}

func (b Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "not a flusher", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := make(chan string, b.queue)
	defer close(ch)

	b.clients.Store(ch, nil)
	defer b.clients.Delete(ch)

	for data := range ch {
		select {
		case <-r.Context().Done():
			return
		default:
			if b.event != "" {
				fmt.Fprintf(w, "event: %s\n", b.event)
			}
			for _, s := range strings.Split(data, "\n") {
				fmt.Fprintf(w, "data: %s\n", s)
			}
			fmt.Fprintf(w, "\n")
			flusher.Flush()
		}
	}
}
