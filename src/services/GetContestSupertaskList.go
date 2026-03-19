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
	*******************************************
	Получение списка задач в контесте.
	Доступ к данным для участников и жюри ограничивается временем старта контеста:
		- для жюри доступ только после старта контеста
		- для участников - только во время участия, вызов сервиса вызывает старт участия, если он возможен
	Доступ для владельца и администраторов не ограничен по времени.
	*******************************************
*/

type GetContestSupertaskListRequest struct {
	ContestId int32 `json:"contestId"`
}

type UserContestRights struct {
	IsOwner       bool `json:"isOwner"`
	IsAdmin       bool `json:"isAdmin"`
	IsJury        bool `json:"isJury"`
	IsParticipant bool `json:"isParticipant"`
}

type GetContestSupertaskListData struct {
	Success           bool                            `json:"success"`
	SupertasksList    []db.SupertaskInContestInfo     `json:"supertaskList"`
	UserContestRights UserContestRights               `json:"userContestRights"`
	Info              db.UserContestParticipationInfo `json:"userContestParticipationInfo"`
}

type GetContestSupertaskListResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    GetContestSupertaskListData `json:"data"`
}

func writeResponse(w http.ResponseWriter, response GetContestSupertaskListResponse) {
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}

func GetContestSupertaskList(w http.ResponseWriter, r *http.Request) {

	var request GetContestSupertaskListRequest
	var response GetContestSupertaskListResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response = GetContestSupertaskListResponse{
			Success: false,
			Message: "Body reading failed",
		}
		writeResponse(w, response)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		response = GetContestSupertaskListResponse{
			Success: false,
			Message: "JSON decoding failed",
		}
		writeResponse(w, response)
		return
	}

	userContestRights, err := db.GetContestUserRights(userId, request.ContestId)
	if err != nil {
		response = GetContestSupertaskListResponse{
			Success: false,
			Message: err.Error(),
		}
		writeResponse(w, response)
		return
	}

	if userContestRights == 0 {
		response = GetContestSupertaskListResponse{
			Success: false,
			Message: "User does not have rights on contest",
		}
		writeResponse(w, response)
		return
	}

	// Проверка доступа к контесту по времени
	access, info, err := CheckContestTimeAccess(userId, request.ContestId)
	if err != nil {
		response = GetContestSupertaskListResponse{
			Success: false,
			Message: err.Error(),
		}
		writeResponse(w, response)
		return
	}

	var supertaskList []db.SupertaskInContestInfo

	if access {
		if supertaskList, err = db.GetContestSupertaskList(request.ContestId); err != nil {
			response = GetContestSupertaskListResponse{
				Success: false,
				Message: err.Error(),
			}
			writeResponse(w, response)
			return
		}
	}

	data := GetContestSupertaskListData{
		Success:        access,
		SupertasksList: supertaskList,
		UserContestRights: UserContestRights{
			IsOwner:       (userContestRights&1 > 0),
			IsAdmin:       (userContestRights&2 > 0),
			IsJury:        (userContestRights&4 > 0),
			IsParticipant: (userContestRights&8 > 0),
		},
		Info: info,
	}
	response = GetContestSupertaskListResponse{
		Success: true,
		Message: "OK",
		Data:    data,
	}

	writeResponse(w, response)
}

/** TO DELETE

func GetContestSupertaskList_OLD(w http.ResponseWriter, r *http.Request) {

	var request GetContestSupertaskListRequest
	var response GetContestSupertaskListResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetContestSupertaskListResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetContestSupertaskListResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {
			userContestRights, err := db.GetContestUserRights(userId, request.ContestId)

			if err != nil {
				response = GetContestSupertaskListResponse{
					Success: false,
					Message: err.Error(),
				}
			} else if userContestRights == 0 {
				response = GetContestSupertaskListResponse{
					Success: false,
					Message: "User does not have rights on contest",
				}
			} else {
				// TODO: здесь нужно добавить проверку возможности доступа к контесту по времени:
				// 1. Данная проверка должна быть реализована только для пользователей, не обладающих правами владельца или администратора контеста (1 и 2)
				// 2. Логика проверки для пользователей с правами "участник" следующая:
				// - если соревнование ещё не началось - нужно вернуть ответ о времени старта соревнования
				// - если началось - нужно проверить, записан ли пользователь в таблицу USER_CONTEST_START_TIME,
				// 			если нет - записать и вернуть список суперзадач, если есть - вернуть список суперзадач
				// 3. Логика проверки для пользователей с правами "жюри" такая же, но пользователь не отмечается в таблице USER_CONTEST_START_TIME
				// между отсчётным стартом контеста и отсчётным стартом + длительность контеста

				access := false

				if userContestRights&3 > 0 { // 3 = 2 + 1 - права владельца и админа
					access = true
				} else {
					if info, err := db.GetUserContestParticipationStatus(userId, request.ContestId); err != nil {
						response = GetContestSupertaskListResponse{
							Success: false,
							Message: err.Error(),
						}
						access = false
					} else {

					}
				}

				if access && err == nil {
					supertaskList, err := db.GetContestSupertaskList(request.ContestId)
					if err != nil {
						response = GetContestSupertaskListResponse{
							Success: false,
							Message: err.Error(),
						}
					} else {
						data := GetContestSupertaskListData{
							SupertasksList: supertaskList,
							UserContestRights: UserContestRights{
								IsOwner:       (userContestRights&1 > 0),
								IsAdmin:       (userContestRights&2 > 0),
								IsJury:        (userContestRights&4 > 0),
								IsParticipant: (userContestRights&8 > 0),
							},
						}
						response = GetContestSupertaskListResponse{
							Success: true,
							Message: "OK",
							Data:    data,
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
*/
