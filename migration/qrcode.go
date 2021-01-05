package migration

import (
	"io/ioutil"

	"github.com/skip2/go-qrcode"
)

func (op *Payload_OtpParameters) QR() ([]byte, error) {
	return qrcode.Encode(op.URL().String(), qrcode.Medium, 256)
}

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
