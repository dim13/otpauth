package migration

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"

	"google.golang.org/protobuf/proto"
)

// ErrUnkown scheme or host
var ErrUnkown = errors.New("unknown")

func unmarshal(u *url.URL) ([]byte, error) {
	if u.Scheme != "otpauth-migration" {
		return nil, fmt.Errorf("scheme %s: %w", u.Scheme, ErrUnkown)
	}
	if u.Host != "offline" {
		return nil, fmt.Errorf("host %s: %w", u.Host, ErrUnkown)
	}
	data := u.Query().Get("data")
	return base64.StdEncoding.DecodeString(data)
}

// Unmarshal otpauth-migration URL
func Unmarshal(u *url.URL) (*Payload, error) {
	data, err := unmarshal(u)
	if err != nil {
		return nil, err
	}
	var p Payload
	if err := proto.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
