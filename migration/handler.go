package migration

import "net/http"

func (op *Payload_OtpParameters) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pic, err := op.QR()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(pic)
}
