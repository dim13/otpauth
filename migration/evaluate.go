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
	"net/url"
	"time"

	"google.golang.org/protobuf/proto"
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
	counter = map[Payload_OtpType]func(*Payload_OtpParameters) int64{
		Payload_OTP_TYPE_UNSPECIFIED: otpTOTP, // default
		Payload_OTP_TYPE_HOTP:        otpHOTP,
		Payload_OTP_TYPE_TOTP:        otpTOTP,
	}
)

// now function for testing purposes
var now = time.Now

func otpHOTP(op *Payload_OtpParameters) int64 { return op.Counter }
func otpTOTP(op *Payload_OtpParameters) int64 { return now().Unix() / 30 }

// Evaluate OTP parameters
func (op *Payload_OtpParameters) Evaluate() int {
	h := hmac.New(hashes[op.Algorithm], op.Secret)
	binary.Write(h, binary.BigEndian, counter[op.Type](op))
	hash := h.Sum(nil)
	off := hash[h.Size()-1] & 15
	header := binary.BigEndian.Uint32(hash[off:]) & (1<<31 - 1)
	return int(header) % digits[op.Digits]
}

// Evaluate otpauth-migration URL
func Evaluate(u *url.URL) error {
	data, err := dataQuery(u)
	if err != nil {
		return err
	}
	var p Payload
	if err := proto.Unmarshal(data, &p); err != nil {
		return err
	}
	for _, op := range p.OtpParameters {
		fmt.Printf("%06d %s\n", op.Evaluate(), op.Name)
	}
	return nil
}
