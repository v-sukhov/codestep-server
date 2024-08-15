package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type SupertaskIdAndVersionNumber struct {
	SupertaskId   int32 `json:"supertaskId"`
	VersionNumber int32 `json:"versionNumber"`
}

type SaveContestRequest struct {
	ContestId       int32                         `json:"contestId"`
	ContestName     string                        `json:"contestName"`
	ContestDesc     string                        `json:"contestDesc"`
	ContestLogoHref string                        `json:"contestLogoHref"`
	SupertaskList   []SupertaskIdAndVersionNumber `json:"supertaskList"`
}

type SaveContestData struct {
	ContestId int32 `json:"contestId"`
}

type SaveContestResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    SaveContestData `json:"data"`
}

/*
Перезаписывает данные контеста. Требует прав владельца или администратора контеста.
*/
func SaveContest(w http.ResponseWriter, r *http.Request) {

	var request SaveContestRequest
	var response SaveContestResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = SaveContestResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = SaveContestResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {
			userContestRights, err := db.GetContestUserRights(userId, request.ContestId)

			if err != nil {
				response = SaveContestResponse{
					Success: false,
					Message: err.Error(),
				}
			} else if userContestRights&3 == 0 {
				response = SaveContestResponse{
					Success: false,
					Message: "User does not have owner or admin rights on the contest",
				}
			} else {

				var supertaskList []db.SupertaskInContestInfo
				for _, row := range request.SupertaskList {
					supertaskList = append(supertaskList,
						db.SupertaskInContestInfo{
							SupertaskId:   row.SupertaskId,
							VersionNumber: row.VersionNumber,
						})
				}
				contest := db.Contest{
					ContestId:       request.ContestId,
					ContestName:     request.ContestName,
					ContestDesc:     request.ContestDesc,
					ContestLogoHref: request.ContestLogoHref,
					SupertaskList:   supertaskList,
				}

				if err := db.SaveContest(&contest, userId); err != nil {
					response = SaveContestResponse{
						Success: false,
						Message: err.Error(),
					}
				} else {
					data := SaveContestData{
						ContestId: contest.ContestId,
					}
					response = SaveContestResponse{
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
