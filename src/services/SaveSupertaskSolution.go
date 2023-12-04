package services

/*
	Сохранение решения задачи
*/

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type SaveSupertaskSolutionRequest struct {
	SupertaskId            int32  `json:"supertaskId"`
	TaskNum                int16  `json:"taskNum"`
	SupertaskVersionNumber int32  `json:"supertaskVersionNumber"`
	ContestId              int32  `json:"contestId"`
	SolutionStatus         int16  `json:"solutionStatus"`
	SolutionKey            int16  `json:"solutionKey"`
	SolutionDesc           string `json:"solutionDesc"`
	SolutionWorkspaceXML   string `json:"solutionWorkspaceXML"`
}

type SaveSupertaskSolutionData struct {
}

type SaveSupertaskSolutionResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    SaveSupertaskSolutionData `json:"data"`
}

func SaveSupertaskSolution(w http.ResponseWriter, r *http.Request) {
	var request SaveSupertaskSolutionRequest
	var response SaveSupertaskSolutionResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = SaveSupertaskSolutionResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = SaveSupertaskSolutionResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {
			err := db.SaveSupertaskSolution(&db.SupertaskSolution{
				SupertaskId:            request.SupertaskId,
				UserId:                 userId,
				SupertaskVersionNumber: request.SupertaskVersionNumber,
				SolutionStatus:         request.SolutionStatus,
				TaskNum:                request.TaskNum,
				SolutionKey:            request.SolutionKey,
				SolutionDesc:           request.SolutionDesc,
				SolutionWorkspaceXML:   request.SolutionWorkspaceXML,
			})

			if err != nil {
				response = SaveSupertaskSolutionResponse{
					Success: false,
					Message: err.Error(),
				}
			} else {
				response = SaveSupertaskSolutionResponse{
					Success: true,
					Message: "",
					Data:    SaveSupertaskSolutionData{},
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
