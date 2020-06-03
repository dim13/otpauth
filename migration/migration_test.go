package migration

import "testing"

func TestConvert(t *testing.T) {
	const (
		migration = "otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC"
		want      = "otpauth://totp/Example:alice@google.com?issuer=Example&secret=JBSWY3DPEHPK3PXP"
	)
	got, err := Convert(migration)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) < 1 {
		t.Fatalf("got lengh %v, want 1", len(got))
	}
	if got[0] != want {
		t.Errorf("got %v, want %v", got[0], want)
	}
}
