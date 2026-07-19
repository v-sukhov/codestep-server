package contest_participation

import (
	"bufio"
	"codestep/db"
	"codestep/security"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const maxUploadingFileSizeBytes = 1 << 20

/*
	Установка добавленного времени участникам соревнования через импорт файла.
	Структура файла (разделитель - табуляции или пробелы):

	login	added_time_sec
*/

func parseAddParticipationTimeInputFile(file io.Reader) (rows []db.SetAddedTimeByLoginRow, err error) {
	rows = make([]db.SetAddedTimeByLoginRow, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		row := db.SetAddedTimeByLoginRow{}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			row.IsError = true
			row.ErrorMessage = "Invalid row format: expected login and added time in seconds"
		} else {
			row.Login = fields[0]
			addedTimeSec, parseErr := strconv.ParseInt(fields[1], 10, 32)
			if parseErr != nil {
				row.IsError = true
				row.ErrorMessage = "Added time parsing error"
			} else if addedTimeSec < 0 {
				row.IsError = true
				row.ErrorMessage = "Added time must be non-negative"
			} else {
				row.AddedTimeSec = int32(addedTimeSec)
			}
		}

		rows = append(rows, row)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}

func AddParticipationTimeFile(w http.ResponseWriter, r *http.Request) {
	var response AddParticipationTimeFileResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if err := r.ParseMultipartForm(maxUploadingFileSizeBytes); err != nil {
		response = AddParticipationTimeFileResponse{
			Success: false,
			Message: err.Error(),
		}
		writeAddParticipationTimeFileResponse(w, response)
		return
	}

	contestId, err := strconv.ParseInt(r.FormValue("contestId"), 10, 32)
	if err != nil {
		response = AddParticipationTimeFileResponse{
			Success: false,
			Message: "Contest Id parsing error",
		}
		writeAddParticipationTimeFileResponse(w, response)
		return
	}

	file, _, err := r.FormFile("participationTimeList")
	if err != nil {
		response = AddParticipationTimeFileResponse{
			Success: false,
			Message: err.Error(),
		}
		writeAddParticipationTimeFileResponse(w, response)
		return
	}
	defer file.Close()

	rows, err := parseAddParticipationTimeInputFile(file)
	if err != nil {
		response = AddParticipationTimeFileResponse{
			Success: false,
			Message: err.Error(),
		}
		writeAddParticipationTimeFileResponse(w, response)
		return
	}

	response = processAddParticipationTimeByLogin(userId, int32(contestId), rows)
	writeAddParticipationTimeFileResponse(w, response)
}

func processAddParticipationTimeByLogin(
	requestUserId int32,
	contestId int32,
	rows []db.SetAddedTimeByLoginRow,
) AddParticipationTimeFileResponse {
	userContestRights, err := db.GetContestUserRights(requestUserId, contestId)
	if err != nil {
		return AddParticipationTimeFileResponse{
			Success: false,
			Message: err.Error(),
		}
	}
	if userContestRights&3 == 0 {
		return AddParticipationTimeFileResponse{
			Success: false,
			Message: "User does not have rights on operation",
		}
	}

	unknownLogins, err := db.SetAddedTimeByLogin(contestId, rows)
	if err != nil {
		return AddParticipationTimeFileResponse{
			Success: false,
			Message: err.Error(),
		}
	}

	rowError := make([]AddParticipationTimeRowError, 0)
	for i, row := range rows {
		if row.IsError {
			rowError = append(rowError, AddParticipationTimeRowError{
				RowNumber:    i + 1,
				ErrorMessage: row.ErrorMessage,
			})
		}
	}

	return AddParticipationTimeFileResponse{
		Success: true,
		Message: "",
		Data: AddParticipationTimeFileData{
			UnknownLogin: unknownLogins,
			RowError:     rowError,
		},
	}
}
