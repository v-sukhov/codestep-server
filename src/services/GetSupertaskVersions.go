package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GetSupertaskVersionsRequest struct {
	SupertaskId int32 `json:"supertaskId"`
}

type GetSupertaskVersionsData struct {
	SupertaskVersionList []db.SupertaskVersionShortInfo `json:"supertaskVersionList"`
}

type GetSupertaskVersionsResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    GetSupertaskVersionsData `json:"data"`
}

func GetSupertaskVersions(w http.ResponseWriter, r *http.Request) {
	var request GetSupertaskVersionsRequest
	var response GetSupertaskVersionsResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := io.ReadAll(r.Body); err != nil {
		response = GetSupertaskVersionsResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetSupertaskVersionsResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {
			supertaskUserRight, err := db.GetSupertaskUserRight(request.SupertaskId, userId)
			if err != nil {
				response = GetSupertaskVersionsResponse{
					Success: false,
					Message: err.Error(),
				}
			} else if supertaskUserRight == 0 {
				response = GetSupertaskVersionsResponse{
					Success: false,
					Message: "User does not have permission on the supertask",
				}
			} else {
				supertaskVersionList, err := db.GetSupertaskAllVersions(request.SupertaskId)
				if err != nil {
					response = GetSupertaskVersionsResponse{
						Success: false,
						Message: err.Error(),
					}
				} else {
					response = GetSupertaskVersionsResponse{
						Success: true,
						Message: "",
						Data: GetSupertaskVersionsData{
							SupertaskVersionList: supertaskVersionList,
						},
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
