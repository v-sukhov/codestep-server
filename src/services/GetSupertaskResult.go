package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type GetSupertaskResultRequest struct {
	SupertaskId int32 `json:"supertaskId"`
	TaskNum     int16 `json:"taskNum"`
	ContestId   int32 `json:"contestId"`
}

type GetSupertaskResultData struct {
	SupertaskId            int32     `json:"supertaskId"`
	TaskNum                int16     `json:"taskNum"`
	ContestId              int32     `json:"contestId"`
	SupertaskVersionNumber int32     `json:"supertaskVersionNumber"`
	Passed                 bool      `json:"passed"`
	Score                  int16     `json:"score"`
	Tries                  int16     `json:"tries"`
	SaveDTM                time.Time `json:"saveDTM"`
}

type GetSupertaskResultResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    GetSupertaskResultData `json:"data"`
}

func GetSupertaskResult(w http.ResponseWriter, r *http.Request) {
	var request GetSupertaskResultRequest
	var response GetSupertaskResultResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetSupertaskResultResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetSupertaskResultResponse{
				Success: false,
				Message: "JSON decoding failed ",
			}
		} else {
			result, err := db.GetSupertaskResult(
				request.SupertaskId,
				userId,
				request.TaskNum,
				request.ContestId,
			)

			if err != nil {
				response = GetSupertaskResultResponse{
					Success: false,
					Message: err.Error(),
				}
			} else {
				response = GetSupertaskResultResponse{
					Success: true,
					Message: "",
					Data: GetSupertaskResultData{
						SupertaskId:            result.SupertaskId,
						TaskNum:                result.TaskNum,
						SupertaskVersionNumber: result.SupertaskVersionNumber,
						ContestId:              result.ContestId,
						Passed:                 result.Passed,
						Score:                  result.Score,
						Tries:                  result.Tries,
						SaveDTM:                result.SaveDTM,
					},
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
