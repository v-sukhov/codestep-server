package services

import (
	"bufio"
	"codestep/db"
	"codestep/security"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var validLoginRegexp = regexp.MustCompile(`[a-z\d-_]+`)

func checkLogin(login string) bool {
	return validLoginRegexp.MatchString(login)
}

func parseInputFile(file io.Reader) (rows []db.AnswerRow, errNum int) {
	rows = make([]db.AnswerRow, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var row db.AnswerRow
		login := strings.TrimSpace(scanner.Text())
		if len(login) > 0 {
			if checkLogin(login) {
				row.Login = login
				row.Password = generatePassword()
			} else {
				row.IsError = true
				row.ErrorMessage = "Неправильный формат логина. Допускаются только символы a-z, 0-9, - и _. Указанный логин: " + login
				errNum++
			}
			rows = append(rows, row)
		}
	}

	return
}

func printLoginPasswordRowsToBytes(rows []db.AnswerRow, errNum int) []byte {
	var bytes []byte
	bytes = make([]byte, 0, 20*len(rows)+100*errNum)

	if errNum > 0 {
		bytes = append(bytes, "!!! При создании учётных записей возникли ошибки !!!\n"...)
	}

	for _, row := range rows {
		var s string
		if row.IsError {
			s = fmt.Sprintf("ERROR!!! %s\n", row.ErrorMessage)
		} else {
			s = fmt.Sprintf("%s\t%s\n", row.Login, row.Password)
		}
		bytes = append(bytes, s...)
	}

	return bytes
}

/*
	Input file structure:

	login
*/

func CreateMultipleUsers(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	userRights, err := db.GetUserRights(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !userRights.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var formFileName = "userList"

	r.ParseMultipartForm(MAX_UPLOADING_FILE_SIZE_BYTES)

	file, handler, err := r.FormFile(formFileName)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	rows, errNum := parseInputFile(file)

	err, dbErrRowsNum := db.CreateMultipleInternalUsers(rows)

	errNum += dbErrRowsNum

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes := printLoginPasswordRowsToBytes(rows, errNum)

	answerFilename := "login-password-" + handler.Filename

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+answerFilename+"\"")
	w.Write(bytes)
}
