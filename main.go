// Google Authenticator migration decoder
//
// convert "otpauth-migration" links to plain "otpauth" links
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dim13/otpauth/migration"
)

func migrationData(fname, link string) ([]byte, error) {
	if link == "" {
		// read from cache
		return os.ReadFile(fname)
	}
	data, err := migration.Data(link)
	if err != nil {
		return nil, err
	}
	// write to cache
	return data, os.WriteFile(fname, data, 0600)
}

func main() {
	var (
		link  = flag.String("link", "", "migration link (required)")
		cache = flag.String("cache", "migration.bin", "cache file")
		http  = flag.String("http", "", "serve http (e.g. localhost:6060)")
		eval  = flag.Bool("eval", false, "evaluate otps")
		qr    = flag.Bool("qr", false, "generate QR-codes")
		info  = flag.Bool("info", false, "display batch info")
	)
	flag.Parse()

	data, err := migrationData(*cache, *link)
	if err != nil {
		log.Fatal("-link parameter or cache file missing: ", err)
	}

	p, err := migration.Unmarshal(data)
	if err != nil {
		log.Fatal("decode data: ", err)
	}

	switch {
	case *http != "":
		if err := serve(*http, p); err != nil {
			log.Fatal("serve http: ", err)
		}
	case *qr:
		for _, op := range p.OtpParameters {
			if err := op.WriteFile(op.FileName() + ".png"); err != nil {
				log.Fatal("write file: ", err)
			}
		}
	case *eval:
		for _, op := range p.OtpParameters {
			fmt.Printf("%06d %s\n", op.Evaluate(), op.Name)
		}
	case *info:
		fmt.Println("version", p.Version)
		fmt.Println("batch size", p.BatchSize)
		fmt.Println("batch index", p.BatchIndex)
		fmt.Println("batch id", p.BatchId)
	default:
		for _, op := range p.OtpParameters {
			fmt.Println(op.URL())
		}
	}
}
