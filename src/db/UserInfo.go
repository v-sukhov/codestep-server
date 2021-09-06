package db

import (
	"database/sql"
	"log"
	"strings"
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
	err := db.QueryRow("SELECT user_id, login, COALESCE(display_name, '') display_name, user_type FROM t_user WHERE login = $1 and password_sha256 = sha256($2)",
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

func FindUserByEmail(email string) (UserInfo, bool) {
	userInfo := UserInfo{}
	ok := true
	err := db.QueryRow("SELECT user_id, login, COALESCE(display_name, '') display_name, user_type FROM t_user WHERE LOWER(email) = $1 ",
		strings.ToLower(email)).
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

func FindUserByLogin(login string) (UserInfo, bool) {
	userInfo := UserInfo{}
	ok := true
	err := db.QueryRow("SELECT user_id, login, COALESCE(display_name, '') display_name, user_type FROM t_user WHERE LOWER(login) = $1 ",
		strings.ToLower(login)).
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

func CreateUser(username string, email string, password string) (int64, error) {

	stmt, err := db.Prepare("INSERT INTO t_user(login, email, password_sha256, user_type) VALUES( $1, $2, $3, $4 )")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(username, email, EncryptionSaltWord+password, 1)
	if err != nil {
		log.Fatal(err)
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return rowCnt, err
}
