// Google Authenticator migration decoder
//
// convert "otpauth-migration" links to plain "otpauth" links
package main

import (
	"encoding/base32"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dim13/otpauth/migration"
	"google.golang.org/protobuf/proto"
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

// processOtpauthFile reads a file containing otpauth:// URLs (one per line),
// groups them in batches, creates migration payloads for each batch,
// and generates QR codes.
func processOtpauthFile(filePath, workdir, batchPrefix string, batchSize int) error {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	// Split the content by lines and filter empty lines
	var urls []string
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasPrefix(line, "otpauth://") {
			urls = append(urls, line)
		}
	}

	if len(urls) == 0 {
		return fmt.Errorf("no valid otpauth:// URLs found in file")
	}

	// Process URLs in batches
	batchCount := (len(urls) + batchSize - 1) / batchSize // Ceiling division

	fmt.Printf("Found %d otpauth URLs, creating %d batches\n", len(urls), batchCount)

	var processedBatches int
	var totalProcessedURLs int

	for batchIdx := 0; batchIdx < batchCount; batchIdx++ {
		start := batchIdx * batchSize
		end := (batchIdx + 1) * batchSize
		if end > len(urls) {
			end = len(urls)
		}

		batchUrls := urls[start:end]
		payload, err := createMigrationPayload(batchUrls, batchIdx, batchCount)
		if err != nil {
			fmt.Printf("Error in batch %d: %v\n", batchIdx+1, err)
			continue
		}

		if len(payload.OtpParameters) == 0 {
			fmt.Printf("Skipping batch %d: No valid OTP parameters found\n", batchIdx+1)
			continue
		}

		// Marshal the payload to protobuf
		data, err := proto.Marshal(payload)
		if err != nil {
			fmt.Printf("Error marshaling payload for batch %d: %v\n", batchIdx+1, err)
			continue
		}

		// Create QR code for the batch
		fileName := fmt.Sprintf("%s_%d.png", batchPrefix, batchIdx+1)
		filePath := filepath.Join(workdir, fileName)
		if err := migration.PNG(filePath, migration.URL(data)); err != nil {
			fmt.Printf("Error generating QR code for batch %d: %v\n", batchIdx+1, err)
			continue
		}

		processedBatches++
		totalProcessedURLs += len(payload.OtpParameters)
		fmt.Printf("Created batch %d QR code: %s (with %d valid URLs)\n",
			batchIdx+1, filePath, len(payload.OtpParameters))
	}

	// Print summary
	fmt.Printf("\n=== SUMMARY ===\n")
	fmt.Printf("Total URLs found: %d\n", len(urls))
	fmt.Printf("Valid URLs processed: %d\n", totalProcessedURLs)
	fmt.Printf("Skipped URLs: %d\n", len(urls) - totalProcessedURLs)
	fmt.Printf("Successfully processed batches: %d out of %d\n", processedBatches, batchCount)

	return nil
}

