package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/dim13/otpauth/migration"
	"github.com/google/uuid"
)

//go:embed static
var static embed.FS

type otp struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Time float64   `json:"time"`
}

func eventStream(event string, p *migration.Payload) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "not a flusher", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for range t.C {
			select {
			case <-r.Context().Done():
				return
			default:
				for _, op := range p.OtpParameters {
					b, _ := json.Marshal(otp{
						ID:   op.UUID(),
						Code: op.EvaluateString(),
						Time: op.Seconds(),
					})
					fmt.Fprintf(w, "event: %s\r\n", event)
					fmt.Fprintf(w, "data: %s\r\n\r\n", string(b))
				}
				flusher.Flush()
			}
		}
	}
}

func indexHandler(t *template.Template, p *migration.Payload) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, p); err != nil {
			log.Println("execute template:", err)
		}
	}
}

func serve(addr string, p *migration.Payload) error {
	t, err := template.ParseFS(static, "static/index.html")
	if err != nil {
		return err
	}
	http.Handle("/", indexHandler(t, p))
	for _, op := range p.OtpParameters {
		http.Handle("/"+op.UUID().String()+".png", op)
	}
	http.Handle("/events", eventStream("otp", p))
	http.Handle("/static/", http.FileServer(http.FS(static)))
	log.Println("listen on", addr)
	return http.ListenAndServe(addr, nil)
}
