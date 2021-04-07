package migration

import (
	"io/ioutil"

	"github.com/skip2/go-qrcode"
)

// QR image bytes as PNG
func (op *Payload_OtpParameters) QR() ([]byte, error) {
	return qrcode.Encode(op.URL().String(), qrcode.Medium, -3)
}

// WriteFile writes QR code as PNG to specified file
func (op *Payload_OtpParameters) WriteFile(fname string) error {
	pic, err := op.QR()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fname, pic, 0600); err != nil {
		return err
	}
	return nil
}
