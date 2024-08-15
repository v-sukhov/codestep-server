package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Contest struct {
	ContestId       int32                    `json:"contestId"`
	ContestName     string                   `json:"contestName"`
	ContestDesc     string                   `json:"contestDesc"`
	ContestLogoHref string                   `json:"contestLogoHref"`
	SupertaskList   []SupertaskInContestInfo `json:"supertaskList"`
}

type ContestWithRights struct {
	ContestId         int32  `json:"contestId"`
	ContestName       string `json:"contestName"`
	ContestDesc       string `json:"contestDesc"`
	ContestLogoHref   string `json:"contestLogoHref"`
	UserRightsBitmask int32  `json:"userRightsBitmask"`
}

type SupertaskInContestInfo struct {
	SupertaskId       int32  `json:"supertaskId"`
	VersionNumber     int32  `json:"versionNumber"`
	SupertaskName     string `json:"supertaskName"`
	SupertaskDesc     string `json:"supertaskDesc"`
	SupertaskLogoHref string `json:"supertaskLogoHref"`
	OrderNumber       int16  `json:"orderNumber"`
}

type SupertaskInContestWithResults struct {
	ContestId           int32                  `json:"contestId"`
	SupertaskInfo       SupertaskInContestInfo `json:"supertaskInfo"`
	TaskResults         []TaskResult           `json:"taskResults"`
	SupertaskObjectJson string                 `json:"supertaskObjectJson"`
}

type ContestUserResult struct {
	UserLogin       string    `json:"userLogin"`
	SupertaskScore  [][]int32 `json:"supertaskScore"`
	SupertaskTries  [][]int32 `json:"supertaskTries"`
	SupertaskPassed [][]bool  `json:"supertaskPassed"`
	TotalPassed     int32     `json:"totalPassed"`
	TotalScore      int32     `json:"totalScore"`
	TotalTries      int32     `json:"totalTries"`
}

type ContestResults struct {
	SupertaskNames    []string            `json:"supertaskNames"`
	TotalTasksNum     int32               `json:"totalTasksNum"`
	MaxPossibleResult ContestUserResult   `json:"maxPossibleResult"`
	UserResults       []ContestUserResult `json:"userResults"`
	Errors            []string            `json:"resultsErrors"`
}

/*
Сохраняет контест
Если contestId = 0, значит создаёт новый
Права пользователя не проверяются
Из списка задач использует только supertaskId и versionNumber
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

		if err != nil {
			return err
		}

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

	err := RewriteContestSupertaskList(contest.ContestId, contest.SupertaskList)

	return err
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
	if err != nil {
		return
	}

	supertaskList, err := GetContestSupertaskList(contestId)
	if err != nil {
		return
	}

	contest.SupertaskList = supertaskList

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
		order by
			ur.contest_id desc
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
func GetContestUserRights(userId int32, contestId int32) (rightsBitmask int32, err error) {
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
					coalesce (max(ORDER_NUMBER), -1)
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
	if err != nil {
		return err
	}

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
	Переписать список задач контеста
	Из списка задач использует только supertaskId и versionNumber
	Порядок создаёт в соответствии с порядком элементов массива, а не orderNumber
*/

func RewriteContestSupertaskList(contestId int32, supertaskList []SupertaskInContestInfo) error {

	ctx := context.TODO()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`LOCK TABLE T_CONTEST_SUPERTASK IN EXCLUSIVE MODE`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
			DELETE FROM T_CONTEST_SUPERTASK
			WHERE
				CONTEST_ID = $1
		`, contestId)
	if err != nil {
		return err
	}

	orderNumber := 0

	stmt, err := tx.Prepare("INSERT INTO T_CONTEST_SUPERTASK(CONTEST_ID, SUPERTASK_ID, SUPERTASK_VERSION_NUMBER, ORDER_NUMBER) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	for _, row := range supertaskList {
		orderNumber++
		_, dbErr := stmt.Exec(contestId, row.SupertaskId, row.VersionNumber, orderNumber)
		if dbErr != nil {
			return dbErr
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

/*
Возвращает список суперзадач контеста,
отсортированный в порядке order_number
*/
func GetContestSupertaskList(contestId int32) (supertaskList []SupertaskInContestInfo, err error) {
	rows, err := db.Query(`
		SELECT
			cs.supertask_id,
			cs.supertask_version_number,
			sv.supertask_name,
			sv.supertask_desc,
			sv.supertask_logo_href,
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

	defer rows.Close()

	supertaskList = make([]SupertaskInContestInfo, 0)

	for rows.Next() {
		var supertaskInfo SupertaskInContestInfo
		err = rows.Scan(
			&supertaskInfo.SupertaskId,
			&supertaskInfo.VersionNumber,
			&supertaskInfo.SupertaskName,
			&supertaskInfo.SupertaskDesc,
			&supertaskInfo.SupertaskLogoHref,
			&supertaskInfo.OrderNumber,
		)
		if err != nil {
			return
		}
		supertaskList = append(supertaskList, supertaskInfo)
	}

	return

}

