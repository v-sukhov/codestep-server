package security

import (
	"encoding/json"
	"io"
	"log"

	"net/http"

	"codestep/db"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserInfo struct {
	UserId      int32  `json:"userId"`
	UserLogin   string `json:"userLogin"`
	DisplayName string `json:"userDisplayName"`
	UserType    int    `json:"userType"`
	IsAdmin     bool   `json:"roleAdmin"`
	IsDeveloper bool   `json:"roleDeveloper"`
	IsJury      bool   `json:"roleJury"`
}

type LoginResponseData struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"userInfo"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	var request LoginRequest
	var response LoginResponse

	if body, err := io.ReadAll(r.Body); err != nil {
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

				/*token := addUser(UserInfoCache{
					Id:       userInfo.UserId,
					UserType: userInfo.UserType,
				})*/

				if token, err := generateJWT(userInfo.UserId); err != nil {
					response = LoginResponse{
						Success: false,
						Message: "Internal server error: " + err.Error(),
					}
				} else {
					response = LoginResponse{
						Success: true,
						Message: "OK",
						Token:   "Bearer " + token,
					}
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
		log.Fatal(err)
	} else {
		w.Write(byteArr)
	}

}
