package contest_participation

/**
Запрос содержания задачи вместе с результатами запрашивающего пользователя
Сервис предназначен для вызова участниками владельцем, администратором, жюри и участниками контеста
Доступ для владельца и администратора не ограничен по времени
Доступ для участника и жюри ограничен по времени
*/

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GetSupertaskInContestWithResultsRequest struct {
	ContestId   int32 `json:"contestId"`
	SupertaskId int32 `json:"supertaskId"`
}

type GetSupertaskInContestWithResultsData struct {
	AccessWithReason              AccessWithReason                 `json:"accessWithReason"`
	Info                          db.UserContestParticipationInfo  `json:"userContestParticipationInfo"`
	SupertaskInContestWithResults db.SupertaskInContestWithResults `json:"supertaskInContestWithResults"`
}

type GetSupertaskInContestWithResultsResponse struct {
	Success bool                                 `json:"success"`
	Message string                               `json:"message"`
	Data    GetSupertaskInContestWithResultsData `json:"data"`
}

func GetSupertaskInContestWithResults(w http.ResponseWriter, r *http.Request) {
	var request GetSupertaskInContestWithResultsRequest
	var response GetSupertaskInContestWithResultsResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response = GetSupertaskInContestWithResultsResponse{
			Success: false,
			Message: "Body reading failed",
		}
		writeGetSupertaskInContestWithResultsResponse(w, response)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		response = GetSupertaskInContestWithResultsResponse{
			Success: false,
			Message: "JSON decoding failed",
		}
		writeGetSupertaskInContestWithResultsResponse(w, response)
		return
	}

	access, info, err := CheckContestAccess(userId, request.ContestId)
	if err != nil {
		response = GetSupertaskInContestWithResultsResponse{
			Success: false,
			Message: err.Error(),
		}
		writeGetSupertaskInContestWithResultsResponse(w, response)
		return
	}

	if info.UserContestRights.IsEmpty() {
		response = GetSupertaskInContestWithResultsResponse{
			Success: false,
			Message: "User does not have rights on contest",
		}
		writeGetSupertaskInContestWithResultsResponse(w, response)
		return
	}

	accessWithReason := ResolveAccessWithReason(access, info)

	data := GetSupertaskInContestWithResultsData{
		AccessWithReason: accessWithReason,
		Info:             info,
	}

	if accessWithReason == AccessAllowed {
		supertaskInContestWithResults, err := db.GetSupertaskInContestWithResults(
			request.ContestId,
			request.SupertaskId,
			userId,
		)
		if err != nil {
			response = GetSupertaskInContestWithResultsResponse{
				Success: false,
				Message: err.Error(),
			}
			writeGetSupertaskInContestWithResultsResponse(w, response)
			return
		}
		data.SupertaskInContestWithResults = supertaskInContestWithResults
	}

	response = GetSupertaskInContestWithResultsResponse{
		Success: true,
		Message: "OK",
		Data:    data,
	}
	writeGetSupertaskInContestWithResultsResponse(w, response)
}

func writeGetSupertaskInContestWithResultsResponse(w http.ResponseWriter, response GetSupertaskInContestWithResultsResponse) {
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
