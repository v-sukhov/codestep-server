package security

import (
	"encoding/json"
	"log"
	"net/http"
)

type LogoutRequest struct {
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var response LogoutResponse

	token := r.Header.Get("Authorization")[7:]

	if valid, _, err := VerifyAndDecodeJWT(token); valid && err == nil {
		//deleteToken(token)
		response = LogoutResponse{
			Success: true,
			Message: "OK",
		}
	} else if !valid {
		response = LogoutResponse{
			Success: false,
			Message: "Invalid token",
		}
	} else {
		response = LogoutResponse{
			Success: false,
			Message: "Internal server error: " + err.Error(),
		}
	}

	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}

}
