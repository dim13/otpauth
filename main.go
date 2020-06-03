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
	mig := flag.String("link", "", "migration link")
	flag.Parse()
	u, err := migration.Convert(*mig)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range u {
		fmt.Println(v)
	}
}
