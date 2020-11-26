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
