package migration

import (
	"io/ioutil"

	"github.com/skip2/go-qrcode"
)

func (op *Payload_OtpParameters) QR() ([]byte, error) {
	return qrcode.Encode(op.URL().String(), qrcode.Medium, 256)
}

func (op *Payload_OtpParameters) WriteFile() error {
	pic, err := op.QR()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(op.FileName()+".png", pic, 0600); err != nil {
		return err
	}
	return nil
}
