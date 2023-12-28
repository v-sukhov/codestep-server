package db

import (
	"context"
	"database/sql"
	"time"
)

type SupertaskVersion struct {
	SupertaskId         int32
	VersionNumber       int32
	ParentVersionNumber int32
	Commited            bool
	AuthorUserId        int32
	CommitMessage       string
	SaveDTM             time.Time
	SupertaskName       string
	SupertaskDesc       string
	SupertaskObjectJson string
}

type SupertaskLastVersionInfo struct {
	SupertaskId            int32     `json:"supertaskId"`
	SupertaskRightTypeId   int32     `json:"supertaskRightTypeId"`
	SupertaskRightTypeName string    `json:"supertaskRightTypeName"`
	OwnerUserId            int32     `json:"ownerUserId"`
	OwnerUserName          string    `json:"ownerUserName"`
	LastVersionNumber      int32     `json:"lastVersionNumber"`
	SupertaskCreateDTM     time.Time `json:"supertaskCreateDTM"`
	SupertaskName          string    `json:"supertaskName"`
	SupertaskDesc          string    `json:"supertaskDesc"`
	VersionAuthorUserName  string    `json:"versionAuthorUserName"`
	ParentVersionNumber    string    `json:"parentVersionNumber"`
	Commited               bool      `json:"commited"`
	CommitMessage          string    `json:"commitMessage"`
	SaveDTM                time.Time `json:"saveDTM"`
}

