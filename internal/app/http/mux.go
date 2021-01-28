package http

import "github.com/go-chi/chi"

var mux *chi.Mux

// Mux for http server
func Mux() (*chi.Mux, func(), error) {
	if mux != nil {
		return mux, func() {}, nil
	}

	mux = chi.NewRouter()

	return mux, func() {}, nil
}
