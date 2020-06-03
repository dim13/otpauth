package migration

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"google.golang.org/protobuf/proto"
)

//go:generate protoc --go_out=paths=source_relative:. migration.proto

var (
	otpTypes = map[Payload_OtpType]string{
		Payload_OTP_TYPE_HOTP: "hotp",
		Payload_OTP_TYPE_TOTP: "totp",
	}
	algorithms = map[Payload_Algorithm]string{
		Payload_ALGORITHM_SHA1:   "SHA1",
		Payload_ALGORITHM_SHA256: "SHA256",
		Payload_ALGORITHM_SHA512: "SHA512",
		Payload_ALGORITHM_MD5:    "MD5",
	}
	digitCounts = map[Payload_DigitCount]string{
		Payload_DIGIT_COUNT_SIX:   "6",
		Payload_DIGIT_COUNT_EIGHT: "8",
	}
)

func otpauth(op *Payload_OtpParameters) *url.URL {
	b := base32.StdEncoding.WithPadding(base32.NoPadding)
	v := make(url.Values)
	v.Add("secret", b.EncodeToString(op.Secret))
	v.Add("issuer", op.Issuer)
	if op.Algorithm != Payload_ALGORITHM_UNSPECIFIED {
		v.Add("algorithm", algorithms[op.Algorithm])
	}
	if op.Digits != Payload_DIGIT_COUNT_UNSPECIFIED {
		v.Add("digits", digitCounts[op.Digits])
	}
	if op.Counter > 0 {
		v.Add("counter", strconv.Itoa(int(op.Counter)))
	}
	return &url.URL{
		Scheme:   "otpauth",
		Host:     otpTypes[op.Type],
		Path:     op.Name,
		RawQuery: v.Encode(),
	}
}

func parseData(s string) ([]byte, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "otpauth-migration" {
		return nil, fmt.Errorf("unknown scheme: %s", u.Scheme)
	}
	if u.Host != "offline" {
		return nil, fmt.Errorf("unknown host: %s", u.Host)
	}
	data := u.Query().Get("data")
	return base64.StdEncoding.DecodeString(data)
}

// Convert otpauth-migration to otpauth links
func Convert(s string) ([]string, error) {
	data, err := parseData(s)
	if err != nil {
		return nil, err
	}
	var p Payload
	if err := proto.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	var ret []string
	for _, v := range p.OtpParameters {
		ret = append(ret, otpauth(v).String())
	}
	return ret, nil
}
