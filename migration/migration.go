package migration

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative migration.proto

var algorithmHash = []func() hash.Hash{
	Payload_OtpParameters_ALGORITHM_UNSPECIFIED: sha1.New,
	Payload_OtpParameters_ALGORITHM_SHA1:        sha1.New,
	Payload_OtpParameters_ALGORITHM_SHA256:      sha256.New,
	Payload_OtpParameters_ALGORITHM_SHA512:      sha512.New,
	Payload_OtpParameters_ALGORITHM_MD5:         md5.New,
}

func (x Payload_OtpParameters_Algorithm) Hash() func() hash.Hash {
	return algorithmHash[x]
}

var algorithmNames = []string{
	Payload_OtpParameters_ALGORITHM_UNSPECIFIED: "SHA1",
	Payload_OtpParameters_ALGORITHM_SHA1:        "SHA1",
	Payload_OtpParameters_ALGORITHM_SHA256:      "SHA256",
	Payload_OtpParameters_ALGORITHM_SHA512:      "SHA512",
	Payload_OtpParameters_ALGORITHM_MD5:         "MD5",
}

func (x Payload_OtpParameters_Algorithm) Name() string {
	return algorithmNames[x]
}

var digitCount = []int{
	Payload_OtpParameters_DIGIT_COUNT_UNSPECIFIED: 6,
	Payload_OtpParameters_DIGIT_COUNT_SIX:         6,
	Payload_OtpParameters_DIGIT_COUNT_EIGHT:       8,
}

func (x Payload_OtpParameters_DigitCount) Count() int {
	return digitCount[x]
}

var otpTypeFunc = []func(*Payload_OtpParameters) uint64{
	Payload_OtpParameters_OTP_TYPE_UNSPECIFIED: totp,
	Payload_OtpParameters_OTP_TYPE_HOTP:        hotp,
	Payload_OtpParameters_OTP_TYPE_TOTP:        totp,
}

func (x Payload_OtpParameters_OtpType) Count(op *Payload_OtpParameters) uint64 {
	return otpTypeFunc[x](op)
}

var otpTypeNames = []string{
	Payload_OtpParameters_OTP_TYPE_UNSPECIFIED: "totp",
	Payload_OtpParameters_OTP_TYPE_HOTP:        "hotp",
	Payload_OtpParameters_OTP_TYPE_TOTP:        "totp",
}

func (x Payload_OtpParameters_OtpType) Name() string {
	return otpTypeNames[x]
}
