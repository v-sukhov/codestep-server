package services

/*
	***********************************************************

	Получение информации по одному или нескольким supertask

	***********************************************************
*/

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"codestep/db"
	"codestep/security"
)

type GetSupertaskRequest struct {
	SupertaskId   int32 `json:"supertaskId"`
	VersionNumber int32 `json:"versionNumber"`
	AuthorUserId  int32 `json:"authorUserId"`
}

type GetSupertaskData struct {
	SupertaskId         int32     `json:"supertaskId"`
	VersionNumber       int32     `json:"versionNumber"`
	ParentVersionNumber int32     `json:"parentVersionNumber"`
	Commited            bool      `json:"commited"`
	AuthorUserId        int32     `json:"authorUserId"`
	SaveDTM             time.Time `json:"saveDTM"`
	CommitMessage       string    `json:"commitMessage"`
	SupertaskName       string    `json:"supertaskName"`
	SupertaskDesc       string    `json:"supertaskDesc"`
	SupertaskObjectJson string    `json:"supertaskObjectJson"`
}

type GetSupertaskResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    GetSupertaskData `json:"data"`
}

func GetSupertask(w http.ResponseWriter, r *http.Request) {
	var request GetSupertaskRequest
	var response GetSupertaskResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		response = GetSupertaskResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = GetSupertaskResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {
			supertaskUserRight, err := db.GetSupertaskUserRight(request.SupertaskId, userId)
			if err != nil {
				response = GetSupertaskResponse{
					Success: false,
					Message: err.Error(),
				}
			} else if supertaskUserRight == 0 {
				response = GetSupertaskResponse{
					Success: false,
					Message: "User does not have permition on the supertask",
				}
			} else {
				supertaskVersion, err := db.GetSupertaskVersion(request.SupertaskId, request.VersionNumber, request.AuthorUserId)

				if err != nil {
					response = GetSupertaskResponse{
						Success: false,
						Message: err.Error(),
					}
				} else {
					data := GetSupertaskData{
						SupertaskId:         supertaskVersion.SupertaskId,
						VersionNumber:       supertaskVersion.VersionNumber,
						ParentVersionNumber: supertaskVersion.ParentVersionNumber,
						AuthorUserId:        supertaskVersion.AuthorUserId,
						Commited:            supertaskVersion.Commited,
						CommitMessage:       supertaskVersion.CommitMessage,
						SaveDTM:             supertaskVersion.SaveDTM,
						SupertaskName:       supertaskVersion.SupertaskName,
						SupertaskDesc:       supertaskVersion.SupertaskDesc,
						SupertaskObjectJson: supertaskVersion.SupertaskObjectJson,
					}
					response = GetSupertaskResponse{
						Success: true,
						Message: "",
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
