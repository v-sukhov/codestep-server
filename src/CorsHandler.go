package main

import (
	"net/http"
)

func CorsHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization , Content-Type")
		if r.Method != http.MethodOptions {
			h.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(fn)
}
