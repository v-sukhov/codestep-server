package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type RemoveSupertaskFromContestRequest struct {
	ContestId   int32 `json:"contestId"`
	SupertaskId int32 `json:"supertaskId"`
}

type RemoveSupertaskFromContestData struct {
}

type RemoveSupertaskFromContestResponse struct {
	Success bool                           `json:"success"`
	Message string                         `json:"message"`
	Data    RemoveSupertaskFromContestData `json:"data"`
}

func RemoveSupertaskFromContest(w http.ResponseWriter, r *http.Request) {

	var request RemoveSupertaskFromContestRequest
	var response RemoveSupertaskFromContestResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = RemoveSupertaskFromContestResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = RemoveSupertaskFromContestResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {

			userContestRights, err := db.GetUserContestRights(userId, request.ContestId)

			if err != nil {
				response = RemoveSupertaskFromContestResponse{
					Success: false,
					Message: err.Error(),
				}
			} else if userContestRights&3 == 0 {
				// 3 = 00..011 = 1 + 2 - owner and admin
				response = RemoveSupertaskFromContestResponse{
					Success: false,
					Message: "User does not have admin right on contest",
				}
			} else {
				err := db.RemoveSupertaskFromContest(request.ContestId, request.SupertaskId)

				if err != nil {
					response = RemoveSupertaskFromContestResponse{
						Success: false,
						Message: err.Error(),
					}
				} else {
					data := RemoveSupertaskFromContestData{}
					response = RemoveSupertaskFromContestResponse{
						Success: true,
						Message: "OK",
						Data:    data,
					}
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
