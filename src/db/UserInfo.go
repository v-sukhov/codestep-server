package db

import (
	"database/sql"
	"log"
	"strings"
)

type UserBasicInfo struct {
	UserId      int32
	UserLogin   string
	DisplayName string
	UserType    int
}

type UserRights struct {
	UserId      int32
	IsAdmin     bool
	IsDeveloper bool
	IsJury      bool
}

func AuthenticateUser(login string, password string) (UserBasicInfo, bool) {
	userInfo := UserBasicInfo{}
	ok := true
	err := db.QueryRow("SELECT user_id, login, COALESCE(display_name, '') display_name, user_type FROM t_user WHERE login = $1 and password_sha256 = sha256($2)",
		login, EncryptionSaltWord+password).
		Scan(&userInfo.UserId, &userInfo.UserLogin, &userInfo.DisplayName, &userInfo.UserType)
	if err != nil {
		ok = false
		if err != sql.ErrNoRows {
			log.Print(err)
		}
	} else {
		ok = true
	}

	return userInfo, ok
}

func FindUserByEmail(email string) (UserBasicInfo, bool) {
	userInfo := UserBasicInfo{}
	ok := true
	err := db.QueryRow("SELECT user_id, login, COALESCE(display_name, '') display_name, user_type FROM t_user WHERE LOWER(email) = $1 ",
		strings.ToLower(email)).
		Scan(&userInfo.UserId, &userInfo.UserLogin, &userInfo.DisplayName, &userInfo.UserType)

	if err != nil {
		ok = false
		if err != sql.ErrNoRows {
			log.Print(err)
		}
	} else {
		ok = true
	}

	return userInfo, ok
}

func FindUserByLogin(login string) (UserBasicInfo, bool) {
	userInfo := UserBasicInfo{}
	ok := true
	err := db.QueryRow("SELECT user_id, login, COALESCE(display_name, '') display_name, user_type FROM t_user WHERE LOWER(login) = $1 ",
		strings.ToLower(login)).
		Scan(&userInfo.UserId, &userInfo.UserLogin, &userInfo.DisplayName, &userInfo.UserType)

	if err != nil {
		ok = false
		if err != sql.ErrNoRows {
			log.Print(err)
		}
	} else {
		ok = true
	}

	return userInfo, ok
}

func CreateUser(username string, email string, password string) (UserRights, error) {

	var userId int32

	err := db.QueryRow(`INSERT INTO t_user(login, email, password_sha256, user_type) 
						VALUES( $1, $2, sha256($3), $4 ) 
						RETURNING user_id`, username, email, EncryptionSaltWord+password, 1).Scan(&userId)
	if err == nil {
		_, err = db.Exec(`INSERT INTO t_user_rights(user_id, is_admin, is_developer, is_jury) 
		VALUES( $1, false, false, false)`, userId)
	}

	return UserRights{
			UserId:      userId,
			IsAdmin:     false,
			IsDeveloper: false,
			IsJury:      false,
		},
		err
}

func GetUserRights(userId int32) (userRights UserRights, err error) {
	err = db.QueryRow(`SELECT user_id, is_admin, is_developer, is_jury FROM t_user_rights WHERE user_id = $1`, userId).
		Scan(&userRights.UserId, &userRights.IsAdmin, &userRights.IsDeveloper, &userRights.IsJury)

	if err == sql.ErrNoRows {
		err = nil
		userRights.UserId = userId
		userRights.IsAdmin = false
		userRights.IsDeveloper = false
		userRights.IsJury = false
	}

	return
}
