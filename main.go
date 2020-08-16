// Google Authenticator migration decoder
//
// convert "otpauth-migration" links to plain "otpauth" links
//
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/dim13/otpauth/migration"
)

func main() {
	link := flag.String("link", "", "migration link (required)")
	eval := flag.Bool("eval", false, "evaluate otps")
	flag.Parse()

	u, err := url.Parse(*link)
	if err != nil {
		log.Fatal(err)
	}

	p, err := migration.Unmarshal(u)
	if err != nil {
		log.Fatal(err)
	}

	for _, op := range p.OtpParameters {
		if *eval {
			fmt.Printf("%06d %s\n", op.Evaluate(), op.Name)
		} else {
			fmt.Println(op.URL())
		}
	}
}
