package migration

import "testing"

func TestFileName(t *testing.T) {
	testCases := []struct {
		op   *Payload_OtpParameters
		want string
	}{
		{
			op: &Payload_OtpParameters{
				Name:   "Test",
				Issuer: "Issuer",
			},
			want: "Test_Issuer",
		},
		{
			op: &Payload_OtpParameters{
				Name:   "A/B:C.d",
				Issuer: "Issuer",
			},
			want: "A_B_C_d_Issuer",
		},
		{
			op: &Payload_OtpParameters{
				Name:   "A/../B:C.d",
				Issuer: "Issuer",
			},
			want: "A____B_C_d_Issuer",
		},
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
