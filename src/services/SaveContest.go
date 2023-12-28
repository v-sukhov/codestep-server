package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type SaveContestRequest struct {
	ContestId       int32  `json:"contestId"`
	ContestName     string `json:"contestName"`
	ContestDesc     string `json:"contestDesc"`
	ContestLogoHref string `json:"contestLogoHref"`
}

type SaveContestData struct {
	ContestId int32 `json:"contestId"`
}

type SaveContestResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    SaveContestData `json:"data"`
}

func SaveContest(w http.ResponseWriter, r *http.Request) {

	/*
		TODO: добавить проверку прав пользователя:
			1. Вообще на операцию сохранения контеста
			2. Конкретно на этот контест
	*/

	var request SaveContestRequest
	var response SaveContestResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = SaveContestResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = SaveContestResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {

			contest := db.Contest{
				ContestId:       request.ContestId,
				ContestName:     request.ContestName,
				ContestDesc:     request.ContestDesc,
				ContestLogoHref: request.ContestLogoHref,
			}

			if err := db.SaveContest(&contest, userId); err != nil {
				response = SaveContestResponse{
					Success: false,
					Message: err.Error(),
				}
			} else {
				data := SaveContestData{
					ContestId: contest.ContestId,
				}
				response = SaveContestResponse{
					Success: true,
					Message: "OK",
					Data:    data,
				}
			}
		}
	}

	if byteArr, err := json.Marshal(response); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
