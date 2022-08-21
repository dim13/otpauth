package migration

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"time"
)

var (
	hashFunc = map[Payload_OtpParameters_Algorithm]func() hash.Hash{
		Payload_OtpParameters_ALGORITHM_UNSPECIFIED: sha1.New, // default
		Payload_OtpParameters_ALGORITHM_SHA1:        sha1.New,
		Payload_OtpParameters_ALGORITHM_SHA256:      sha256.New,
		Payload_OtpParameters_ALGORITHM_SHA512:      sha512.New,
		Payload_OtpParameters_ALGORITHM_MD5:         md5.New,
	}
	digitCount = map[Payload_OtpParameters_DigitCount]int{
		Payload_OtpParameters_DIGIT_COUNT_UNSPECIFIED: 6, // default
		Payload_OtpParameters_DIGIT_COUNT_SIX:         6,
		Payload_OtpParameters_DIGIT_COUNT_EIGHT:       8,
	}
	countFunc = map[Payload_OtpParameters_OtpType]func(*Payload_OtpParameters) uint64{
		Payload_OtpParameters_OTP_TYPE_UNSPECIFIED: totp, // default
		Payload_OtpParameters_OTP_TYPE_HOTP:        hotp,
		Payload_OtpParameters_OTP_TYPE_TOTP:        totp,
	}
)

const (
	offset = 5 * time.Second  // offset into future
	period = 30 * time.Second // default value period
)

var now = func() time.Time { return time.Now().Add(offset) }

func hotp(op *Payload_OtpParameters) uint64 {
	op.Counter++ // pre-increment rfc4226 section 7.2.
	return op.Counter
}

func totp(op *Payload_OtpParameters) uint64 {
	return uint64(now().Unix()) / uint64(period.Seconds())
}

// Seconds of current validity frame
func (op *Payload_OtpParameters) Seconds() float64 {
	return now().Sub(now().Truncate(period)).Seconds()
}

// Evaluate OTP parameters
func (op *Payload_OtpParameters) Evaluate() int {
	h := hmac.New(hashFunc[op.Algorithm], op.Secret)
	binary.Write(h, binary.BigEndian, countFunc[op.Type](op))
	hashed := h.Sum(nil)
	offset := hashed[h.Size()-1] & 15
	result := binary.BigEndian.Uint32(hashed[offset:]) & (1<<31 - 1)
	return int(result) % int(math.Pow10(digitCount[op.Digits]))
}

// EvaluateString returns OTP as formatted string
func (op *Payload_OtpParameters) EvaluateString() string {
	return fmt.Sprintf("%0*d", digitCount[op.Digits], op.Evaluate())
}
