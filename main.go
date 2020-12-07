// Google Authenticator migration decoder
//
// convert "otpauth-migration" links to plain "otpauth" links
//
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dim13/otpauth/migration"
)

func main() {
	link := flag.String("link", "", "migration link (required)")
	eval := flag.Bool("eval", false, "evaluate otps")
	qr := flag.Bool("qr", false, "generate QR-codes")
	http := flag.String("http", "", "serve http (e.g. localhost:6060)")
	flag.Parse()

	p, err := migration.UnmarshalURL(*link)
	if err != nil {
		log.Fatal(err)
	}

	for _, op := range p.OtpParameters {
		switch {
		case *eval:
			fmt.Printf("%06d %s\n", op.Evaluate(time.Now()), op.Name)
		case *qr:
			if err := op.WriteFile(); err != nil {
				log.Fatal(err)
			}
		default:
			fmt.Println(op.URL())
		}
	}

	if *http != "" {
		if err := serve(*http, p); err != nil {
			log.Fatal(err)
		}
	}
}
