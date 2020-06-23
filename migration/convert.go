package migration

import (
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"

	"google.golang.org/protobuf/proto"
)

//go:generate protoc --go_out=paths=source_relative:. migration.proto

// Errors
var (
	ErrUnkown = errors.New("unknown")
)

var (
	typeString = map[Payload_OtpType]string{
		Payload_OTP_TYPE_HOTP: "hotp",
		Payload_OTP_TYPE_TOTP: "totp",
	}
	algString = map[Payload_Algorithm]string{
		Payload_ALGORITHM_SHA1:   "SHA1",
		Payload_ALGORITHM_SHA256: "SHA256",
		Payload_ALGORITHM_SHA512: "SHA512",
		Payload_ALGORITHM_MD5:    "MD5",
	}
	digitsString = map[Payload_DigitCount]string{
		Payload_DIGIT_COUNT_SIX:   "6",
		Payload_DIGIT_COUNT_EIGHT: "8",
	}
)

// URL of otp parameters
func (op *Payload_OtpParameters) URL() *url.URL {
	b := base32.StdEncoding.WithPadding(base32.NoPadding)
	v := make(url.Values)
	// required
	v.Add("secret", b.EncodeToString(op.Secret))
	// strongly recommended
	if op.Issuer != "" {
		v.Add("issuer", op.Issuer)
	}
	// optional
	if op.Algorithm != Payload_ALGORITHM_UNSPECIFIED {
		v.Add("algorithm", algString[op.Algorithm])
	}
	// optional
	if op.Digits != Payload_DIGIT_COUNT_UNSPECIFIED {
		v.Add("digits", digitsString[op.Digits])
	}
	// required if type is hotp
	if op.Type == Payload_OTP_TYPE_HOTP {
		v.Add("counter", fmt.Sprint(op.Counter))
	}
	// optional if type is totp
	if op.Type == Payload_OTP_TYPE_TOTP {
		v.Add("period", "30") // default value
	}
	return &url.URL{
		Scheme:   "otpauth",
		Host:     typeString[op.Type],
		Path:     op.Name,
		RawQuery: v.Encode(),
	}
}

func dataQuery(u *url.URL) ([]byte, error) {
	if u.Scheme != "otpauth-migration" {
		return nil, fmt.Errorf("scheme %s: %w", u.Scheme, ErrUnkown)
	}
	if u.Host != "offline" {
		return nil, fmt.Errorf("host %s: %w", u.Host, ErrUnkown)
	}
	data := u.Query().Get("data")
	return base64.StdEncoding.DecodeString(data)
}

// Convert otpauth-migration URL to otpauth URL
func Convert(u *url.URL) ([]*url.URL, error) {
	data, err := dataQuery(u)
	if err != nil {
		return nil, err
	}
	var p Payload
	if err := proto.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	var ret []*url.URL
	for _, op := range p.OtpParameters {
		ret = append(ret, op.URL())
	}
	return ret, nil
}
