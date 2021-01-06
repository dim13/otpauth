package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/dim13/otpauth/migration"
	"github.com/dim13/sse"
	"github.com/google/uuid"
)

//go:embed index.tmpl
var tmpl string

type otp struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Time float64   `json:"time"`
}

func serve(addr string, p *migration.Payload) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()
	log.Println("listen on", l.Addr())
	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, p)
	})
	for _, op := range p.OtpParameters {
		http.Handle("/"+op.UUID().String()+".png", op)
	}
	events := sse.New("otp", 100)
	http.Handle("/events", events)
	go func() {
		enc := json.NewEncoder(events)
		t := time.NewTicker(time.Second / 2)
		defer t.Stop()
		for range t.C {
			for _, op := range p.OtpParameters {
				enc.Encode(otp{
					ID:   op.UUID(),
					Code: op.EvaluateString(),
					Time: op.Second(),
				})
			}
		}
	}()
	return http.Serve(l, nil)
}
