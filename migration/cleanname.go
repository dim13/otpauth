package migration

import (
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// FileName returns sanitized filename without path delimiters
func (op *Payload_OtpParameters) FileName() string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return '_'
	}, op.Name+"_"+op.Issuer)
}

// UUID of OTP parameter
func (op *Payload_OtpParameters) UUID() uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, op.Secret)
}
