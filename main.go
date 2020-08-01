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
	mig := flag.String("link", "", "migration link")
	eval := flag.Bool("eval", false, "evaluate otps")
	flag.Parse()
	u, err := url.Parse(*mig)
	if err != nil {
		log.Fatal(err)
	}
	if *eval {
		migration.Evaluate(u)
	} else {
		p, err := migration.Convert(u)
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range p {
			fmt.Println(v)
		}
	}
}
