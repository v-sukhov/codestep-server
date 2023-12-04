package db

import (
	"context"
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
	*/

	updateQuery := `
					UPDATE T_SUPERTASK_RESULT
					SET
						SUPERTASK_VERSION_NUMBER = $5,
						PASSED = $6,
						SCORE = $7,
						TRIES = $8,
						SAVE_DTM = CURRENT_TIMESTAMP
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

	result.SupertaskId = supertaskId
	result.UserId = userId
	result.TaskNum = taskNum
	result.ContestId = contestId

	return
}
