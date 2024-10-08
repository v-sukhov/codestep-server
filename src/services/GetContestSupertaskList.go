package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

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
	SupertasksList    []db.SupertaskInContestInfo `json:"supertaskList"`
	UserContestRights UserContestRights           `json:"userContestRights"`
}

type GetContestSupertaskListResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    GetContestSupertaskListData `json:"data"`
}

func GetContestSupertaskList(w http.ResponseWriter, r *http.Request) {

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
					Message: "JSON decoding failed",
				}
			} else if userContestRights == 0 {
				response = GetContestSupertaskListResponse{
					Success: false,
					Message: "User does not have rights on contest",
				}
			} else {
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

	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
