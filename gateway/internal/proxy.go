package internal

import (
	"net/http"

	"github.com/rs/zerolog"
)

func NewProxy(
	config *Config,
	logger *zerolog.Logger,
) http.Handler {
	mux := http.NewServeMux()
	addRoute(mux)

	var handler http.Handler
	handler = mux
	return handler
}

func addRoute(mux *http.ServeMux) {

	mux.Handle("/", GreetHandler())
}
