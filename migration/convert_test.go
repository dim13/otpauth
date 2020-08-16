package migration

import (
	"net/url"
	"testing"
)

func TestConvert(t *testing.T) {
	const (
		migration = "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"
		want      = "otpauth://totp/Example:alice@google.com?issuer=Example&period=30&secret=JBSWY3DPEHPK3PXP"
	)
	u, err := url.Parse(migration)
	if err != nil {
		t.Fatal(err)
	}
	p, err := Unmarshal(u)
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
