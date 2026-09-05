package internal

import "net/http"

func GreetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HelloWorld"))
	})
}