/*
Возвращает актуальную версию суперзадачи в контесте
вместе с результатами данного пользователя
*/
func GetSupertaskInContestWithResults(contestId int32, supertaskId int32, userId int32) (supertaskInContest SupertaskInContestWithResults, err error) {

	supertaskInContest.ContestId = contestId

	err = db.QueryRow(`
					SELECT
						cs.supertask_id,
						cs.supertask_version_number,
						sv.supertask_name,
						sv.supertask_desc,
						sv.supertask_logo_href,
						cs.order_number,
						sv.supertask_object_json
					FROM
						t_contest_supertask cs join
						t_supertask_version sv on cs.supertask_id = sv.supertask_id and cs.supertask_version_number = sv.version_number
					WHERE 
						cs.contest_id = $1 and
						cs.supertask_id = $2
					`, contestId, supertaskId).Scan(
		&supertaskInContest.SupertaskInfo.SupertaskId,
		&supertaskInContest.SupertaskInfo.VersionNumber,
		&supertaskInContest.SupertaskInfo.SupertaskName,
		&supertaskInContest.SupertaskInfo.SupertaskDesc,
		&supertaskInContest.SupertaskInfo.SupertaskLogoHref,
		&supertaskInContest.SupertaskInfo.OrderNumber,
		&supertaskInContest.SupertaskObjectJson,
	)

	if err != nil {
		return
	}

	allTasksResults, err := GetSupertaskAllTasksResults(contestId, supertaskId, userId)

	if err != nil {
		return
	}

	supertaskInContest.TaskResults = allTasksResults.TaskResults

	return
}

/*
	Возвращает результаты контеста по всем пользователям
*/

