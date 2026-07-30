package migration

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative migration.proto

func (x Payload_OtpParameters_Algorithm) Hash() func() hash.Hash {
	switch x {
	case Payload_OtpParameters_ALGORITHM_SHA1:
		return sha1.New
	case Payload_OtpParameters_ALGORITHM_SHA256:
		return sha256.New
	case Payload_OtpParameters_ALGORITHM_SHA512:
		return sha512.New
	case Payload_OtpParameters_ALGORITHM_MD5:
		return md5.New
	default:
		return sha1.New
	}
}

func (x Payload_OtpParameters_Algorithm) Name() string {
	switch x {
	case Payload_OtpParameters_ALGORITHM_SHA1:
		return "SHA1"
	case Payload_OtpParameters_ALGORITHM_SHA256:
		return "SHA256"
	case Payload_OtpParameters_ALGORITHM_SHA512:
		return "SHA512"
	case Payload_OtpParameters_ALGORITHM_MD5:
		return "MD5"
	default:
		return "SHA1"
	}
}

func (x Payload_OtpParameters_DigitCount) Count() int {
	switch x {
	case Payload_OtpParameters_DIGIT_COUNT_SIX:
		return 6
	case Payload_OtpParameters_DIGIT_COUNT_EIGHT:
		return 8
	default:
		return 6
	}
}

func (x Payload_OtpParameters_OtpType) Count(op *Payload_OtpParameters) uint64 {
	switch x {
	case Payload_OtpParameters_OTP_TYPE_HOTP:
		return hotp(op)
	case Payload_OtpParameters_OTP_TYPE_TOTP:
		return totp()
	default:
		return totp()
	}
}

func (x Payload_OtpParameters_OtpType) Name() string {
	switch x {
	case Payload_OtpParameters_OTP_TYPE_HOTP:
		return "hotp"
	case Payload_OtpParameters_OTP_TYPE_TOTP:
		return "totp"
	default:
		return "totp"
	}
}
