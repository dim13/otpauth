package migration

import (
	"net/url"
	"os"

	"github.com/skip2/go-qrcode"
)

// QR image bytes as PNG
func QR(u *url.URL) ([]byte, error) {
	return qrcode.Encode(u.String(), qrcode.Medium, -3)
}

// PNG writes QR code as PNG to specified file
func PNG(filename string, u *url.URL) error {
	pic, err := QR(u)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, pic, 0600)
}
