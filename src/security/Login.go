package security

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"

	"net/http"

	"codestep/db"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

func generateToken() string {
	len := 256
	b := make([]byte, len)

	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func Login(w http.ResponseWriter, r *http.Request) {

	var request LoginRequest
	var response LoginResponse

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		response = LoginResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		err := json.Unmarshal(body, &request)
		if err != nil {
			response = LoginResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {
			if userInfo, ok := db.AuthenticateUser(request.Login, request.Password); ok {
				var token string
				for token = ""; token == ""; {
					token = generateToken()
					if _, ok := getUserByToken(token); !ok {
						addUserByToken(token, UserInfoCache{
							Id:       userInfo.UserId,
							UserType: userInfo.UserType,
						})
					} else {
						token = ""
					}
				}

				response = LoginResponse{
					Success: true,
					Message: "OK",
					Token:   "Bearer " + token,
				}
			} else {
				response = LoginResponse{
					Success: false,
					Message: "Incorrect login/password",
				}
			}
		}
	}

	if byteArr, err := json.Marshal(response); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}

}
