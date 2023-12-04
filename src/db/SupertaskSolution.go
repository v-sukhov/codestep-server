package db

import (
	"context"
	"time"
)

type SupertaskSolution struct {
	SupertaskSolutionId    int32
	SupertaskId            int32
	UserId                 int32
	TaskNum                int16
	SupertaskVersionNumber int32
	ContestId              int32
	SolutionStatus         int16
	SolutionKey            int16
	SolutionDesc           string
	SaveDTM                time.Time
	SolutionWorkspaceXML   string
}

/*
	Сохраняет решение задачи

	Обновляет предыдущее сохранённое решение по ключу
		SUPERTASK_ID, USER_ID, TASK_NUM, CONTEST_ID, SOLUTION_STATUS, SOLUTION_KEY
	или сохраняет новое, если такого ещё нет.
*/

func SaveSupertaskSolution(supertaskSolution *SupertaskSolution) error {
	/*
		Сначала делает попытку обновления не в рамках транзакции.
		Затем, если обновление не удалось, открывает транзакцию и снова делает попытку обновления уже в рамках транзакции и затем уже вставку, если обновление снова не удалось.
		В подавляющем большинстве случаев происходит просто первое обновление вне транзакции.
	*/

	updateQuery := `
					UPDATE T_SUPERTASK_SOLUTION
					SET
						SUPERTASK_VERSION_NUMBER = $4,
						SOLUTION_DESC = $8,
						SOLUTION_WORKSPACE_XML = $9,
						SAVE_DTM = CURRENT_TIMESTAMP
					WHERE
						SUPERTASK_ID = $1 AND
						USER_ID = $2 AND
						TASK_NUM = $3 AND
						CONTEST_ID = $5 AND
						SOLUTION_STATUS = $6 AND
						SOLUTION_KEY = $7
					`

	res, err := db.Exec(updateQuery,
		supertaskSolution.SupertaskId,
		supertaskSolution.UserId,
		supertaskSolution.TaskNum,
		supertaskSolution.SupertaskVersionNumber,
		supertaskSolution.ContestId,
		supertaskSolution.SolutionStatus,
		supertaskSolution.SolutionKey,
		supertaskSolution.SolutionDesc,
		supertaskSolution.SolutionWorkspaceXML,
	)

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
			supertaskSolution.SupertaskId,
			supertaskSolution.UserId,
			supertaskSolution.TaskNum,
			supertaskSolution.SupertaskVersionNumber,
			supertaskSolution.ContestId,
			supertaskSolution.SolutionStatus,
			supertaskSolution.SolutionKey,
			supertaskSolution.SolutionDesc,
			supertaskSolution.SolutionWorkspaceXML,
		)

		if err != nil {
			return err
		}

		rowsAffected, err := res.RowsAffected()

		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			_, err = tx.Exec(`
						INSERT INTO T_SUPERTASK_SOLUTION
						(
							supertask_id,
							user_id,
							task_num,
							supertask_version_number,
							contest_id,
							solution_status,
							solution_key,
							solution_desc,
							solution_workspace_xml
						)
						VALUES
						($1, $2, $3, $4, $5, $6, $7, $8, $9)
					`,
				supertaskSolution.SupertaskId,
				supertaskSolution.UserId,
				supertaskSolution.TaskNum,
				supertaskSolution.SupertaskVersionNumber,
				supertaskSolution.ContestId,
				supertaskSolution.SolutionStatus,
				supertaskSolution.SolutionKey,
				supertaskSolution.SolutionDesc,
				supertaskSolution.SolutionWorkspaceXML,
			)

			if err != nil {
				return err
			}
		}

		tx.Commit()

		return err
	} else {
		return nil
	}
}

/*
Возвращает решение задачи по ключу

	SUPERTASK_ID, USER_ID, TASK_NUM, CONTEST_ID, SOLUTION_STATUS, SOLUTION_KEY
*/
func GetSupertaskSolution(
	supertaskId int32,
	userId int32,
	taskNum int16,
	contestId int32,
	solutionStatus int16,
	solutionKey int16) (supertaskSolution SupertaskSolution, err error) {
	err = db.QueryRow(`
					SELECT
						solution_desc,
						save_dtm,
						supertask_version_number,
						solution_workspace_xml
					FROM t_supertask_solution
					WHERE
						supertask_id = $1 AND
						user_id = $2 AND
						task_num = $3 AND
						contest_id = $4 AND
						solution_status = $5 AND
						solution_key = $6`,
		supertaskId, userId, taskNum, contestId, solutionStatus, solutionKey).
		Scan(
			&supertaskSolution.SolutionDesc,
			&supertaskSolution.SaveDTM,
			&supertaskSolution.SupertaskVersionNumber,
			&supertaskSolution.SolutionWorkspaceXML,
		)

	supertaskSolution.SupertaskId = supertaskId
	supertaskSolution.UserId = userId
	supertaskSolution.TaskNum = taskNum
	supertaskSolution.ContestId = contestId
	supertaskSolution.SolutionStatus = solutionStatus
	supertaskSolution.SolutionKey = solutionKey

	return
}
