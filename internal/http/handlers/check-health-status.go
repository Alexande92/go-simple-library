package handlers

import (
	"net/http"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Body", "Healthy")

	w.Write([]byte("Healthy"))
}
