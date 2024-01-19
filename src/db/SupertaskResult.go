package db

import (
	"context"
	"database/sql"
	"time"
)

type SupertaskResult struct {
	SupertaskResultId      int32
	SupertaskId            int32
	UserId                 int32
	TaskNum                int16
	ContestId              int32
	SupertaskVersionNumber int32
	Passed                 bool
	Score                  int16
	Tries                  int16
	SaveDTM                time.Time
}

type TaskResult struct {
	TaskNum                int16     `json:"taskNum"`
	SupertaskVersionNumber int32     `json:"supertaskVersionNumber"`
	Passed                 bool      `json:"passed"`
	Score                  int16     `json:"score"`
	Tries                  int16     `json:"tries"`
	SaveDTM                time.Time `json:"saveDTM"`
}

type SupertaskAllTasksResults struct {
	SupertaskId int32
	UserId      int32
	ContestId   int32
	TaskResults []TaskResult
}

/*
	Сохраняет результат по задаче

	Обновляет предыдущий результат по ключу
		SUPERTASK_ID, USER_ID, TASK_NUM, CONTEST_ID
	или создаёт новую запись, если такой ещё нет
*/

func SaveSupertaskResult(supertaskResult *SupertaskResult) error {

	/*
		Пытается сначала выполнить простой UPDATE без транзакции - этим запросом чаще всего работа и будет ограничиваться.
		Если не получилось - открывает транзакцию и снова пытается выполнить UPDATE, и уже если опять не получилось - делает вставку.
		Реального обновления не должно происходить, если количество попыток в обновлении меньше последнего записанного.
	*/

	updateQuery := `
					UPDATE T_SUPERTASK_RESULT T
					SET
						SUPERTASK_VERSION_NUMBER = (CASE WHEN T.TRIES < $8 THEN $5 ELSE SUPERTASK_VERSION_NUMBER END),
						PASSED = (CASE WHEN T.TRIES < $8 THEN $6 ELSE PASSED END),
						SCORE = (CASE WHEN T.TRIES < $8 THEN $7 ELSE SCORE END),
						TRIES = (CASE WHEN T.TRIES < $8 THEN $8 ELSE TRIES END),
						SAVE_DTM = (CASE WHEN T.TRIES < $8 THEN CURRENT_TIMESTAMP ELSE SAVE_DTM END)
					WHERE
						SUPERTASK_ID = $1 AND
						USER_ID = $2 AND
						TASK_NUM = $3 AND
						CONTEST_ID = $4
					`

	res, err := db.Exec(updateQuery,
		supertaskResult.SupertaskId,
		supertaskResult.UserId,
		supertaskResult.TaskNum,
		supertaskResult.ContestId,
		supertaskResult.SupertaskVersionNumber,
		supertaskResult.Passed,
		supertaskResult.Score,
		supertaskResult.Tries)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		ctx := context.TODO()
		tx, err := db.BeginTx(ctx, nil)

		if err != nil {
			return err
		}
		defer tx.Rollback()

		res, err := tx.Exec(updateQuery,
			supertaskResult.SupertaskId,
			supertaskResult.UserId,
			supertaskResult.TaskNum,
			supertaskResult.ContestId,
			supertaskResult.SupertaskVersionNumber,
			supertaskResult.Passed,
			supertaskResult.Score,
			supertaskResult.Tries)

		if err != nil {
			return err
		}

		rowsAffected, err := res.RowsAffected()

		if rowsAffected == 0 {
			_, err = tx.Exec(`
						INSERT INTO T_SUPERTASK_RESULT
						(
							supertask_id,
							user_id,
							task_num,
							contest_id,
							supertask_version_number,
							passed,
							score,
							tries
						)
						VALUES
						($1, $2, $3, $4, $5, $6, $7, $8)
					`,
				supertaskResult.SupertaskId,
				supertaskResult.UserId,
				supertaskResult.TaskNum,
				supertaskResult.ContestId,
				supertaskResult.SupertaskVersionNumber,
				supertaskResult.Passed,
				supertaskResult.Score,
				supertaskResult.Tries,
			)

			if err != nil {
				return err
			}
		}

		tx.Commit()

		return nil

	} else {
		return nil
	}
}

/*
	Возвращает результат решения задачи по ключу
		SUPERTASK_ID, USER_ID, TASK_NUM, CONTEST_ID
	Если сохранённого результата по ключу нет, возвращает нулевой результат
*/

func GetSupertaskResult(
	supertaskId int32,
	userId int32,
	taskNum int16,
	contestId int32) (result SupertaskResult, err error) {

	err = db.QueryRow(`
			SELECT
				supertask_result_id,
				supertask_version_number,
				passed,
				score,
				tries,
				save_dtm
			FROM
				t_supertask_result
			WHERE
				supertask_id = $1 and
				user_id = $2 and
				task_num = $3 and
				contest_id = $4
		`, supertaskId, userId, taskNum, contestId).Scan(
		&result.SupertaskResultId,
		&result.SupertaskVersionNumber,
		&result.Passed,
		&result.Score,
		&result.Tries,
		&result.SaveDTM,
	)

	if err == sql.ErrNoRows {
		result.Passed = false
		result.Score = 0
		result.Tries = 0
		err = nil
	}

	result.SupertaskId = supertaskId
	result.UserId = userId
	result.TaskNum = taskNum
	result.ContestId = contestId

	return
}

/*
	Возвращает результаты решения всех задач суперзадачи по ключу
		SUPERTASK_ID, USER_ID, CONTEST_ID
*/

func GetSupertaskAllTasksResults(
	contestId int32,
	supertaskId int32,
	userId int32,
) (allTasksResult SupertaskAllTasksResults, err error) {

	rows, err := db.Query(`
			SELECT
				task_num,
				supertask_version_number,
				passed,
				score,
				tries,
				save_dtm
			FROM
				t_supertask_result
			WHERE
				supertask_id = $1 and
				user_id = $2 and
				contest_id = $3
		`, supertaskId, userId, contestId)

	defer rows.Close()

	if err != nil {
		return
	}

	allTasksResult.TaskResults = make([]TaskResult, 0)

	for rows.Next() {
		var taskResult TaskResult
		if err = rows.Scan(
			&taskResult.TaskNum,
			&taskResult.SupertaskVersionNumber,
			&taskResult.Passed,
			&taskResult.Score,
			&taskResult.Tries,
			&taskResult.SaveDTM,
		); err != nil {
			return
		}
		allTasksResult.TaskResults = append(allTasksResult.TaskResults, taskResult)
	}

	allTasksResult.SupertaskId = supertaskId
	allTasksResult.UserId = userId
	allTasksResult.ContestId = contestId

	return
}
