package main

import (
	"log"
	"net/http"
)

var CorsLogHttp bool

func CorsHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if CorsLogHttp {
			log.Println("Http request: ")
			log.Println(r)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		if r.Method != http.MethodOptions {
			h.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(fn)
}
