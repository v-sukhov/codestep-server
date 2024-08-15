package services

/*
	***********************************************************

	Получение информации по контесту
	Данные доступны только владельцам и администраторам контеста

	***********************************************************
*/

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"codestep/db"
	"codestep/security"
)

type GetContestRequest struct {
	ContestId int32 `json:"contestId"`
}

type GetContestResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    db.Contest `json:"data"`
}

func GetContest(w http.ResponseWriter, r *http.Request) {
	var request GetContestRequest
	var response GetContestResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetContestResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetContestResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {
			contestUserRights, err := db.GetContestUserRights(userId, request.ContestId)
			if err != nil {
				response = GetContestResponse{
					Success: false,
					Message: err.Error(),
				}
			} else if contestUserRights&3 == 0 { // 3 = 1 + 2 - owner and admin rights
				response = GetContestResponse{
					Success: false,
					Message: "User does not have permition on the contest",
				}
			} else {
				contest, err := db.GetContest(request.ContestId)

				if err != nil {
					response = GetContestResponse{
						Success: false,
						Message: err.Error(),
					}
				} else {
					response = GetContestResponse{
						Success: true,
						Message: "",
						Data:    contest,
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
