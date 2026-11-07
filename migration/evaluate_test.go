package migration

import (
	"testing"
	"time"
)

func TestEvaluate(t *testing.T) {
	// fake time
	now = func() time.Time { return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC) }
	const testData = "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"
	p, err := UnmarshalURL(testData)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.OtpParameters) < 1 {
		t.Fatalf("got lengh %v, want 1", len(p.OtpParameters))
	}
	res := p.OtpParameters[0].Evaluate()
	if res != 528064 {
		t.Errorf("got %v", res)
	}
}

func TestEvaluateMD5(t *testing.T) {
	// An MD5 digest is 16 bytes, so a truncation offset above 12 leaves
	// fewer than 4 bytes to read.
	testCases := []struct {
		secret byte
		offset int
		want   string
	}{
		{secret: 1, offset: 1, want: "727980"},
		{secret: 6, offset: 12, want: "083708"},
		{secret: 13, offset: 13, want: "909933"},
		{secret: 0, offset: 14, want: "286446"},
		{secret: 37, offset: 15, want: "505023"},
	}
	for _, tc := range testCases {
		op := &Payload_OtpParameters{
			Secret:    []byte{tc.secret, 1, 2, 3},
			Algorithm: Payload_OtpParameters_ALGORITHM_MD5,
			Digits:    Payload_OtpParameters_DIGIT_COUNT_SIX,
			Type:      Payload_OtpParameters_OTP_TYPE_HOTP,
		}
		if got := op.EvaluateString(); got != tc.want {
			t.Errorf("offset %d: got %v, want %v", tc.offset, got, tc.want)
		}
	}
}
