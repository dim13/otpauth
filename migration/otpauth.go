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

func ProcessOtpauthFile(filePath, workdir, batchPrefix string, batchSize int) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

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

	batchCount := (len(urls) + batchSize - 1) / batchSize

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

		data, err := proto.Marshal(payload)
		if err != nil {
			fmt.Printf("Error marshaling payload for batch %d: %v\n", batchIdx+1, err)
			continue
		}

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

	fmt.Printf("\n=== SUMMARY ===\n")
	fmt.Printf("Total URLs found: %d\n", len(urls))
	fmt.Printf("Valid URLs processed: %d\n", totalProcessedURLs)
	fmt.Printf("Skipped URLs: %d\n", len(urls)-totalProcessedURLs)
	fmt.Printf("Successfully processed batches: %d out of %d\n", processedBatches, batchCount)

	return nil
}

func CreateMigrationPayload(urls []string, batchIndex, batchCount int) (*Payload, error) {
	payload := &Payload{
		Version:    1,
		BatchSize:  int32(batchCount),
		BatchIndex: int32(batchIndex),
		BatchId:    1,
	}

	for _, urlStr := range urls {
		u, err := url.Parse(urlStr)
		if err != nil {
			fmt.Printf("SKIPPING - Invalid URL format: %s\nError: %v\n", urlStr, err)
			continue
		}

		if u.Scheme != "otpauth" {
			fmt.Printf("SKIPPING - Invalid URL scheme: %s\nURL: %s\n", u.Scheme, urlStr)
			continue
		}

		values := u.Query()
		secretBase32 := values.Get("secret")
		if secretBase32 == "" {
			fmt.Printf("SKIPPING - Missing secret parameter in URL: %s\n", urlStr)
			continue
		}

		secretBase32 = AddBase32Padding(secretBase32)

		secret, err := base32.StdEncoding.DecodeString(secretBase32)
		if err != nil {
			fmt.Printf("WARNING - Invalid Base32 secret in URL: %s\nError: %v\n", urlStr, err)
			fmt.Printf("Using secret as plain text instead of Base32-encoded data\n")
			secret = []byte(secretBase32)
		}

		param := &Payload_OtpParameters{
			Secret: secret,
			Name:   u.Path,
			Issuer: values.Get("issuer"),
		}

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

		switch values.Get("digits") {
		case "6", "":
			param.Digits = Payload_OtpParameters_DIGIT_COUNT_SIX
		case "8":
			param.Digits = Payload_OtpParameters_DIGIT_COUNT_EIGHT
		default:
			param.Digits = Payload_OtpParameters_DIGIT_COUNT_SIX
		}

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

		param.Name = strings.TrimPrefix(param.Name, "/")

		if param.Issuer != "" && strings.HasPrefix(param.Name, param.Issuer+":") {
			param.Name = param.Name[len(param.Issuer)+1:]
		}

		payload.OtpParameters = append(payload.OtpParameters, param)
	}

	return payload, nil
}

func AddBase32Padding(s string) string {
	if len(s)%8 == 0 {
		return s
	}

	padLen := 8 - (len(s) % 8)
	return s + strings.Repeat("=", padLen)
}
