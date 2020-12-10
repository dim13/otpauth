package migration

import "testing"

func TestFileName(t *testing.T) {
	testCases := []struct {
		op   *Payload_OtpParameters
		want string
	}{
		{&Payload_OtpParameters{Name: "Test"}, "Test"},
		{&Payload_OtpParameters{Name: "A/B:C.d"}, "A_B:C.d"},
		{&Payload_OtpParameters{Name: "A/../B:C.d"}, "A_.._B:C.d"},
	}

	for _, tc := range testCases {
		t.Run(tc.op.Name, func(t *testing.T) {
			got := tc.op.FileName()
			if got != tc.want {
				t.Errorf("got %v; want %v", got, tc.want)
			}
		})
	}

}
