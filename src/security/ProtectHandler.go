package security

import (
	"context"
	"net/http"
)

type ContextKey string

const ContextUserIdKey ContextKey = "userId"

func ProtectHandler(h http.Handler) http.Handler {
	// TODO Add authorization check
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header["Authorization"][0][7:]

		if valid, userId, err := VerifyAndDecodeJWT(token); valid && err == nil {
			ctx := context.WithValue(r.Context(), ContextUserIdKey, userId)
			h.ServeHTTP(w, r.WithContext(ctx))
		} else if !valid {
			http.Error(w, "Unauthorized request", http.StatusUnauthorized)
		} else {
			http.Error(w, "JWT decode eror", http.StatusInternalServerError)
		}
	}

	return http.HandlerFunc(fn)
}
