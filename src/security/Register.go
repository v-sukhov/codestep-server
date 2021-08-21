package security

import "net/http"

type RegisterRequest struct {
	EMail string `json:"email"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Register(w http.ResponseWriter, r *http.Request) {

}
