package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

/*
	Сохраняет результат тестирования суперзадачи
*/

type SaveSupertaskResultRequest struct {
	SupertaskId            int32 `json:"supertaskId"`
	TaskNum                int16 `json:"taskNum"`
	ContestId              int32 `json:"contestId"`
	SupertaskVersionNumber int32 `json:"supertaskVersionNumber"`
	Passed                 bool  `json:"passed"`
	Score                  int16 `json:"score"`
	Tries                  int16 `json:"tries"`
}

type SaveSupertaskResultData struct {
}

type SaveSupertaskResultResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    SaveSupertaskResultData `json:"data"`
}

func SaveSupertaskResult(w http.ResponseWriter, r *http.Request) {

	var request SaveSupertaskResultRequest
	var response SaveSupertaskResultResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = SaveSupertaskResultResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = SaveSupertaskResultResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {

			supertaskResult := db.SupertaskResult{
				SupertaskId:            request.SupertaskId,
				UserId:                 userId,
				TaskNum:                request.TaskNum,
				ContestId:              request.ContestId,
				SupertaskVersionNumber: request.SupertaskVersionNumber,
				Passed:                 request.Passed,
				Score:                  request.Score,
				Tries:                  request.Tries,
			}

			if err := db.SaveSupertaskResult(&supertaskResult); err != nil {
				response = SaveSupertaskResultResponse{
					Success: false,
					Message: err.Error(),
				}
			} else {
				data := SaveSupertaskResultData{}
				response = SaveSupertaskResultResponse{
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
