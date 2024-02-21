package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GetUserContestListRequest struct {
}

type GetUserContestListData struct {
	ContestList []db.ContestWithRights `json:"contestList"`
}

type GetUserContestListResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    GetUserContestListData `json:"data"`
}

func GetUserContestList(w http.ResponseWriter, r *http.Request) {

	var request GetUserContestListRequest
	var response GetUserContestListResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetUserContestListResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetUserContestListResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {

			contestList, err := db.GetUserContestList(userId)
			if err != nil {
				response = GetUserContestListResponse{
					Success: false,
					Message: err.Error(),
				}
			} else {
				data := GetUserContestListData{
					ContestList: contestList,
				}
				response = GetUserContestListResponse{
					Success: true,
					Message: "OK",
					Data:    data,
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
