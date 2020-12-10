package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/dim13/otpauth/migration"
	"github.com/dim13/otpauth/sse"
	"github.com/google/uuid"
)

const tmpl = `<!DOCTYPE html>
<html>
<header>
<style>
	body {
		font-family: 'Go', sans-serif;
		color: #0f0e09;
	}
	.code {
		color: #e76c1a;
	}
	.caption {
		color: #2f3795;
	}
	section {
		float: left;
		background-color: #e5e8f1;
		border: thick solid #847a73;
		border-radius: 1ex;
		margin: 1ex;
		padding: 1ex;
	}
</style>
<script>
	var events = new EventSource("/events");
	events.addEventListener("data", function(e) {
		var data = JSON.parse(e.data);
		var code = document.getElementById(data.ID).getElementsByClassName('code')[0];
		var time = document.getElementById(data.ID).getElementsByClassName('time')[0];
		code.innerHTML = data.Code;
		time.value = data.Time;
	});
</script>
</header>
<body>
{{range .OtpParameters}}
<section id="{{.UUID}}">
	<h4>{{.Title}}<h4>
	<figure><img src="{{.UUID}}.png" alt="{{.URL}}"></figure>
	<label class="code">{{.EvaluateString}}</label>
	<progress class="time" max="30"></progress>
</section>
{{end}}
</body>
</html>
`

type Code struct {
	ID   uuid.UUID
	Code string
	Time float64
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
	events := sse.NewBroker("data", 100)
	http.Handle("/events", events)
	go func() {
		enc := json.NewEncoder(events)
		t := time.NewTicker(time.Second / 2)
		defer t.Stop()
		for range t.C {
			for _, op := range p.OtpParameters {
				enc.Encode(Code{
					ID:   op.UUID(),
					Code: op.EvaluateString(),
					Time: op.Second(),
				})
			}
		}
	}()
	return http.Serve(l, nil)
}
