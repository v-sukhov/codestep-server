package security

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	JwtSecret                      string
	JwtRegisterSecret              string
	JwtTokenLifetimeMinute         int
	JwtRegisterTokenLifetimeMinute int
)

type JWTClaims struct {
	UserId int32 `json:"userId"`
	jwt.StandardClaims
}

func generateJWT(userId int32) (string, error) {
	var signingKey = []byte(JwtSecret)

	claims := JWTClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(JwtTokenLifetimeMinute)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключем
	tokenString, err := token.SignedString(signingKey)

	return tokenString, err
}

func VerifyAndDecodeJWT(tokenStr string) (bool, int32, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecret), nil
	})

	var userId = int32(0)
	var valid = false
	if token != nil && err == nil {
		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			userId = claims.UserId
			valid = true
		}
	}

	return valid, userId, err
}
