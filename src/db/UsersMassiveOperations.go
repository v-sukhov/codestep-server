package db

import (
	"context"

	"github.com/lib/pq"
)

type CreateUserAnswerRow struct {
	IsError      bool
	ErrorMessage string
	Login        string
	Password     string
}

type ManageUserContestRightsRow struct {
	IsError      bool
	ErrorMessage string
	Login        string
	ContestId    int32
	IsRevoke     bool
	Rights       []int16
}

type ManageUserContestRightsOutputRow struct {
}

/*
Создаёт пользователей с user_type=2
*/
func CreateMultipleInternalUsers(rows []CreateUserAnswerRow) (err error, errorRowsNum int) {
	stmt, err := db.Prepare("INSERT INTO T_USER(login, password_sha256, user_type) VALUES($1, sha256($2), 2)")

	if err != nil {
		return
	}

	for i, row := range rows {
		if !row.IsError {
			_, dbErr := stmt.Exec(row.Login, EncryptionSaltWord+row.Password)
			if dbErr != nil {
				rows[i].IsError = true
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

/*
	Задаёт права пользователей на контест
*/

func ManageUserContestRights(rows []ManageUserContestRightsRow) error {

	ctx := context.TODO()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		CREATE TEMP TABLE TMP_USER_CONTEST_RIGHT
		(
			ROW_NUM					INTEGER,
			LOGIN					VARCHAR(256),
			USER_ID					INTEGER,
			CONTEST_ID				INTEGER,
			CONTEST_RIGHT_TYPE_ID 	SMALLINT
		)
	`); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO TMP_USER_CONTEST_RIGHT
		(
			ROW_NUM,
			LOGIN,
			CONTEST_ID,
			CONTEST_RIGHT_TYPE_ID
		)
		VALUES($1, $2, $3, $4)
	`)

	if err != nil {
		return err
	}

	for i, row := range rows {
		if !row.IsError {
			if row.IsRevoke {
				if _, err := stmt.Exec(i, row.Login, row.ContestId, 0); err != nil {
					return err
				}
			} else {
				for _, rightTypeId := range row.Rights {
					if _, err := stmt.Exec(i, row.Login, row.ContestId, rightTypeId); err != nil {
						return err
					}
				}
			}
		}
	}

	if _, err := tx.Exec(`
		UPDATE TMP_USER_CONTEST_RIGHT T
			SET USER_ID = U.USER_ID
		FROM
			T_USER U
		WHERE
			T.LOGIN = U.LOGIN 		
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		DELETE FROM T_CONTEST_USER_RIGHT CR
		USING TMP_USER_CONTEST_RIGHT T
		WHERE CR.USER_ID = T.USER_ID AND CR.CONTEST_ID = T.CONTEST_ID AND CR.CONTEST_RIGHT_TYPE_ID > 1	
	`); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		INSERT INTO T_CONTEST_USER_RIGHT
		(
			USER_ID,
			CONTEST_ID,
			CONTEST_RIGHT_TYPE_ID 
		)
		SELECT
			USER_ID,
			CONTEST_ID,
			CONTEST_RIGHT_TYPE_ID
		FROM
			TMP_USER_CONTEST_RIGHT
		WHERE
			USER_ID IS NOT NULL AND CONTEST_RIGHT_TYPE_ID > 0
		GROUP BY
			USER_ID,
			CONTEST_ID,
			CONTEST_RIGHT_TYPE_ID	
	`); err != nil {
		return err
	}

	/*
		Выявляем неправильные логины и отмечаем в массиве входных данных
	*/

	unknownLogins, err := tx.Query(`
		SELECT
			ROW_NUM
		FROM
			TMP_USER_CONTEST_RIGHT
		WHERE
			USER_ID IS NULL
	`)
	if err != nil {
		return err
	}

	for unknownLogins.Next() {
		var row_num int
		unknownLogins.Scan(&row_num)
		rows[row_num].IsError = true
		rows[row_num].ErrorMessage = "Unknown login: " + rows[row_num].Login
	}

	if _, err := tx.Exec(`DROP TABLE TMP_USER_CONTEST_RIGHT`); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
