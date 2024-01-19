package services

import (
	"codestep/db"
	"codestep/security"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type GetSupertaskSolutionRequest struct {
	SupertaskId    int32 `json:"supertaskId"`
	TaskNum        int16 `json:"taskNum"`
	ContestId      int32 `json:"contestId"`
	SolutionStatus int16 `json:"solutionStatus"`
	SolutionKey    int16 `json:"solutionKey"`
}

type GetSupertaskSolutionData struct {
	Exists                 bool      `json:"exists"`
	SupertaskId            int32     `json:"supertaskId"`
	TaskNum                int16     `json:"taskNum"`
	SupertaskVersionNumber int32     `json:"supertaskVersionNumber"`
	ContestId              int32     `json:"contestId"`
	SolutionStatus         int16     `json:"solutionStatus"`
	SolutionKey            int16     `json:"solutionKey"`
	SaveDTM                time.Time `json:"saveDTM"`
	SolutionDesc           string    `json:"solutionDesc"`
	SolutionWorkspaceXML   string    `json:"solutionWorkspaceXML"`
}

type GetSupertaskSolutionResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    GetSupertaskSolutionData `json:"data"`
}

func GetSupertaskSolution(w http.ResponseWriter, r *http.Request) {
	var request GetSupertaskSolutionRequest
	var response GetSupertaskSolutionResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetSupertaskSolutionResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetSupertaskSolutionResponse{
				Success: false,
				Message: "JSON decoding failed ",
			}
		} else {
			solution, err := db.GetSupertaskSolution(
				request.SupertaskId,
				userId,
				request.TaskNum,
				request.ContestId,
				request.SolutionStatus,
				request.SolutionKey,
			)

			if err != nil && err != sql.ErrNoRows {
				response = GetSupertaskSolutionResponse{
					Success: false,
					Message: err.Error(),
				}
			} else {
				response = GetSupertaskSolutionResponse{
					Success: true,
					Message: "",
					Data: GetSupertaskSolutionData{
						Exists:                 (err == nil),
						SupertaskId:            solution.SupertaskId,
						TaskNum:                solution.TaskNum,
						SupertaskVersionNumber: solution.SupertaskVersionNumber,
						ContestId:              solution.ContestId,
						SolutionStatus:         solution.SolutionStatus,
						SolutionKey:            solution.SolutionKey,
						SaveDTM:                solution.SaveDTM,
						SolutionDesc:           solution.SolutionDesc,
						SolutionWorkspaceXML:   solution.SolutionWorkspaceXML,
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
