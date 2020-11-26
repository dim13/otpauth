package migration

import (
	"testing"
)

func TestConvert(t *testing.T) {
	const (
		testData = "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"
		want     = "otpauth://totp/Example:alice@google.com?issuer=Example&period=30&secret=JBSWY3DPEHPK3PXP"
	)
	p, err := UnmarshalURL(testData)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.OtpParameters) < 1 {
		t.Fatalf("got lengh %v, want 1", len(p.OtpParameters))
	}
	if p.OtpParameters[0].URL().String() != want {
		t.Errorf("got %v, want %v", p.OtpParameters[0].URL(), want)
	}
}
