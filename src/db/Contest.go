package db

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
)

type Contest struct {
	ContestId       int32
	ContestName     string
	ContestDesc     string
	ContestLogoHref string
}

type ContestWithRights struct {
	ContestId         int32  `json:"contestId"`
	ContestName       string `json:"contestName"`
	ContestDesc       string `json:"contestDesc"`
	ContestLogoHref   string `json:"contestLogoHref"`
	UserRightsBitmask int32  `json:"userRightsBitmask"`
}

type SupertaskInContestInfo struct {
	SupertaskId   int32  `json:"supertaskId"`
	VersionNumber int32  `json:"versionNumber"`
	SupertaskName string `json:"supertaskName"`
	SupertaskDesc string `json:"supertaskDesc"`
	OrderNumber   int16  `json:"orderNumber"`
}

/*
Сохраняет контест
Если contestId = 0, значит создаёт новый
Права пользователя не проверяются
*/
func SaveContest(contest *Contest, userId int32) error {
	if contest.ContestId == 0 {
		ctx := context.TODO()
		tx, err := db.BeginTx(ctx, nil)

		if err != nil {
			return err
		}
		defer tx.Rollback()

		err = tx.QueryRow(`INSERT INTO T_CONTEST(CONTEST_NAME, CONTEST_DESC, CONTEST_LOGO_HREF)
							VALUES($1,$2,$3) RETURNING CONTEST_ID`,
			contest.ContestName, contest.ContestDesc, contest.ContestLogoHref).
			Scan(&contest.ContestId)

		if err != nil {
			return err
		}

		_, err = tx.Exec(`INSERT INTO T_CONTEST_USER_RIGHT(USER_ID, CONTEST_ID, CONTEST_RIGHT_TYPE_ID)
						VALUES($1, $2, 1)`, userId, contest.ContestId)

		tx.Commit()
	} else {
		res, err := db.Exec(`UPDATE T_CONTEST
						SET CONTEST_NAME = $1,
							CONTEST_DESC = $2,
							CONTEST_LOGO_HREF = $3
						WHERE
							CONTEST_ID = $4
						`, contest.ContestName, contest.ContestDesc, contest.ContestLogoHref, contest.ContestId)
		if err != nil {
			return err
		}

		rowsAffected, err := res.RowsAffected()

		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			err = errors.New("No contest found with givven id: " + strconv.Itoa(int(contest.ContestId)))
			return err
		}
	}

	return nil
}

/*
Получить информацию о контесте, права пользователя не проверяет
*/
func GetContest(contestId int32) (contest Contest, err error) {
	err = db.QueryRow(`
			SELECT
				CONTEST_ID,
				CONTEST_NAME,
				CONTEST_DESC,
				CONTEST_LOGO_HREF
			FROM
				T_CONTEST
			WHERE
				CONTEST_ID = $1
		`, contestId).Scan(
		&contest.ContestId,
		&contest.ContestName,
		&contest.ContestDesc,
		&contest.ContestLogoHref,
	)
	return
}

func GetUserContestList(userId int32) (contests []ContestWithRights, err error) {
	rows, err := db.Query(`
		select
			ur.contest_id,
			c.contest_name,
			c.contest_desc,
			c.contest_logo_href,
			ur.user_contest_rights_bitmask
		from
			t_contest c join
			(
			select
				user_id,
				contest_id,
				sum(1 << (contest_right_type_id - 1)) as user_contest_rights_bitmask
			from
				t_contest_user_right
			where
				user_id = $1
			group by
				user_id,
				contest_id
			) ur on c.contest_id = ur.contest_id
		`, userId)

	if err != nil {
		return
	}

	contests = make([]ContestWithRights, 0)

	for rows.Next() {
		var contest ContestWithRights

		if err = rows.Scan(
			&contest.ContestId,
			&contest.ContestName,
			&contest.ContestDesc,
			&contest.ContestLogoHref,
			&contest.UserRightsBitmask,
		); err != nil {
			return
		}

		contests = append(contests, contest)
	}

	return
}

/*
Запрос прав пользователя на контест - возвращает в виде битовой маски прав
*/
func GetUserContestRights(userId int32, contestId int32) (rightsBitmask int32, err error) {
	err = db.QueryRow(`
			select
				sum(1 << (contest_right_type_id - 1)) as user_contest_rights_bitmask
			from
				t_contest_user_right
			where
				user_id = $1
			group by
				user_id,
				contest_id	
		`, userId).Scan(&rightsBitmask)

	if err == sql.ErrNoRows {
		rightsBitmask = 0
		err = nil
	}

	return
}

