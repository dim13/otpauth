// Google Authenticator migration decoder
//
// convert "otpauth-migration" links to plain "otpauth" links
//
package main

import (
	"flag"
	"fmt"
	"github.com/skip2/go-qrcode"
	"log"
	"net/url"

	"github.com/dim13/otpauth/migration"
)

func main() {
	link := flag.String("link", "", "migration link (required)")
	eval := flag.Bool("eval", false, "evaluate otps")
	qr := flag.Bool("qr", false, "qrcode generate otps")
	flag.Parse()

	u, err := url.Parse(*link)
	if err != nil {
		log.Fatal(err)
	}

	p, err := migration.Unmarshal(u)
	if err != nil {
		log.Fatal(err)
	}

	for idx, op := range p.OtpParameters {
		if *eval {
			fmt.Printf("%06d %s\n", op.Evaluate(), op.Name)
		} else {
			res := op.URL().String()
			fmt.Println(res)
			if *qr {
				fn :=fmt.Sprintf("qr_%d.png", idx)
				if err := qrcode.WriteFile(res, qrcode.Medium, 256, fn); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	if *qr {
		fmt.Println("Don't forgot delete qr_*.png from disk finally.")
	}
}
