package migration

import (
	"path/filepath"
	"strings"
)

func (op *Payload_OtpParameters) CleanName() string {
	return strings.Map(func(r rune) rune {
		if r == filepath.Separator {
			return '_'
		}
		return r
	}, op.Name)
}
