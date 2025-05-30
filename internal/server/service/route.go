package service

import "net/http"

func NewServeMux(update *UpdateHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/update/", update)

	return mux
}
