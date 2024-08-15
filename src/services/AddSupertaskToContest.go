package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type AddSupertaskToContestRequest struct {
	ContestId              int32 `json:"contestId"`
	SupertaskId            int32 `json:"supertaskId"`
	SupertaskVersionNumber int32 `json:"supertaskVersionNumber"`
}

type AddSupertaskToContestData struct {
	OrderNumber int16 `json:"orderNumber"`
}

type AddSupertaskToContestResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    AddSupertaskToContestData `json:"data"`
}

func AddSupertaskToContest(w http.ResponseWriter, r *http.Request) {

	var request AddSupertaskToContestRequest
	var response AddSupertaskToContestResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = AddSupertaskToContestResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = AddSupertaskToContestResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {

			userContestRights, err := db.GetContestUserRights(userId, request.ContestId)

			if err != nil {
				response = AddSupertaskToContestResponse{
					Success: false,
					Message: err.Error(),
				}
			} else if userContestRights&3 == 0 {
				// 3 = 00..011 = 1 + 2 - owner and admin
				response = AddSupertaskToContestResponse{
					Success: false,
					Message: "User does not have admin right on contest",
				}
			} else {
				orderNumber, err := db.AddSupertaskToContest(request.ContestId, request.SupertaskId, request.SupertaskVersionNumber)

				if err != nil {
					response = AddSupertaskToContestResponse{
						Success: false,
						Message: err.Error(),
					}
				} else {
					data := AddSupertaskToContestData{
						OrderNumber: orderNumber,
					}
					response = AddSupertaskToContestResponse{
						Success: true,
						Message: "OK",
						Data:    data,
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
