package admin

import (
	"codestep/db"
	"codestep/security"
	"codestep/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type CreateContestUsersRequest struct {
	ContestId    int32    `json:"contestId"`
	LoginPrefix  string   `json:"loginPrefix"`
	NumberDigits int32    `json:"numberDigits"`
	StartNumber  int32    `json:"startNumber"`
	UserCount    int32    `json:"userCount"`
	Roles        []string `json:"roles"`
}

const maxContestUsersPerRequest int32 = 1000

func parseContestUserRoles(roles []string) ([]int16, string) {
	rightIds := make([]int16, 0, len(roles))
	seen := make(map[int16]bool)

	for _, role := range roles {
		var rightId int16
		switch strings.ToUpper(strings.TrimSpace(role)) {
		case "ADMIN":
			rightId = 2
		case "JURY":
			rightId = 3
		case "PARTICIPANT":
			rightId = 4
		default:
			return nil, "Unknown role: " + role
		}
		if !seen[rightId] {
			seen[rightId] = true
			rightIds = append(rightIds, rightId)
		}
	}

	if len(rightIds) == 0 {
		return nil, "At least one role must be specified"
	}

	return rightIds, ""
}

func generateContestUserLogins(
	loginPrefix string,
	numberDigits int32,
	startNumber int32,
	userCount int32,
) []db.CreateUserAnswerRow {
	rows := make([]db.CreateUserAnswerRow, 0, userCount)

	for i := int32(0); i < userCount; i++ {
		login := fmt.Sprintf("%s%0*d", loginPrefix, numberDigits, startNumber+i)
		row := db.CreateUserAnswerRow{
			Login:    login,
			Password: utils.GeneratePassword(),
		}

		if !checkLogin(login) {
			row.IsError = true
			row.ErrorMessage = "Неправильный формат логина. Допускаются только символы a-z, 0-9, - и _. Указанный логин: " + login
		}

		rows = append(rows, row)
	}

	return rows
}

func isUserAlreadyExistsRow(row db.CreateUserAnswerRow) bool {
	return row.IsError && strings.Contains(row.ErrorMessage, "уже существует")
}

func printCreateContestUsersRowsToBytes(rows []db.CreateUserAnswerRow) []byte {
	existingRows := make([]db.CreateUserAnswerRow, 0)
	otherErrorRows := make([]db.CreateUserAnswerRow, 0)
	successRows := make([]db.CreateUserAnswerRow, 0)

	for _, row := range rows {
		switch {
		case isUserAlreadyExistsRow(row):
			existingRows = append(existingRows, row)
		case row.IsError:
			otherErrorRows = append(otherErrorRows, row)
		default:
			successRows = append(successRows, row)
		}
	}

	errNum := len(existingRows) + len(otherErrorRows)
	bytes := make([]byte, 0, 20*len(rows)+100*errNum)

	if errNum > 0 {
		bytes = append(bytes, "!!! При создании учётных записей возникли ошибки !!!\n"...)
	}

	for _, row := range existingRows {
		bytes = append(bytes, fmt.Sprintf(
			"ERROR!!! %s: пользователь уже существует, права на контест не предоставлены\n",
			row.Login,
		)...)
	}

	for _, row := range otherErrorRows {
		bytes = append(bytes, fmt.Sprintf("ERROR!!! %s\n", row.ErrorMessage)...)
	}

	for _, row := range successRows {
		bytes = append(bytes, fmt.Sprintf("%s\t%s\n", row.Login, row.Password)...)
	}

	return bytes
}

func writeCreateContestUsersJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"success": false,
		"message": message,
	}
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.Write([]byte(`{"success":false,"message":"Response marshal failed"}`))
	} else {
		w.Write(byteArr)
	}
}

/*
	Массовое создание пользователей и предоставление им прав на контест.
	Логины формируются как <loginPrefix><номер>, номер дополняется нулями слева.
	Существующим пользователям права не предоставляются.
*/
func CreateContestUsers(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	var request CreateContestUsersRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeCreateContestUsersJSONError(w, "Body reading failed", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		writeCreateContestUsersJSONError(w, "JSON decoding failed", http.StatusBadRequest)
		return
	}

	if request.ContestId <= 0 {
		writeCreateContestUsersJSONError(w, "Contest Id is required", http.StatusBadRequest)
		return
	}

	loginPrefix := strings.TrimSpace(request.LoginPrefix)
	if loginPrefix == "" {
		writeCreateContestUsersJSONError(w, "Login prefix is required", http.StatusBadRequest)
		return
	}

	if request.NumberDigits <= 0 {
		writeCreateContestUsersJSONError(w, "Number digits must be greater than zero", http.StatusBadRequest)
		return
	}

	if request.StartNumber < 0 {
		writeCreateContestUsersJSONError(w, "Start number must be non-negative", http.StatusBadRequest)
		return
	}

	if request.UserCount <= 0 {
		writeCreateContestUsersJSONError(w, "User count must be greater than zero", http.StatusBadRequest)
		return
	}

	if request.UserCount > maxContestUsersPerRequest {
		writeCreateContestUsersJSONError(w, "User count must not exceed 1000", http.StatusBadRequest)
		return
	}

	rightIds, roleErr := parseContestUserRoles(request.Roles)
	if roleErr != "" {
		writeCreateContestUsersJSONError(w, roleErr, http.StatusBadRequest)
		return
	}

	userContestRights, err := db.GetContestUserRights(userId, request.ContestId)
	if err != nil {
		writeCreateContestUsersJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userContestRights&3 == 0 {
		writeCreateContestUsersJSONError(w, "User does not have permition on the contest", http.StatusForbidden)
		return
	}

	if _, err = db.GetContest(request.ContestId); err != nil {
		writeCreateContestUsersJSONError(w, "Contest not found", http.StatusBadRequest)
		return
	}

	rows := generateContestUserLogins(
		loginPrefix,
		request.NumberDigits,
		request.StartNumber,
		request.UserCount,
	)

	_, err = db.CreateMultipleInternalUsers(rows)
	if err != nil {
		log.Println(err)
		writeCreateContestUsersJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rightsRows := make([]db.ManageUserContestRightsRow, 0)
	for _, row := range rows {
		if !row.IsError {
			rightsRows = append(rightsRows, db.ManageUserContestRightsRow{
				Login:     row.Login,
				ContestId: request.ContestId,
				Rights:    rightIds,
			})
		}
	}

	if len(rightsRows) > 0 {
		if err = db.ManageUserContestRights(rightsRows); err != nil {
			log.Println(err)
			writeCreateContestUsersJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rightsErrorsByLogin := make(map[string]string)
		for _, rightsRow := range rightsRows {
			if rightsRow.IsError {
				rightsErrorsByLogin[rightsRow.Login] = rightsRow.ErrorMessage
			}
		}

		for i, row := range rows {
			if errorMessage, ok := rightsErrorsByLogin[row.Login]; ok {
				rows[i].IsError = true
				rows[i].ErrorMessage = errorMessage
			}
		}
	}

	answerFilename := fmt.Sprintf("contest-%d-login-password.txt", request.ContestId)
	bytes := printCreateContestUsersRowsToBytes(rows)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+answerFilename+"\"")
	w.Write(bytes)
}
