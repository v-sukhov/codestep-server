package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GetSupertaskAllTasksResultsRequest struct {
	SupertaskId int32 `json:"supertaskId"`
	ContestId   int32 `json:"contestId"`
}

type GetSupertaskAllTasksResultsData struct {
	SupertaskId  int32           `json:"supertaskId"`
	ContestId    int32           `json:"contestId"`
	TasksResults []db.TaskResult `json:"tasksResults"`
}

type GetSupertaskAllTasksResultsResponse struct {
	Success bool                            `json:"success"`
	Message string                          `json:"message"`
	Data    GetSupertaskAllTasksResultsData `json:"data"`
}

func GetSupertaskAllTasksResults(w http.ResponseWriter, r *http.Request) {
	var request GetSupertaskAllTasksResultsRequest
	var response GetSupertaskAllTasksResultsResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetSupertaskAllTasksResultsResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetSupertaskAllTasksResultsResponse{
				Success: false,
				Message: "JSON decoding failed ",
			}
		} else {
			allResults, err := db.GetSupertaskAllTasksResults(
				request.SupertaskId,
				userId,
				request.ContestId,
			)

			if err != nil {
				response = GetSupertaskAllTasksResultsResponse{
					Success: false,
					Message: err.Error(),
				}
			} else {
				response = GetSupertaskAllTasksResultsResponse{
					Success: true,
					Message: "",
					Data: GetSupertaskAllTasksResultsData{
						SupertaskId:  allResults.SupertaskId,
						ContestId:    allResults.ContestId,
						TasksResults: allResults.TaskResults,
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
