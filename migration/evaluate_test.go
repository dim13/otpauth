package migration

import (
	"fmt"
	"log"
	"net/url"
	"time"
)

func ExampleEvaluate() {
	// fake time
	now = func() time.Time { return time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC) }
	u, err := url.Parse("otpauth-migration://offline?data=CjEKCkhlbGxvId6tvu8SGEV4YW1wbGU6YWxpY2VAZ29vZ2xlLmNvbRoHRXhhbXBsZTAC")
	if err != nil {
		log.Fatal(err)
	}
	p, err := Unmarshal(u)
	if err != nil {
		log.Fatal(err)
	}
	for _, op := range p.OtpParameters {
		fmt.Printf("%06d %s\n", op.Evaluate(), op.Name)
	}
	// Output:
	// 528064 Example:alice@google.com
}
