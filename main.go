// Google Authenticator migration decoder
//
// convert "otpauth-migration" links to plain "otpauth" links
//
package main

import (
	"flag"
	"fmt"
	"log"

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

	switch {
	case *http != "":
		if err := serve(*http, p); err != nil {
			log.Fatal(err)
		}
	case *qr:
		for _, op := range p.OtpParameters {
			if err := op.WriteFile(op.FileName() + ".png"); err != nil {
				log.Fatal(err)
			}
		}
	case *eval:
		for _, op := range p.OtpParameters {
			fmt.Printf("%06d %s\n", op.Evaluate(), op.Name)
		}
	default:
		for _, op := range p.OtpParameters {
			fmt.Println(op.URL())
		}
	}
}
