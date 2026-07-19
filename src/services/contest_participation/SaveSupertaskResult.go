package contest_participation

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
	AccessWithReason AccessWithReason                `json:"accessWithReason"`
	Info             db.UserContestParticipationInfo `json:"userContestParticipationInfo"`
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response = SaveSupertaskResultResponse{
			Success: false,
			Message: "Body reading failed",
		}
		writeSaveSupertaskResultResponse(w, response)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		response = SaveSupertaskResultResponse{
			Success: false,
			Message: "JSON decoding failed",
		}
		writeSaveSupertaskResultResponse(w, response)
		return
	}

	access, info, err := CheckContestAccess(userId, request.ContestId)
	if err != nil {
		response = SaveSupertaskResultResponse{
			Success: false,
			Message: err.Error(),
		}
		writeSaveSupertaskResultResponse(w, response)
		return
	}

	if info.UserContestRights.IsEmpty() {
		response = SaveSupertaskResultResponse{
			Success: false,
			Message: "User does not have rights on contest",
		}
		writeSaveSupertaskResultResponse(w, response)
		return
	}

	accessWithReason := ResolveAccessWithReason(access, info)

	data := SaveSupertaskResultData{
		AccessWithReason: accessWithReason,
		Info:             info,
	}

	if accessWithReason == AccessAllowed {
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
			writeSaveSupertaskResultResponse(w, response)
			return
		}
	}

	response = SaveSupertaskResultResponse{
		Success: true,
		Message: "OK",
		Data:    data,
	}
	writeSaveSupertaskResultResponse(w, response)
}

func writeSaveSupertaskResultResponse(w http.ResponseWriter, response SaveSupertaskResultResponse) {
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
