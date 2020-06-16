package migration

import (
	"net/url"
	"time"
)

func ExampleEvaluate() {
	// fake time
	now = func() time.Time { return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC) }
	u, _ := url.Parse("otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC")
	Evaluate(u)
	// Output:
	// 528064 Example:alice@google.com
}