// createMigrationPayload converts a list of otpauth:// URLs into a migration.Payload
// Returns the payload and a slice of invalid URLs
func createMigrationPayload(urls []string, batchIndex, batchCount int) (*migration.Payload, error) {
	payload := &migration.Payload{
		Version:    1,
		BatchSize:  int32(batchCount),
		BatchIndex: int32(batchIndex),
		BatchId:    1, // Using a default batch ID
	}

	for _, urlStr := range urls {
		// Parse the otpauth URL
		u, err := url.Parse(urlStr)
		if err != nil {
			fmt.Printf("SKIPPING - Invalid URL format: %s\nError: %v\n", urlStr, err)
			continue
		}

		if u.Scheme != "otpauth" {
			fmt.Printf("SKIPPING - Invalid URL scheme: %s\nURL: %s\n", u.Scheme, urlStr)
			continue
		}

		// Extract parameters
		values := u.Query()
		secretBase32 := values.Get("secret")
		if secretBase32 == "" {
			fmt.Printf("SKIPPING - Missing secret parameter in URL: %s\n", urlStr)
			continue
		}

		// Add padding if needed (Google Authenticator uses NoPadding)
		secretBase32 = addBase32Padding(secretBase32)

		// Try to decode secret from base32
		secret, err := base32.StdEncoding.DecodeString(secretBase32)
		if err != nil {
			// If Base32 decoding fails, use the secret as plain text
			fmt.Printf("WARNING - Invalid Base32 secret in URL: %s\nError: %v\n", urlStr, err)
			fmt.Printf("Using secret as plain text instead of Base32-encoded data\n")
			secret = []byte(secretBase32) // Use the raw string as the secret
		}

		// Create OtpParameters
		param := &migration.Payload_OtpParameters{
			Secret: secret,
			Name:   u.Path,
			Issuer: values.Get("issuer"),
		}

		// Set algorithm
		switch strings.ToUpper(values.Get("algorithm")) {
		case "SHA1", "":
			param.Algorithm = migration.Payload_OtpParameters_ALGORITHM_SHA1
		case "SHA256":
			param.Algorithm = migration.Payload_OtpParameters_ALGORITHM_SHA256
		case "SHA512":
			param.Algorithm = migration.Payload_OtpParameters_ALGORITHM_SHA512
		case "MD5":
			param.Algorithm = migration.Payload_OtpParameters_ALGORITHM_MD5
		default:
			param.Algorithm = migration.Payload_OtpParameters_ALGORITHM_SHA1
		}

		// Set digit count
		switch values.Get("digits") {
		case "6", "":
			param.Digits = migration.Payload_OtpParameters_DIGIT_COUNT_SIX
		case "8":
			param.Digits = migration.Payload_OtpParameters_DIGIT_COUNT_EIGHT
		default:
			param.Digits = migration.Payload_OtpParameters_DIGIT_COUNT_SIX
		}

		// Set OTP type
		switch u.Host {
		case "hotp":
			param.Type = migration.Payload_OtpParameters_OTP_TYPE_HOTP
			counter, _ := strconv.ParseUint(values.Get("counter"), 10, 64)
			param.Counter = counter
		case "totp", "":
			param.Type = migration.Payload_OtpParameters_OTP_TYPE_TOTP
		default:
			param.Type = migration.Payload_OtpParameters_OTP_TYPE_TOTP
		}

		// If path starts with a slash, remove it
		if strings.HasPrefix(param.Name, "/") {
			param.Name = param.Name[1:]
		}

		// If name starts with the issuer name followed by a colon, remove it
		if param.Issuer != "" && strings.HasPrefix(param.Name, param.Issuer+":") {
			param.Name = param.Name[len(param.Issuer)+1:]
		}

		payload.OtpParameters = append(payload.OtpParameters, param)
	}

	return payload, nil
}

// addBase32Padding adds padding to a base32 string if needed
func addBase32Padding(s string) string {
	if len(s)%8 == 0 {
		return s
	}

	padLen := 8 - (len(s) % 8)
	return s + strings.Repeat("=", padLen)
}

func main() {
	var (
		link        = flag.String("link", "", "migration link (required)")
		workdir     = flag.String("workdir", "", "working directory")
		http        = flag.String("http", "", "serve http (e.g. localhost:6060)")
		eval        = flag.Bool("eval", false, "evaluate otps")
		qr          = flag.Bool("qr", false, "generate QR-codes (optauth://)")
		rev         = flag.Bool("rev", false, "reverse QR-code (otpauth-migration://)")
		info        = flag.Bool("info", false, "display batch info")
		inputFile   = flag.String("file", "", "input file with otpauth:// URLs (one per line)")
		batchPrefix = flag.String("batch-prefix", "batch", "prefix for batch QR code filenames")
		batchSize   = flag.Int("batch-size", 7, "number of URLs to include in each batch (default: 7)")
	)
	flag.Parse()

	if *workdir != "" {
		if err := os.MkdirAll(*workdir, 0700); err != nil {
			log.Fatal("error creating working directory: ", err)
		}
	}

	// Handle input file with otpauth URLs
	if *inputFile != "" {
		if err := processOtpauthFile(*inputFile, *workdir, *batchPrefix, *batchSize); err != nil {
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
