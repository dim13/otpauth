package migration

import (
	"encoding/base32"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
)

// ProcessOtpauthFile reads a file containing otpauth:// URLs (one per line),
// groups them in batches, creates migration payloads for each batch,
// and generates QR codes.
func ProcessOtpauthFile(filePath, workdir, batchPrefix string, batchSize int) error {
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
		payload, err := CreateMigrationPayload(batchUrls, batchIdx, batchCount)
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
		if err := PNG(filePath, URL(data)); err != nil {
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
	fmt.Printf("Skipped URLs: %d\n", len(urls)-totalProcessedURLs)
	fmt.Printf("Successfully processed batches: %d out of %d\n", processedBatches, batchCount)

	return nil
}

// CreateMigrationPayload converts a list of otpauth:// URLs into a migration.Payload
// Returns the payload and any error encountered
func CreateMigrationPayload(urls []string, batchIndex, batchCount int) (*Payload, error) {
	payload := &Payload{
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
		secretBase32 = AddBase32Padding(secretBase32)

		// Try to decode secret from base32
		secret, err := base32.StdEncoding.DecodeString(secretBase32)
		if err != nil {
			// If Base32 decoding fails, use the secret as plain text
			fmt.Printf("WARNING - Invalid Base32 secret in URL: %s\nError: %v\n", urlStr, err)
			fmt.Printf("Using secret as plain text instead of Base32-encoded data\n")
			secret = []byte(secretBase32) // Use the raw string as the secret
		}

		// Create OtpParameters
		param := &Payload_OtpParameters{
			Secret: secret,
			Name:   u.Path,
			Issuer: values.Get("issuer"),
		}

		// Set algorithm
		switch strings.ToUpper(values.Get("algorithm")) {
		case "SHA1", "":
			param.Algorithm = Payload_OtpParameters_ALGORITHM_SHA1
		case "SHA256":
			param.Algorithm = Payload_OtpParameters_ALGORITHM_SHA256
		case "SHA512":
			param.Algorithm = Payload_OtpParameters_ALGORITHM_SHA512
		case "MD5":
			param.Algorithm = Payload_OtpParameters_ALGORITHM_MD5
		default:
			param.Algorithm = Payload_OtpParameters_ALGORITHM_SHA1
		}

		// Set digit count
		switch values.Get("digits") {
		case "6", "":
			param.Digits = Payload_OtpParameters_DIGIT_COUNT_SIX
		case "8":
			param.Digits = Payload_OtpParameters_DIGIT_COUNT_EIGHT
		default:
			param.Digits = Payload_OtpParameters_DIGIT_COUNT_SIX
		}

		// Set OTP type
		switch u.Host {
		case "hotp":
			param.Type = Payload_OtpParameters_OTP_TYPE_HOTP
			counter, _ := strconv.ParseUint(values.Get("counter"), 10, 64)
			param.Counter = counter
		case "totp", "":
			param.Type = Payload_OtpParameters_OTP_TYPE_TOTP
		default:
			param.Type = Payload_OtpParameters_OTP_TYPE_TOTP
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

// AddBase32Padding adds padding to a base32 string if needed
func AddBase32Padding(s string) string {
	if len(s)%8 == 0 {
		return s
	}

	padLen := 8 - (len(s) % 8)
	return s + strings.Repeat("=", padLen)
}