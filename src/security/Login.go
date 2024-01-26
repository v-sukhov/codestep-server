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
	IsAdmin     bool   `json:"isAdmin"`
	IsDeveloper bool   `json:"isDeveloper"`
	IsJury      bool   `json:"isJury"`
}

type LoginResponseData struct {
	Token    string   `json:"token"`
	UserInfo UserInfo `json:"userInfo"`
}

type LoginResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    LoginResponseData `json:"data"`
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
			if userBasicInfo, ok := db.AuthenticateUser(request.Login, request.Password); ok {
				if token, err := generateJWT(userBasicInfo.UserId); err != nil {
					response = LoginResponse{
						Success: false,
						Message: "Internal server error: " + err.Error(),
					}
				} else {
					if userRights, err := db.GetUserRights(userBasicInfo.UserId); err != nil {
						response = LoginResponse{
							Success: false,
							Message: "Internal server error: " + err.Error(),
						}
					} else {
						userInfo := UserInfo{
							UserId:      userBasicInfo.UserId,
							UserLogin:   userBasicInfo.UserLogin,
							DisplayName: userBasicInfo.DisplayName,
							UserType:    userBasicInfo.UserType,
							IsAdmin:     userRights.IsAdmin,
							IsDeveloper: userRights.IsDeveloper,
							IsJury:      userRights.IsJury,
						}
						data := LoginResponseData{
							UserInfo: userInfo,
							Token:    "Bearer " + token,
						}
						response = LoginResponse{
							Success: true,
							Message: "OK",
							Data:    data,
						}
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
		log.Print(err)
	} else {
		w.Write(byteArr)
	}

}
