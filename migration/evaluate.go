package migration

import (
	"crypto/hmac"
	"encoding/binary"
	"fmt"
	"math"
	"time"
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
	h := hmac.New(op.Algorithm.Hash(), op.Secret)
	binary.Write(h, binary.BigEndian, op.Type.Count(op))
	hashed := h.Sum(nil)
	offset := hashed[h.Size()-1] & 15
	result := binary.BigEndian.Uint32(hashed[offset:]) & (1<<31 - 1)
	return int(result) % int(math.Pow10(op.Digits.Count()))
}

// EvaluateString returns OTP as formatted string
func (op *Payload_OtpParameters) EvaluateString() string {
	return fmt.Sprintf("%0*d", op.Digits.Count(), op.Evaluate())
}
