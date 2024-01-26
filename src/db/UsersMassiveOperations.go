package db

import (
	"github.com/lib/pq"
)

type AnswerRow struct {
	IsError      bool
	ErrorMessage string
	Login        string
	Password     string
}

/*
Создаёт пользователей с user_type=2
*/
func CreateMultipleInternalUsers(rows []AnswerRow) (err error, errorRowsNum int) {
	stmt, err := db.Prepare("INSERT INTO T_USER(login, password_sha256, user_type) VALUES($1, sha256($2), 2)")

	if err != nil {
		return
	}

	for i, row := range rows {
		if !row.IsError {
			_, dbErr := stmt.Exec(row.Login, EncryptionSaltWord+row.Password)
			if dbErr != nil {
				rows[i].IsError = true
				//rows[i].ErrorMessage = "Пользователь с таким логином уже существует"
				pqErr := dbErr.(*pq.Error)
				if pqErr.Code == "23505" {
					rows[i].ErrorMessage = row.Login + ": пользователь с таким логином уже существует"
				} else {
					rows[i].ErrorMessage = row.Login + ": " + dbErr.Error() + "; ERROR CODE: " + string(pqErr.Code)
				}
				errorRowsNum++
			}
		}
	}

	return
}
