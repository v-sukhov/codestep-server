package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GetContestResultsRequest struct {
	ContestId int32 `json:"contestId"`
}

type GetContestResultsResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    db.ContestResults `json:"data"`
}

func GetContestResults(w http.ResponseWriter, r *http.Request) {
	var request GetContestResultsRequest
	var response GetContestResultsResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetContestResultsResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetContestResultsResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {
			userRights, err := db.GetUserRights(userId)

			if err != nil {
				response = GetContestResultsResponse{
					Success: false,
					Message: err.Error(),
				}
			} else if !userRights.IsAdmin {
				response = GetContestResultsResponse{
					Success: false,
					Message: "User does not have admin rights to request contest result",
				}
			} else {
				results, err := db.GetContestResults(request.ContestId)

				if err != nil {
					response = GetContestResultsResponse{
						Success: false,
						Message: err.Error(),
					}
				} else {
					response = GetContestResultsResponse{
						Success: true,
						Data:    results,
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
