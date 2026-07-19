package contest_participation

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
	AccessWithReason AccessWithReason                `json:"accessWithReason"`
	Info             db.UserContestParticipationInfo `json:"userContestParticipationInfo"`
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response = SaveSupertaskSolutionResponse{
			Success: false,
			Message: "Body reading failed",
		}
		writeSaveSupertaskSolutionResponse(w, response)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		response = SaveSupertaskSolutionResponse{
			Success: false,
			Message: "JSON decoding failed",
		}
		writeSaveSupertaskSolutionResponse(w, response)
		return
	}

	access, info, err := CheckContestAccess(userId, request.ContestId)
	if err != nil {
		response = SaveSupertaskSolutionResponse{
			Success: false,
			Message: err.Error(),
		}
		writeSaveSupertaskSolutionResponse(w, response)
		return
	}

	if info.UserContestRights.IsEmpty() {
		response = SaveSupertaskSolutionResponse{
			Success: false,
			Message: "User does not have rights on contest",
		}
		writeSaveSupertaskSolutionResponse(w, response)
		return
	}

	accessWithReason := ResolveAccessWithReason(access, info)

	data := SaveSupertaskSolutionData{
		AccessWithReason: accessWithReason,
		Info:             info,
	}

	if accessWithReason == AccessAllowed {
		err := db.SaveSupertaskSolution(&db.SupertaskSolution{
			ContestId:              request.ContestId,
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
			writeSaveSupertaskSolutionResponse(w, response)
			return
		}
	}

	response = SaveSupertaskSolutionResponse{
		Success: true,
		Message: "OK",
		Data:    data,
	}
	writeSaveSupertaskSolutionResponse(w, response)
}

func writeSaveSupertaskSolutionResponse(w http.ResponseWriter, response SaveSupertaskSolutionResponse) {
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
