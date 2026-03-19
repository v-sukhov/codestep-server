package services

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

type GetSupertaskInContestWithResultsResponse struct {
	Success bool                             `json:"success"`
	Message string                           `json:"message"`
	Data    db.SupertaskInContestWithResults `json:"data"`
}

func GetSupertaskInContestWithResults(w http.ResponseWriter, r *http.Request) {
	var request GetSupertaskInContestWithResultsRequest
	var response GetSupertaskInContestWithResultsResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetSupertaskInContestWithResultsResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetSupertaskInContestWithResultsResponse{
				Success: false,
				Message: "JSON decoding failed ",
			}
		} else {
			// Проверка прав пользователя на контест
			userContestRights, err := db.GetContestUserRights(userId, request.ContestId)
			if err != nil {
				response = GetSupertaskInContestWithResultsResponse{
					Success: false,
					Message: err.Error(),
				}
			} else if userContestRights == 0 {
				response = GetSupertaskInContestWithResultsResponse{
					Success: false,
					Message: "User does not have rights on contest",
				}
			} else {
				// Проверка доступа к контесту по времени
				access, _, err := CheckContestTimeAccess(userId, request.ContestId)
				if err != nil {
					response = GetSupertaskInContestWithResultsResponse{
						Success: false,
						Message: err.Error(),
					}
				} else if !access {
					response = GetSupertaskInContestWithResultsResponse{
						Success: false,
						Message: "Access to contest is not allowed at this time",
					}
				} else {
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
					} else {
						response = GetSupertaskInContestWithResultsResponse{
							Success: true,
							Message: "",
							Data:    supertaskInContestWithResults,
						}
					}
				}
			}
		}
	}

	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
