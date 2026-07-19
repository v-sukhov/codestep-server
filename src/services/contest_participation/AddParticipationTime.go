package contest_participation

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type AddParticipationTimeItem struct {
	UserId       int32 `json:"userId"`
	AddedTimeSec int32 `json:"addedTimeSec"`
}

type AddParticipationTimeRequest struct {
	ContestId int32                      `json:"contestId"`
	Rows      []AddParticipationTimeItem `json:"rows"`
}

type AddParticipationTimeRowError struct {
	RowNumber    int    `json:"rowNumber"`
	ErrorMessage string `json:"errorMessage"`
}

type AddParticipationTimeData struct {
	RowError []AddParticipationTimeRowError `json:"rowError"`
}

type AddParticipationTimeFileData struct {
	UnknownLogin []string                       `json:"unknownLogin"`
	RowError     []AddParticipationTimeRowError `json:"rowError"`
}

type AddParticipationTimeResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    AddParticipationTimeData `json:"data"`
}

type AddParticipationTimeFileResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    AddParticipationTimeFileData `json:"data"`
}

func AddParticipationTime(w http.ResponseWriter, r *http.Request) {
	var request AddParticipationTimeRequest
	var response AddParticipationTimeResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response = AddParticipationTimeResponse{
			Success: false,
			Message: "Body reading failed",
		}
		writeAddParticipationTimeResponse(w, response)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		response = AddParticipationTimeResponse{
			Success: false,
			Message: "JSON decoding failed",
		}
		writeAddParticipationTimeResponse(w, response)
		return
	}

	response = processAddParticipationTimeByUserId(userId, request.ContestId, request.Rows)
	writeAddParticipationTimeResponse(w, response)
}

func processAddParticipationTimeByUserId(
	requestUserId int32,
	contestId int32,
	items []AddParticipationTimeItem,
) AddParticipationTimeResponse {
	userContestRights, err := db.GetContestUserRights(requestUserId, contestId)
	if err != nil {
		return AddParticipationTimeResponse{
			Success: false,
			Message: err.Error(),
		}
	}
	if userContestRights&3 == 0 { // 3 = Owner + Admin
		return AddParticipationTimeResponse{
			Success: false,
			Message: "User does not have rights on operation",
		}
	}

	setRows := make([]db.SetAddedTimeRow, 0, len(items))
	rowError := make([]AddParticipationTimeRowError, 0)

	for i, item := range items {
		if item.UserId <= 0 {
			rowError = append(rowError, AddParticipationTimeRowError{
				RowNumber:    i + 1,
				ErrorMessage: "User Id is required",
			})
			continue
		}
		if item.AddedTimeSec < 0 {
			rowError = append(rowError, AddParticipationTimeRowError{
				RowNumber:    i + 1,
				ErrorMessage: "Added time must be non-negative",
			})
			continue
		}
		setRows = append(setRows, db.SetAddedTimeRow{
			UserId:       item.UserId,
			AddedTimeSec: item.AddedTimeSec,
		})
	}

	if len(setRows) > 0 {
		if err := db.SetAddedTime(contestId, setRows); err != nil {
			return AddParticipationTimeResponse{
				Success: false,
				Message: err.Error(),
			}
		}
	}

	return AddParticipationTimeResponse{
		Success: true,
		Message: "",
		Data: AddParticipationTimeData{
			RowError: rowError,
		},
	}
}

func writeAddParticipationTimeResponse(w http.ResponseWriter, response AddParticipationTimeResponse) {
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}

func writeAddParticipationTimeFileResponse(w http.ResponseWriter, response AddParticipationTimeFileResponse) {
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
