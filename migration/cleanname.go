package migration

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func (op *Payload_OtpParameters) FileName() string {
	return strings.Map(func(r rune) rune {
		if r == filepath.Separator {
			return '_'
		}
		return r
	}, op.Name)
}

func (op *Payload_OtpParameters) UUID() uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, op.Secret)
}

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
