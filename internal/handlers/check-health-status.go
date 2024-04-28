package handlers

import (
	"net/http"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy"))
}
