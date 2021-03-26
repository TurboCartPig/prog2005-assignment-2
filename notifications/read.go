package notifications

import "net/http"

func NewReadHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func NewReadAllHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}