func GetContestResults(contestId int32) (results ContestResults, err error) {

	results.Errors = make([]string, 0)
	results.UserResults = make([]ContestUserResult, 0)

	/*
		Формируем заголовочную часть
	*/

	rows, err := db.Query(`
		SELECT
			tsv.supertask_name,
			tsv.tasks_num,
			tsv.max_total_score,
			tsv.max_task_score
		FROM
			t_contest_supertask tcs join
			t_supertask_version tsv on tsv.supertask_id = tcs.supertask_id and tsv.version_number = tcs.supertask_version_number
		WHERE
			tcs.contest_id = $1
		ORDER BY
			tcs.order_number
		`, contestId)

	if err != nil {
		return
	}

	defer rows.Close()

	results.MaxPossibleResult.UserLogin = "TOTAL"
	supertaskNum := 0
	for rows.Next() {
		var supertask_name, max_task_score string
		var tasks_num, max_total_score int32

		if err = rows.Scan(
			&supertask_name,
			&tasks_num,
			&max_total_score,
			&max_task_score,
		); err != nil {
			return
		}

		supertaskNum++
		results.SupertaskNames = append(results.SupertaskNames, supertask_name)
		results.TotalTasksNum += tasks_num
		results.MaxPossibleResult.TotalScore += max_total_score
		results.MaxPossibleResult.TotalTries += tasks_num
		results.MaxPossibleResult.TotalPassed += tasks_num

		results.MaxPossibleResult.SupertaskScore = append(results.MaxPossibleResult.SupertaskScore, make([]int32, tasks_num))
		results.MaxPossibleResult.SupertaskTries = append(results.MaxPossibleResult.SupertaskTries, make([]int32, tasks_num))
		results.MaxPossibleResult.SupertaskPassed = append(results.MaxPossibleResult.SupertaskPassed, make([]bool, tasks_num))

		taskMaxScore := strings.Split(max_task_score, " ")

		for i := 0; i < min(int(tasks_num), len(taskMaxScore)); i++ {
			score, e := strconv.Atoi(taskMaxScore[i])
			if e != nil {
				score = 0
			}
			results.MaxPossibleResult.SupertaskScore[supertaskNum-1][i] = int32(score)
			results.MaxPossibleResult.SupertaskTries[supertaskNum-1][i] = 1
			results.MaxPossibleResult.SupertaskPassed[supertaskNum-1][i] = true
		}
	}

	rows.Close()

	/*
		Формируем результаты по всем пользователям
	*/

	rowsUsers, err := db.Query(`
		SELECT
			r.user_id,
			tu.login,
			tcs.order_number,
			r.task_num,
			r.passed,
			r.score,
			r.tries
		FROM
			t_supertask_result r join
			t_contest_supertask tcs on tcs.contest_id = r.contest_id and tcs.supertask_id = r.supertask_id join 
			t_user tu on tu.user_id = r.user_id 
		WHERE
			r.contest_id = $1
		ORDER BY
			r.user_id
		`, contestId)

	if err != nil {
		return
	}

	defer rowsUsers.Close()

	prevUserId := int32(-1)
	userNumber := -1
	for rowsUsers.Next() {
		var user_id, order_number, task_num, score, tries int32
		var login string
		var passed bool

		if err = rowsUsers.Scan(
			&user_id,
			&login,
			&order_number,
			&task_num,
			&passed,
			&score,
			&tries,
		); err != nil {
			return
		}

		// Если новый очередной пользователь - формируем для него структуру
		if user_id != prevUserId {
			results.UserResults = append(results.UserResults, ContestUserResult{})
			userNumber++

			results.UserResults[userNumber].SupertaskScore = make([][]int32, len(results.SupertaskNames))
			results.UserResults[userNumber].SupertaskTries = make([][]int32, len(results.SupertaskNames))
			results.UserResults[userNumber].SupertaskPassed = make([][]bool, len(results.SupertaskNames))

			for i := 0; i < supertaskNum; i++ {
				results.UserResults[userNumber].SupertaskScore[i] = make([]int32, len(results.MaxPossibleResult.SupertaskScore[i]))
				results.UserResults[userNumber].SupertaskTries[i] = make([]int32, len(results.MaxPossibleResult.SupertaskScore[i]))
				results.UserResults[userNumber].SupertaskPassed[i] = make([]bool, len(results.MaxPossibleResult.SupertaskScore[i]))
			}

			results.UserResults[userNumber].UserLogin = login

			prevUserId = user_id
		}

		if int(task_num) < len(results.MaxPossibleResult.SupertaskScore[order_number]) {
			results.UserResults[userNumber].SupertaskScore[order_number][task_num] = score
			results.UserResults[userNumber].SupertaskTries[order_number][task_num] = tries
			results.UserResults[userNumber].SupertaskPassed[order_number][task_num] = passed

			results.UserResults[userNumber].TotalScore += score
			results.UserResults[userNumber].TotalTries += tries
			if passed {
				results.UserResults[userNumber].TotalPassed++
			}
		} else {
			results.Errors = append(results.Errors, fmt.Sprintf("Internal data error: task_num is larger than supertask tasks num: supertask order_number = %d, task_num = %d, user_id = %d, user_login = %s", order_number, task_num, user_id, login))
		}
	}

	sort.Slice(results.UserResults, func(i, j int) bool {
		return (results.UserResults[i].TotalScore > results.UserResults[j].TotalScore) ||
			(results.UserResults[i].TotalScore == results.UserResults[j].TotalScore && results.UserResults[i].TotalPassed > results.UserResults[j].TotalPassed) ||
			(results.UserResults[i].TotalScore == results.UserResults[j].TotalScore && results.UserResults[i].TotalPassed == results.UserResults[j].TotalPassed &&
				results.UserResults[i].TotalTries < results.UserResults[j].TotalTries)
	})

	return
}