/*
*
Сохраняет версию суперзадачи с обработкой различных ситуаций:
  - создание новой суперзадачи
  - сохранение существующей
  - сохранение закоммиченной версии
  - сохранение незакоммиченной версии
  - VersionNumber -
    в случае коммита - вычисляемое поле: вычисляется и записывается в структуру, входящее значение не используется
    если нет коммита - записывается входящее значение
*/
func SaveSupertaskVersion(supertaskVersion *SupertaskVersion) error {

	if supertaskVersion.SupertaskId == 0 {
		// Новая задача - создаём необходимые записи в таблицах T_SUPERTASK, T_SUPERTASK_USER_RIGHT, T_SUPERTASK_LAST_VERSION

		var supertaskId int32

		err := db.QueryRow("INSERT INTO T_SUPERTASK(SUPERTASK_STATUS_ID) VALUES(1) RETURNING SUPERTASK_ID").Scan(&supertaskId)

		if err != nil {
			return err
		}

		supertaskVersion.SupertaskId = supertaskId

		_, err = db.Exec("INSERT INTO T_SUPERTASK_USER_RIGHT(SUPERTASK_ID, USER_ID, SUPERTASK_RIGHT_TYPE_ID) VALUES($1, $2, 1)", supertaskId, supertaskVersion.AuthorUserId)

		if err != nil {
			return err
		}

		_, err = db.Exec("INSERT INTO T_SUPERTASK_LAST_VERSION(SUPERTASK_ID, LAST_VERSION_NUMBER) VALUES($1, 0)", supertaskId)

		if err != nil {
			return err
		}

	}

	// Теперь сохраняем объект supertaskVersion в T_SUPERTASK_VERSION - в объекте уже в любом случае есть supertaskId
	// Отдельно обрабатываем коммит и простое сохранение

	if supertaskVersion.Commited {
		// Создаём новую закоммиченную версию задачи: увеличиваем номер максимальной версии и сохраняем в данную версию как закоммиченную.
		// Рабочую версию данного пользователя удаляем.
		//var ctx context.Context
		//tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: 0, ReadOnly: false})
		ctx := context.TODO()
		tx, err := db.BeginTx(ctx, nil)

		if err != nil {
			return err
		}
		defer tx.Rollback()

		var lastVersion int32
		err = tx.QueryRow("SELECT LAST_VERSION_NUMBER FROM T_SUPERTASK_LAST_VERSION WHERE SUPERTASK_ID = $1 FOR UPDATE", supertaskVersion.SupertaskId).Scan(&lastVersion)

		if err != nil {
			return err
		}

		newVersionNumber := lastVersion + 1
		supertaskVersion.VersionNumber = newVersionNumber

		_, err = tx.Exec("UPDATE T_SUPERTASK_LAST_VERSION SET LAST_VERSION_NUMBER = $1 WHERE SUPERTASK_ID = $2", newVersionNumber, supertaskVersion.SupertaskId)

		if err != nil {
			return err
		}

		_, err = tx.Exec(`INSERT INTO T_SUPERTASK_VERSION
						  (
							SUPERTASK_ID,
							VERSION_NUMBER,
							PARENT_VERSION_NUMBER,
							COMMITED,
							AUTHOR_USER_ID,
							COMMIT_MESSAGE,
							SAVE_DTM,
							SUPERTASK_NAME,
							SUPERTASK_DESC,
							SUPERTASK_OBJECT_JSON
						  )
						  VALUES
						  (
							$1, $2, $3, $4, $5, $6, NOW(), $7, $8, $9
						  )`,
			supertaskVersion.SupertaskId,
			newVersionNumber,
			supertaskVersion.ParentVersionNumber,
			supertaskVersion.Commited,
			supertaskVersion.AuthorUserId,
			supertaskVersion.CommitMessage,
			supertaskVersion.SupertaskName,
			supertaskVersion.SupertaskDesc,
			supertaskVersion.SupertaskObjectJson,
		)

		if err != nil {
			return err
		}

		_, err = tx.Exec("DELETE FROM T_SUPERTASK_VERSION WHERE SUPERTASK_ID = $1 AND AUTHOR_USER_ID = $2 AND COMMITED = FALSE",
			supertaskVersion.SupertaskId, supertaskVersion.AuthorUserId)

		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

	} else {
		// Просто сохраняем текущую рабочую версию данного пользователя

		ctx := context.TODO()
		tx, err := db.BeginTx(ctx, nil)

		if err != nil {
			return err
		}
		defer tx.Rollback()

		if _, err := tx.Exec(`CREATE TEMP TABLE TMP_NEW_SUPERTASK_VERSION
				(
				SUPERTASK_ID INTEGER,
				VERSION_NUMBER INTEGER,
				PARENT_VERSION_NUMBER INTEGER,
				COMMITED BOOLEAN,
				AUTHOR_USER_ID INTEGER,
				COMMIT_MESSAGE VARCHAR(512),
				SUPERTASK_NAME VARCHAR(256),
				SUPERTASK_DESC VARCHAR(256),
				SUPERTASK_OBJECT_JSON TEXT
				)
		`); err != nil {
			return err
		}

		if _, err := tx.Exec(`INSERT INTO TMP_NEW_SUPERTASK_VERSION (
				SUPERTASK_ID,
				VERSION_NUMBER,
				PARENT_VERSION_NUMBER,
				COMMITED,
				AUTHOR_USER_ID,
				COMMIT_MESSAGE,
				SUPERTASK_NAME,
				SUPERTASK_DESC,
				SUPERTASK_OBJECT_JSON
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			supertaskVersion.SupertaskId,
			supertaskVersion.VersionNumber,
			supertaskVersion.ParentVersionNumber,
			supertaskVersion.Commited,
			supertaskVersion.AuthorUserId,
			supertaskVersion.CommitMessage,
			supertaskVersion.SupertaskName,
			supertaskVersion.SupertaskDesc,
			supertaskVersion.SupertaskObjectJson,
		); err != nil {
			return err
		}

		if _, err := tx.Exec(`	MERGE INTO T_SUPERTASK_VERSION sv
		USING TMP_NEW_SUPERTASK_VERSION as new_data				
		ON sv.SUPERTASK_ID = new_data.SUPERTASK_ID AND sv.AUTHOR_USER_ID = new_data.AUTHOR_USER_ID AND sv.COMMITED = FALSE
		WHEN MATCHED THEN
		UPDATE SET
			VERSION_NUMBER = new_data.VERSION_NUMBER, 
			PARENT_VERSION_NUMBER = new_data.PARENT_VERSION_NUMBER,
			COMMITED = new_data.COMMITED,
			AUTHOR_USER_ID = new_data.AUTHOR_USER_ID,
			COMMIT_MESSAGE = new_data.COMMIT_MESSAGE,
			SAVE_DTM = NOW(),
			SUPERTASK_NAME = new_data.SUPERTASK_NAME,
			SUPERTASK_DESC = new_data.SUPERTASK_DESC,
			SUPERTASK_OBJECT_JSON = new_data.SUPERTASK_OBJECT_JSON
		WHEN NOT MATCHED THEN 
		INSERT 
			(
				SUPERTASK_ID, 
				VERSION_NUMBER, 
				PARENT_VERSION_NUMBER,
				COMMITED,
				AUTHOR_USER_ID,
				COMMIT_MESSAGE,
				SAVE_DTM,
				SUPERTASK_NAME,
				SUPERTASK_DESC,
				SUPERTASK_OBJECT_JSON
			) 
			VALUES 
			(
				new_data.SUPERTASK_ID, 
				new_data.VERSION_NUMBER, 
				new_data.PARENT_VERSION_NUMBER,
				new_data.COMMITED,
				new_data.AUTHOR_USER_ID,
				new_data.COMMIT_MESSAGE,
				NOW(),
				new_data.SUPERTASK_NAME,
				new_data.SUPERTASK_DESC,
				new_data.SUPERTASK_OBJECT_JSON
			)
		`); err != nil {
			return err
		}

		if _, err := tx.Exec(`DROP TABLE TMP_NEW_SUPERTASK_VERSION`); err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func GetSupertaskVersion(supertaskId int32, supertaskVersionNumber int32) (supertaskVersion SupertaskVersion, err error) {

	err = db.QueryRow(`
					SELECT
						supertask_id,
						version_number,
						parent_version_number,
						commited,
						author_user_id,
						commit_message,
						save_dtm,
						supertask_name,
						supertask_desc,
						supertask_object_json
					FROM
						t_supertask_version
					WHERE
						supertask_id = $1 and version_number = $2`, supertaskId, supertaskVersionNumber).
		Scan(
			&supertaskVersion.SupertaskId,
			&supertaskVersion.VersionNumber,
			&supertaskVersion.ParentVersionNumber,
			&supertaskVersion.Commited,
			&supertaskVersion.AuthorUserId,
			&supertaskVersion.CommitMessage,
			&supertaskVersion.SaveDTM,
			&supertaskVersion.SupertaskName,
			&supertaskVersion.SupertaskDesc,
			&supertaskVersion.SupertaskObjectJson)

	return
}

func GetSupertaskWorkVersion(supertaskId int32, AuthorUserId int32) (supertaskVersion SupertaskVersion, err error) {

	err = db.QueryRow(`
					SELECT
						supertask_id,
						version_number,
						parent_version_number,
						commited,
						author_user_id,
						commit_message,
						save_dtm,
						supertask_name,
						supertask_desc,
						supertask_object_json
					FROM
						t_supertask_version
					WHERE
						supertask_id = $1 and commited = false and author_user_id = $2`, supertaskId, AuthorUserId).
		Scan(
			&supertaskVersion.SupertaskId,
			&supertaskVersion.VersionNumber,
			&supertaskVersion.ParentVersionNumber,
			&supertaskVersion.Commited,
			&supertaskVersion.AuthorUserId,
			&supertaskVersion.CommitMessage,
			&supertaskVersion.SaveDTM,
			&supertaskVersion.SupertaskName,
			&supertaskVersion.SupertaskDesc,
			&supertaskVersion.SupertaskObjectJson)

	return
}

func GetSupertaskUserRight(supertaskId int32, userId int32) (supertaskUserRight int32, err error) {
	// 0 - no right, > 0 - according to t_supertask_right_type

	err = db.QueryRow(`SELECT supertask_right_type_id FROM t_supertask_user_right WHERE supertask_id = $1 and user_id = $2`, supertaskId, userId).Scan(&supertaskUserRight)

	if err == sql.ErrNoRows {
		err = nil
		supertaskUserRight = 0
	}

	return
}

func GetSupertaskOwnerId(supertaskId int32) (ownerUserId int32, err error) {
	err = db.QueryRow(`SELECT user FROM t_supertask_user_right WHERE supertask_id = $1 and supertask_right_type_id = 1`, supertaskId).Scan(&ownerUserId)

	if err == sql.ErrNoRows {
		err = nil
		ownerUserId = 0
	}

	return
}

/*
*
Сформировать список суперзадач, доступных пользователю

TODO: подумать, что делать, если предоставили пользователю права на задачу, у которой нет ещё ни одной закоммиченной версии
*/
func GetUserSupertaskList(userId int32) ([]SupertaskLastVersionInfo, error) {
	rows, err := db.Query(`
					select
						ur.supertask_id,
						tsrt.supertask_right_type_id,
						tsrt.supertask_right_type_name,
						uo.user_id as owner_user_id,
						uo.display_name as owner_user_name,
						tslv.last_version_number,
						ts.create_dtm as supertask_create_dtm,
						coalesce(stu.supertask_name, stc.supertask_name, 'Новая задача') as supertask_name,
						coalesce(stu.supertask_desc, stc.supertask_desc, 'В задаче отсутствуют закоммиченные версии и Ваша рабочая версия') as supertask_desc,
						coalesce(u.display_name, uc.display_name, '') as version_author_user_name,
						coalesce(stu.parent_version_number , stc.parent_version_number, 0) as parent_version_number,
						coalesce(stu.commited, stc.commited, false) as commited,
						coalesce(stu.commit_message, stc.commit_message, '') as commit_message,
						coalesce(stu.save_dtm, stc.save_dtm, ts.create_dtm) as save_dtm
					from 
						t_supertask_user_right ur join
						t_supertask_right_type tsrt on tsrt.supertask_right_type_id = ur.supertask_right_type_id join
						t_supertask ts on ts.supertask_id = ur.supertask_id join 
						t_supertask_user_right uro on uro.supertask_id = ur.supertask_id and uro.supertask_right_type_id = 1 join
						t_user uo on uo.user_id = uro.user_id join
						t_supertask_last_version tslv on tslv.supertask_id = ur.supertask_id left join 
						t_supertask_version stc on stc.supertask_id = ur.supertask_id and (stc.version_number = tslv.last_version_number and stc.version_number > 0) left join
						t_supertask_version stu on stu.supertask_id = ur.supertask_id and stu.version_number = 0 and stu.author_user_id = ur.user_id left join 
						t_user uc on uc.user_id = stc.author_user_id left join 
						t_user u on u.user_id = stu.author_user_id 
					where
						ts.supertask_status_id = 1 and ur.supertask_right_type_id > 0 and ur.user_id = $1
					order by 
						coalesce(stu.save_dtm, stc.save_dtm) desc	
				`, userId)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var supertaskList []SupertaskLastVersionInfo

	for rows.Next() {
		var supertaskInfo SupertaskLastVersionInfo
		if err := rows.Scan(
			&supertaskInfo.SupertaskId,
			&supertaskInfo.SupertaskRightTypeId,
			&supertaskInfo.SupertaskRightTypeName,
			&supertaskInfo.OwnerUserId,
			&supertaskInfo.OwnerUserName,
			&supertaskInfo.LastVersionNumber,
			&supertaskInfo.SupertaskCreateDTM,
			&supertaskInfo.SupertaskName,
			&supertaskInfo.SupertaskDesc,
			&supertaskInfo.VersionAuthorUserName,
			&supertaskInfo.ParentVersionNumber,
			&supertaskInfo.Commited,
			&supertaskInfo.CommitMessage,
			&supertaskInfo.SaveDTM,
		); err != nil {
			return supertaskList, err
		}
		supertaskList = append(supertaskList, supertaskInfo)
	}

	if err = rows.Err(); err != nil {
		return supertaskList, err
	}

	return supertaskList, nil
}
