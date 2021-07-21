package migration

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

// FileName returns sanitized filename without path delimiters
func (op *Payload_OtpParameters) FileName() string {
	return strings.Map(func(r rune) rune {
		switch r {
		case os.PathSeparator, os.PathListSeparator:
			return '_'
		default:
			return r
		}
	}, op.Name+"_"+op.Issuer)
}

// UUID of OTP parameter
func (op *Payload_OtpParameters) UUID() uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, op.Secret)
}

// Title of OTP parameter
func (op *Payload_OtpParameters) Title() string {
	// strip issuer from Name
	name := op.Name
	if i := strings.Index(name, ":"); i > 0 {
		name = name[i+1:]
	}
	if op.Issuer == "" {
		return name
	}
	return op.Issuer + " (" + name + ")"
}
