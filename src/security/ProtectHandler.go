package security

import (
	"fmt"
	"net/http"
)

func ProtectHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Protected!")
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