/*
	Если данной суперзадачи ещё нет в контесте - добавляет её
	Права пользователя не проверяются
	supertaskVersionNumber должно быть > 0
*/

func AddSupertaskToContest(contestId int32, supertaskId int32, supertaskVersionNumber int32) (orderNumber int16, err error) {

	if supertaskVersionNumber <= 0 {
		err = errors.New("Supertask version number should be larger than zero")
		return
	}

	ctx := context.TODO()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return
	}
	defer tx.Rollback()

	var count int32
	_, err = tx.Exec(`LOCK TABLE T_CONTEST_SUPERTASK IN EXCLUSIVE MODE`)

	if err != nil {
		return
	}

	err = tx.QueryRow(`
		SELECT
			count(*)
		FROM
			T_CONTEST_SUPERTASK
		WHERE
			CONTEST_ID = $1 and SUPERTASK_ID = $2
	`, contestId, supertaskId).Scan(&count)

	if err != nil {
		return
	}

	if count == 0 {
		var maxOrderNumber int16
		err = tx.QueryRow(`
				SELECT
					coalesce (max(ORDER_NUMBER), 0)
				FROM
					T_CONTEST_SUPERTASK
				WHERE
					CONTEST_ID = $1
			`, contestId).Scan(&maxOrderNumber)

		if err != nil {
			return
		}

		orderNumber = maxOrderNumber + 1

		_, err = tx.Exec(`
			INSERT INTO T_CONTEST_SUPERTASK(CONTEST_ID, SUPERTASK_ID, SUPERTASK_VERSION_NUMBER, ORDER_NUMBER)
			VALUES($1, $2, $3, $4)
		`, contestId, supertaskId, supertaskVersionNumber, orderNumber)

		if err != nil {
			return
		}

		if err = tx.Commit(); err != nil {
			return
		}
	}

	return
}

/*
	Удаляет суперзадачу из контеста и перенумеровывает остальные суперзадачи
	Права пользователя не проверяются
*/

func RemoveSupertaskFromContest(contestId int32, supertaskId int32) error {
	ctx := context.TODO()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`LOCK TABLE T_CONTEST_SUPERTASK IN EXCLUSIVE MODE`)

	var orderNumber int32
	err = tx.QueryRow(`
		SELECT
			ORDER_NUMBER
		FROM
			T_CONTEST_SUPERTASK
		WHERE
			CONTEST_ID = $1 and SUPERTASK_ID = $2
	`, contestId, supertaskId).Scan(&orderNumber)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if orderNumber > 0 {

		_, err = tx.Exec(`
			DELETE FROM T_CONTEST_SUPERTASK
			WHERE
				CONTEST_ID = $1 and SUPERTASK_ID = $2
		`, contestId, supertaskId)

		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			UPDATE T_CONTEST_SUPERTASK
				SET ORDER_NUMBER = ORDER_NUMBER - 1
			WHERE
				CONTEST_ID = $1 and ORDER_NUMBER > $2
		`, contestId, orderNumber)

		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

/*
Возвращает список суперзадач контеста
*/
func GetContestSupertasksList(contestId int32) (supertaskList []SupertaskInContestInfo, err error) {
	rows, err := db.Query(`
		SELECT
			cs.supertask_id,
			cs.supertask_version_number,
			sv.supertask_name,
			sv.supertask_desc,
			cs.order_number
		FROM
			t_contest_supertask cs join
			t_supertask_version sv on cs.supertask_id = sv.supertask_id and cs.supertask_version_number = sv.version_number
		WHERE 
			cs.contest_id = $1
		ORDER BY
			cs.order_number
		`, contestId)

	if err != nil {
		return
	}

	supertaskList = make([]SupertaskInContestInfo, 0)

	for rows.Next() {
		var supertaskInfo SupertaskInContestInfo
		err = rows.Scan(
			&supertaskInfo.SupertaskId,
			&supertaskInfo.VersionNumber,
			&supertaskInfo.SupertaskName,
			&supertaskInfo.SupertaskDesc,
			&supertaskInfo.OrderNumber,
		)
		if err != nil {
			return
		}
		supertaskList = append(supertaskList, supertaskInfo)
	}

	return

}
