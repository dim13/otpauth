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
	"github.com/skip2/go-qrcode"
)

func main() {
	link := flag.String("link", "", "migration link (required)")
	eval := flag.Bool("eval", false, "evaluate otps")
	qr := flag.Bool("qr", false, "generate QR-codes")
	flag.Parse()

	p, err := migration.UnmarshalURL(*link)
	if err != nil {
		log.Fatal(err)
	}

	for i, op := range p.OtpParameters {
		switch {
		case *eval:
			fmt.Printf("%06d %s\n", op.Evaluate(), op.Name)
		case *qr:
			fname := fmt.Sprintf("%d_%s.png", i+1, op.CleanName())
			fmt.Println("write", fname)
			err := qrcode.WriteFile(op.URL().String(), qrcode.Medium, 256, fname)
			if err != nil {
				log.Fatal(err)
			}
		default:
			fmt.Println(op.URL())
		}
	}
}
