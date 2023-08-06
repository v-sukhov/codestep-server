package security

import (
	"codestep/db"
	"codestep/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type RegisterRequest struct {
	EMail    string `json:"email"`
	Username string `json:"username"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ValidateRequest struct {
	Code string `json:"code"`
}

type ValidateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResumeRegisterRequest struct {
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordconfirmation"`
}

type ResumeRegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
	Login   string `json:"login"`
}

type RegisterMapClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func getRegisterJWT(email string, username string) string {

	var signingKey = []byte(JwtRegisterSecret)

	claims := RegisterMapClaims{
		email,
		username,
		//time.Now().Add(time.Hour * 24).Unix(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(JwtRegisterTokenLifetimeMinute)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключем
	tokenString, _ := token.SignedString(signingKey)

	return tokenString
}

func verifyJWT(tokenStr string) (bool, string, string, error) {

	token, err := jwt.ParseWithClaims(tokenStr, &RegisterMapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtRegisterSecret), nil
	})

	var email = ""
	var username = ""
	var valid = false
	if claims, ok := token.Claims.(*RegisterMapClaims); ok && token.Valid {
		email = claims.Email
		username = claims.Username
		valid = true
	}

	return valid, email, username, err
}

func Register(w http.ResponseWriter, r *http.Request) {

	var request RegisterRequest
	var response RegisterResponse
	if r.Method == "POST" {

		if body, err := ioutil.ReadAll(r.Body); err != nil {
			response = RegisterResponse{
				Success: false,
				Message: "Body reading failed",
			}
		} else {
			err := json.Unmarshal(body, &request)
			if err != nil {
				response = RegisterResponse{
					Success: false,
					Message: "JSON decoding failed",
				}
			} else if !utils.ValidEmail(request.EMail) {
				response = RegisterResponse{
					Success: false,
					Message: "Invalid e-mail address",
				}
			} else if _, ok := db.FindUserByEmail(request.EMail); ok {
				response = RegisterResponse{
					Success: false,
					Message: "User with this email-address already exists",
				}
			} else if _, ok := db.FindUserByLogin(request.Username); ok {
				response = RegisterResponse{
					Success: false,
					Message: "User with this login already exists",
				}
			} else {
				subject := "Subject: Подтверждение регистрации на Code-Step!\n"
				var tokenStr = getRegisterJWT(request.EMail, request.Username)
				var protocol, host, port = utils.GetClientUrl()
				var confirmUrl = fmt.Sprintf("%s://%s:%s/confirmation/%s", protocol, host, port, tokenStr)
				var msgStr = fmt.Sprintf("Ваша ссылка для подтверждения e-mail адреса\n\n%s", confirmUrl)

				msg := []byte(subject + msgStr)

				var to = []string{request.EMail}

				utils.SendEmail(to, msg)
				response = RegisterResponse{
					Success: true,
					Message: "Confirmation mail sent",
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
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func ValidateRegisterCode(w http.ResponseWriter, r *http.Request) {

	var request ValidateRequest
	var response ValidateResponse
	if r.Method == "POST" {

		if body, err := ioutil.ReadAll(r.Body); err != nil {
			response = ValidateResponse{
				Success: false,
				Message: "Body reading failed",
			}
		} else {
			err := json.Unmarshal(body, &request)
			if err != nil {
				response = ValidateResponse{
					Success: false,
					Message: "JSON decoding failed",
				}
			} else {
				var valid, _, _, err = verifyJWT(request.Code)

				if valid {
					response = ValidateResponse{
						Success: true,
						Message: "The verification code is correct ",
					}
				} else {
					errMes := fmt.Sprintf("%s", err)
					if err == nil {
						errMes = "The verification code is not correct or has expired "
					}
					response = ValidateResponse{
						Success: false,
						Message: errMes,
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
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func ResumeRegister(w http.ResponseWriter, r *http.Request) {

	var request ResumeRegisterRequest
	var response ResumeRegisterResponse
	if r.Method == "POST" {

		if body, err := ioutil.ReadAll(r.Body); err != nil {
			response = ResumeRegisterResponse{
				Success: false,
				Message: "Body reading failed",
			}
		} else {
			err := json.Unmarshal(body, &request)
			if err != nil {
				response = ResumeRegisterResponse{
					Success: false,
					Message: "JSON decoding failed",
				}
			} else {
				valid, email, username, verifyErr := verifyJWT(request.Token)

				if !valid {
					errMes := fmt.Sprintf("%s", verifyErr)
					if err == nil {
						errMes = "The verification code is not correct or has expired "
					}
					response = ResumeRegisterResponse{
						Success: false,
						Message: errMes,
					}
				} else if !utils.ValidEmail(email) {
					response = ResumeRegisterResponse{
						Success: false,
						Message: "Invalid e-mail address",
					}
				} else if _, ok := db.FindUserByEmail(email); ok {
					response = ResumeRegisterResponse{
						Success: false,
						Message: "User with this email-address already exists",
					}
				} else if _, ok := db.FindUserByLogin(username); ok {
					response = ResumeRegisterResponse{
						Success: false,
						Message: "User with this login already exists",
					}
				} else if request.Password != request.PasswordConfirmation {
					response = ResumeRegisterResponse{
						Success: false,
						Message: "Password mismatch ",
					}
				} else {
					userId, err := db.CreateUser(username, email, request.Password)

					if err != nil {
						response = ResumeRegisterResponse{
							Success: false,
							Message: "Failed to create user ",
						}
					} else {
						/*token := addUser(UserInfoCache{
							Id:       userId,
							UserType: REGULAR,
						})*/
						if token, err := generateJWT(userId); err != nil {
							response = ResumeRegisterResponse{
								Success: false,
								Message: "Inernal server error: " + err.Error(),
							}
						} else {
							response = ResumeRegisterResponse{
								Success: true,
								Message: "Registration completed successfully ",
								Token:   "Bearer " + token,
								Login:   username,
							}
						}
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
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
