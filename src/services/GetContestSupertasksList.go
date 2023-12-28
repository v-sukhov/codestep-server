package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GetContestSupertasksListRequest struct {
	ContestId int32 `json:"contestId"`
}

type GetContestSupertasksListData struct {
	SupertasksList []db.SupertaskInContestInfo `json:"supertasksList"`
}

type GetContestSupertasksListResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    GetContestSupertasksListData `json:"data"`
}

func GetContestSupertasksList(w http.ResponseWriter, r *http.Request) {

	var request GetContestSupertasksListRequest
	var response GetContestSupertasksListResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetContestSupertasksListResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetContestSupertasksListResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {
			userContestRights, err := db.GetUserContestRights(userId, request.ContestId)

			if err != nil {
				response = GetContestSupertasksListResponse{
					Success: false,
					Message: "JSON decoding failed",
				}
			} else if userContestRights == 0 {
				response = GetContestSupertasksListResponse{
					Success: false,
					Message: "User does not have rights on contest",
				}
			} else {
				supertasksList, err := db.GetContestSupertasksList(request.ContestId)
				if err != nil {
					response = GetContestSupertasksListResponse{
						Success: false,
						Message: err.Error(),
					}
				} else {
					data := GetContestSupertasksListData{
						SupertasksList: supertasksList,
					}
					response = GetContestSupertasksListResponse{
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
