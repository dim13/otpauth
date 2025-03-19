// Google Authenticator migration decoder
//
// convert "otpauth-migration" links to plain "otpauth" links
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/dim13/otpauth/migration"
)

const (
	cacheFilename = "migration.bin"
	revFile       = "otpauth-migration.png"
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
		link                    = flag.String("link", "", "migration link (required)")
		workdir                 = flag.String("workdir", "", "working directory")
		http                    = flag.String("http", "", "serve http (e.g. localhost:6060)")
		eval                    = flag.Bool("eval", false, "evaluate otps")
		qr                      = flag.Bool("qr", false, "generate QR-codes (optauth://)")
		rev                     = flag.Bool("rev", false, "reverse QR-code (otpauth-migration://)")
		info                    = flag.Bool("info", false, "display batch info")
		otpauthUrlsFile         = flag.String("file", "", "input file with otpauth:// URLs (one per line)")
		migrationBatchImgPrefix = flag.String("migration-batch-img-prefix", "batch", "prefix for batch QR code filenames")
		migrationBatchSize      = flag.Int("migration-batch-size", 7, "number of URLs to include in each batch (default: 7)")
	)
	flag.Parse()

	if *workdir != "" {
		if err := os.MkdirAll(*workdir, 0700); err != nil {
			log.Fatal("error creating working directory: ", err)
		}
	}

	if *otpauthUrlsFile != "" {
		if err := migration.ProcessOtpauthFile(*otpauthUrlsFile, *workdir, *migrationBatchImgPrefix, *migrationBatchSize); err != nil {
			log.Fatal("processing input file: ", err)
		}
		return
	}

	cacheFile := filepath.Join(*workdir, cacheFilename)
	data, err := migrationData(cacheFile, *link)
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
			fileName := op.FileName() + ".png"
			qrFile := filepath.Join(*workdir, fileName)
			if err := migration.PNG(qrFile, op.URL()); err != nil {
				log.Fatal("write file: ", err)
			}
		}
	case *rev:
		revFile := filepath.Join(*workdir, revFile)
		if err := migration.PNG(revFile, migration.URL(data)); err != nil {
			log.Fatal(err)
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
