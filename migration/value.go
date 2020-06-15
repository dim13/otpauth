package migration

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"hash"
	"time"
)

var (
	hashes = map[Payload_Algorithm]func() hash.Hash{
		Payload_ALGORITHM_UNSPECIFIED: sha1.New, // default
		Payload_ALGORITHM_SHA1:        sha1.New,
		Payload_ALGORITHM_SHA256:      sha256.New,
		Payload_ALGORITHM_SHA512:      sha512.New,
		Payload_ALGORITHM_MD5:         md5.New,
	}
	digits = map[Payload_DigitCount]int{
		Payload_DIGIT_COUNT_UNSPECIFIED: 1e6, // default
		Payload_DIGIT_COUNT_SIX:         1e6,
		Payload_DIGIT_COUNT_EIGHT:       1e8,
	}
)

func (op *Payload_OtpParameters) value(c int64) int {
	h := hmac.New(hashes[op.Algorithm], op.Secret)
	binary.Write(h, binary.BigEndian, c)
	hash := h.Sum(nil)
	off := hash[h.Size()-1] & 15
	header := binary.BigEndian.Uint32(hash[off:]) & (1<<31 - 1)
	return int(header) % digits[op.Digits]
}

// Value of OTP
func (op *Payload_OtpParameters) Value() int {
	switch op.Type {
	case Payload_OTP_TYPE_HOTP:
		return op.value(op.Counter) // TODO increment counter
	case Payload_OTP_TYPE_TOTP:
		return op.value(time.Now().Unix() / 30) // default period 30s
	}
	return 0
}
