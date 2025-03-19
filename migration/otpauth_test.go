package migration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddBase32Padding(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"ABCDEFGH", "ABCDEFGH"},
		{"ABCDEF", "ABCDEF=="},
		{"A", "A======="},
		{"AB", "AB======"},
		{"ABC", "ABC====="},
		{"ABCD", "ABCD===="},
		{"ABCDE", "ABCDE==="},
		{"ABCDEFG", "ABCDEFG="},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := AddBase32Padding(tc.input)
			if result != tc.expected {
				t.Errorf("AddBase32Padding(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestCreateMigrationPayload(t *testing.T) {
	testCases := []struct {
		name         string
		urls         []string
		batchIndex   int
		batchCount   int
		expectedSize int
	}{
		{
			name: "valid URLs",
			urls: []string{
				"otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example",
				"otpauth://totp/Example:bob@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example",
			},
			batchIndex:   0,
			batchCount:   1,
			expectedSize: 2,
		},
		{
			name: "some invalid URLs",
			urls: []string{
				"otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example",
				"invalid://url",
				"otpauth://totp/Missing?issuer=Secret",
			},
			batchIndex:   0,
			batchCount:   1,
			expectedSize: 1,
		},
		{
			name: "different algorithms and digits",
			urls: []string{
				"otpauth://totp/Test1?secret=JBSWY3DPEHPK3PXP&algorithm=SHA1",
				"otpauth://totp/Test2?secret=JBSWY3DPEHPK3PXP&algorithm=SHA256&digits=8",
				"otpauth://totp/Test3?secret=JBSWY3DPEHPK3PXP&algorithm=SHA512",
				"otpauth://totp/Test4?secret=JBSWY3DPEHPK3PXP&algorithm=MD5",
			},
			batchIndex:   0,
			batchCount:   1,
			expectedSize: 4,
		},
		{
			name: "totp and hotp mixed",
			urls: []string{
				"otpauth://totp/Test1?secret=JBSWY3DPEHPK3PXP",
				"otpauth://hotp/Test2?secret=JBSWY3DPEHPK3PXP&counter=10",
			},
			batchIndex:   0,
			batchCount:   1,
			expectedSize: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := CreateMigrationPayload(tc.urls, tc.batchIndex, tc.batchCount)
			if err != nil {
				t.Fatalf("CreateMigrationPayload failed: %v", err)
			}

			if len(payload.OtpParameters) != tc.expectedSize {
				t.Errorf("Expected %d parameters, got %d", tc.expectedSize, len(payload.OtpParameters))
			}

			if payload.BatchIndex != int32(tc.batchIndex) {
				t.Errorf("Expected BatchIndex %d, got %d", tc.batchIndex, payload.BatchIndex)
			}

			if payload.BatchSize != int32(tc.batchCount) {
				t.Errorf("Expected BatchSize %d, got %d", tc.batchCount, payload.BatchSize)
			}
		})
	}
}

func TestProcessOtpauthFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "otpauth-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test-urls.txt")
	
	var testURLsLines []string
	for i := 1; i <= 15; i++ {
		testURLsLines = append(testURLsLines, 
			"otpauth://totp/Example:user"+string(rune('a'+i-1))+"@example.com?secret=JBSWY3DPEHPK3PXP&issuer=Example")
	}
	testURLsLines = append(testURLsLines, "invalid line")
	testURLsLines = append(testURLsLines, "not an otpauth line")
	testURLsLines = append(testURLsLines, "otpauth://totp/Missing?issuer=Secret")
	
	err = os.WriteFile(testFile, []byte(join(testURLsLines, "\n")), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	batchSize := 5
	batchPrefix := "test-batch"
	outputDir := filepath.Join(tempDir, "output")
	
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	err = ProcessOtpauthFile(testFile, outputDir, batchPrefix, batchSize)
	if err != nil {
		t.Fatalf("ProcessOtpauthFile failed: %v", err)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("Output directory was not created")
	}

	for i := 1; i <= 3; i++ {
		batchFile := filepath.Join(outputDir, "test-batch_"+itoa(i)+".png")
		if _, err := os.Stat(batchFile); os.IsNotExist(err) {
			t.Errorf("Batch %d file was not created", i)
		}
	}
}

func join(elements []string, separator string) string {
	var result string
	for i, element := range elements {
		if i > 0 {
			result += separator
		}
		result += element
	}
	return result
}

func itoa(i int) string {
	digits := "0123456789"
	if i == 0 {
		return "0"
	}
	var result string
	for i > 0 {
		result = string(digits[i%10]) + result
		i /= 10
	}
	return result
}