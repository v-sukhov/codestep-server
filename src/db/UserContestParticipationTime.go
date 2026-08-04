package db

import (
	"context"
	"fmt"
)

type SetAddedTimeRow struct {
	UserId       int32
	AddedTimeSec int32
}

type SetAddedTimeByLoginRow struct {
	IsError      bool
	ErrorMessage string
	Login        string
	AddedTimeSec int32
}

/*
SetAddedTime устанавливает добавленное время участникам соревнования.
Если пользователя ещё нет в T_USER_CONTEST_START, создаёт запись с NULL во временах старта.
При повторной установке значение added_time_sec перезаписывается.
*/
func SetAddedTime(contestId int32, rows []SetAddedTimeRow) error {
	if len(rows) == 0 {
		return nil
	}

	ctx := context.TODO()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		CREATE TEMP TABLE TMP_USER_CONTEST_ADDED_TIME
		(
			ROW_NUM INTEGER,
			USER_ID INTEGER,
			ADDED_TIME_SEC INTEGER
		)
	`); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO TMP_USER_CONTEST_ADDED_TIME (ROW_NUM, USER_ID, ADDED_TIME_SEC)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, row := range rows {
		if _, err := stmt.Exec(i, row.UserId, row.AddedTimeSec); err != nil {
			return err
		}
	}

	if _, err := tx.Exec(`
		INSERT INTO t_user_contest_start (user_id, contest_id, user_start_time, user_fact_start_time, added_time_sec)
		SELECT
			user_id,
			$1,
			NULL,
			NULL,
			added_time_sec
		FROM (
			SELECT DISTINCT ON (user_id)
				user_id,
				added_time_sec
			FROM TMP_USER_CONTEST_ADDED_TIME
			ORDER BY user_id, row_num DESC
		) src
		ON CONFLICT (user_id, contest_id) DO UPDATE
		SET added_time_sec = EXCLUDED.added_time_sec
	`, contestId); err != nil {
		return err
	}

	if _, err := tx.Exec(`DROP TABLE TMP_USER_CONTEST_ADDED_TIME`); err != nil {
		return err
	}

	return tx.Commit()
}

/*
SetAddedTimeByLogin устанавливает добавленное время по логинам пользователей.
Несуществующие логины пропускаются и возвращаются в списке unknownLogins.
*/
func SetAddedTimeByLogin(contestId int32, rows []SetAddedTimeByLoginRow) (unknownLogins []string, err error) {
	unknownLogins = make([]string, 0)

	ctx := context.TODO()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
		CREATE TEMP TABLE TMP_USER_CONTEST_ADDED_TIME_BY_LOGIN
		(
			ROW_NUM INTEGER,
			LOGIN VARCHAR(256),
			USER_ID INTEGER,
			ADDED_TIME_SEC INTEGER
		)
	`); err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO TMP_USER_CONTEST_ADDED_TIME_BY_LOGIN (ROW_NUM, LOGIN, ADDED_TIME_SEC)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rowNum := 0
	for _, row := range rows {
		if row.IsError {
			continue
		}
		fmt.Println(row.Login, row.AddedTimeSec)
		if _, err := stmt.Exec(rowNum, row.Login, row.AddedTimeSec); err != nil {
			return nil, err
		}
		rowNum++
	}

	if _, err := tx.Exec(`
		UPDATE TMP_USER_CONTEST_ADDED_TIME_BY_LOGIN tmp
		SET USER_ID = u.user_id
		FROM t_user u
		WHERE LOWER(u.login) = LOWER(tmp.LOGIN)
	`); err != nil {
		return nil, err
	}

	unknownRows, err := tx.Query(`
		SELECT DISTINCT LOGIN
		FROM TMP_USER_CONTEST_ADDED_TIME_BY_LOGIN
		WHERE USER_ID IS NULL
		ORDER BY LOGIN
	`)
	if err != nil {
		return nil, err
	}
	defer unknownRows.Close()

	for unknownRows.Next() {
		var login string
		if err := unknownRows.Scan(&login); err != nil {
			return nil, err
		}
		unknownLogins = append(unknownLogins, login)
	}

	setRows := make([]SetAddedTimeRow, 0)
	knownUsers, err := tx.Query(`
		SELECT DISTINCT ON (user_id)
			user_id,
			added_time_sec
		FROM TMP_USER_CONTEST_ADDED_TIME_BY_LOGIN
		WHERE USER_ID IS NOT NULL
		ORDER BY user_id, row_num DESC
	`)
	if err != nil {
		return nil, err
	}
	defer knownUsers.Close()

	for knownUsers.Next() {
		var row SetAddedTimeRow
		if err := knownUsers.Scan(&row.UserId, &row.AddedTimeSec); err != nil {
			return nil, err
		}
		setRows = append(setRows, row)
	}

	if _, err := tx.Exec(`DROP TABLE TMP_USER_CONTEST_ADDED_TIME_BY_LOGIN`); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if len(setRows) == 0 {
		return unknownLogins, nil
	}

	return unknownLogins, SetAddedTime(contestId, setRows)
}
