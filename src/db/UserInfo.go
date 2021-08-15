package db

import (
	"database/sql"
	"log"
)

type UserInfo struct {
	UserId      int
	UserLogin   string
	DisplayName string
	UserType    int
}

func AuthenticateUser(login string, password string) (UserInfo, bool) {
	userInfo := UserInfo{}
	ok := true
	err := db.QueryRow("SELECT user_id, login, display_name, user_type FROM t_user WHERE login = $1 and password_sha256 = sha256($2)",
		login, EncryptionSaltWord+password).
		Scan(&userInfo.UserId, &userInfo.UserLogin, &userInfo.DisplayName, &userInfo.UserType)
	if err != nil {
		ok = false
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
	} else {
		ok = true
	}

	return userInfo, ok
}
