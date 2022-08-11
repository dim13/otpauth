package migration

import "net/http"

func (op *Payload_OtpParameters) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pic, err := QR(op.URL())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(pic)
}
