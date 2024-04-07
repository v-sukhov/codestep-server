package services

import (
	"bufio"
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ManageMultipleUsersContestRightsRowError struct {
	RowNumber    int    `json:"rowNumber"`
	ErrorMessage string `json:"errorMessage"`
}

type ManageMultipleUsersContestRightsData struct {
	RowError []ManageMultipleUsersContestRightsRowError `json:"rowError"`
}

type ManageMultipleUsersContestRightsResponse struct {
	Success bool                                 `json:"success"`
	Message string                               `json:"message"`
	Data    ManageMultipleUsersContestRightsData `json:"data"`
}

/*
	Управление правами пользователей на контест через импорт файла
	Структура файла (разделитель - табуляции или пробелы):

	login	contest_id	[revoke|admin|jury|participant]*

	revoke - отнимаются все права, с опцией revoke недопустимо задавать другие варианты прав
	По умолчанию (если права не перечислены, считается, что заданы права participant)
	Если один и тот же логин встречается в файле несколько раз, то каждая комбинация логин-контест-тип прав учитывается один раз
*/

/*
	TODO:

	Добавить проверку на существование контеста contest_id
*/

func parseManageUserContestRightsInputFile(file io.Reader) (rows []db.ManageUserContestRightsRow) {
	rows = make([]db.ManageUserContestRightsRow, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var row db.ManageUserContestRightsRow
		row.Rights = make([]int16, 0)
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		for i, token := range fields {
			if i == 0 {
				row.Login = token
			} else if i == 1 {
				id, err := strconv.Atoi(token)
				if err == nil {
					row.ContestId = int32(id)
				} else {
					row.IsError = true
					row.ErrorMessage = "Contest Id parsing error"
				}
			} else if token == "revoke" {
				row.IsRevoke = true
				if len(row.Rights) > 0 {
					row.IsError = true
					row.ErrorMessage = "Setting right type with revoke directive is not allowed"
				}
			} else {
				if row.IsRevoke {
					row.IsError = true
					row.ErrorMessage = "Setting right type with revoke directive is not allowed"
				} else if token == "admin" {
					row.Rights = append(row.Rights, 2)
				} else if token == "jury" {
					row.Rights = append(row.Rights, 3)
				} else if token == "participant" {
					row.Rights = append(row.Rights, 4)
				} else {
					row.IsError = true
					row.ErrorMessage = "Unknown right type: " + token
				}
			}
		}
		if len(row.Rights) == 0 && !row.IsError {
			// По умолчанию добавляем права участника
			row.Rights = append(row.Rights, 4)
		}

		rows = append(rows, row)
	}

	return rows
}

func ManageMultipleUsersContestRights(w http.ResponseWriter, r *http.Request) {

	var response ManageMultipleUsersContestRightsResponse

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	userRights, err := db.GetUserRights(userId)
	if err != nil {
		response = ManageMultipleUsersContestRightsResponse{
			Success: false,
			Message: err.Error(),
		}
	} else if !userRights.IsAdmin {
		response = ManageMultipleUsersContestRightsResponse{
			Success: false,
			Message: "User does not have permition on operation",
		}
	} else {
		var formFileName = "contestUserList"

		r.ParseMultipartForm(MAX_UPLOADING_FILE_SIZE_BYTES)

		file, _, err := r.FormFile(formFileName)

		if err != nil {
			response = ManageMultipleUsersContestRightsResponse{
				Success: false,
				Message: err.Error(),
			}
		} else {
			defer file.Close()
			rows := parseManageUserContestRightsInputFile(file)
			err = db.ManageUserContestRights(rows)
			if err != nil {
				response = ManageMultipleUsersContestRightsResponse{
					Success: false,
					Message: err.Error(),
				}
			} else {
				var rowError []ManageMultipleUsersContestRightsRowError
				rowError = make([]ManageMultipleUsersContestRightsRowError, 0)

				for i, row := range rows {
					if row.IsError {
						rowError = append(rowError, ManageMultipleUsersContestRightsRowError{
							RowNumber:    i + 1,
							ErrorMessage: row.ErrorMessage,
						})
					}
				}

				response = ManageMultipleUsersContestRightsResponse{
					Success: true,
					Message: "",
					Data: ManageMultipleUsersContestRightsData{
						RowError: rowError,
					},
				}
			}
		}

	}

	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}

}
