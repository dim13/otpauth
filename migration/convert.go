package migration

import (
	"encoding/base32"
	"fmt"
	"net/url"
)

var (
	typeString = map[Payload_OtpParameters_OtpType]string{
		Payload_OtpParameters_OTP_TYPE_HOTP: "hotp",
		Payload_OtpParameters_OTP_TYPE_TOTP: "totp",
	}
	algString = map[Payload_OtpParameters_Algorithm]string{
		Payload_OtpParameters_ALGORITHM_SHA1:   "SHA1",
		Payload_OtpParameters_ALGORITHM_SHA256: "SHA256",
		Payload_OtpParameters_ALGORITHM_SHA512: "SHA512",
		Payload_OtpParameters_ALGORITHM_MD5:    "MD5",
	}
	digitsString = map[Payload_OtpParameters_DigitCount]string{
		Payload_OtpParameters_DIGIT_COUNT_SIX:   "6",
		Payload_OtpParameters_DIGIT_COUNT_EIGHT: "8",
	}
)

// SecretString returns Secret as a base32 encoded String
func (op *Payload_OtpParameters) SecretString() string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(op.Secret)
}

// SecretTuples returns Secret as a base32 string splitted into tuples of 4
func (op *Payload_OtpParameters) SecretTuples() []string {
	secret := op.SecretString()
	// pad secret to multiple of 4×4
	tuples := make([]string, len(secret)/4)
	for i := range tuples {
		tuples[i] = secret[4*i : 4*i+4]
	}
	return tuples
}

// URL of otp parameters
func (op *Payload_OtpParameters) URL() *url.URL {
	v := make(url.Values)
	// required
	v.Add("secret", op.SecretString())
	// strongly recommended
	if op.Issuer != "" {
		v.Add("issuer", op.Issuer)
	}
	// optional
	if op.Algorithm != Payload_OtpParameters_ALGORITHM_UNSPECIFIED {
		v.Add("algorithm", algString[op.Algorithm])
	}
	// optional
	if op.Digits != Payload_OtpParameters_DIGIT_COUNT_UNSPECIFIED {
		v.Add("digits", digitsString[op.Digits])
	}
	// required if type is hotp
	if op.Type == Payload_OtpParameters_OTP_TYPE_HOTP {
		v.Add("counter", fmt.Sprint(op.Counter))
	}
	// optional if type is totp
	if op.Type == Payload_OtpParameters_OTP_TYPE_TOTP {
		v.Add("period", "30") // default value
	}
	return &url.URL{
		Scheme:   "otpauth",
		Host:     typeString[op.Type],
		Path:     op.Name,
		RawQuery: v.Encode(),
	}
}
