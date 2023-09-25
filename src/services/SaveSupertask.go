package services

/*
	***************************************
		Создание и сохранение supertask

		Если SupertaskId = 0 - создаётся новый supertask

		Если Commited = true создаётся закоммиченная неизменяемая версия supertask.
		Ей присвивается новый номер версии. При этом указывается ParentVersionNumber - от какой версии произведена.
		Если ParentVersionNumber = 0 - то это считается первоначальной версией, то есть которая ни от чего не произведена.

		Если Commited = false создаётся незакоммиченная версия данного пользователя.
		Незакоммиченная версия у данного ползователя может быть только одна.

	***************************************
*/

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"codestep/db"
	"codestep/security"
)

type SaveSupertaskRequest struct {
	SupertaskId         int32  `json:"supertaskId"`
	ParentVersionNumber int32  `json:"parentVersionNumber"`
	Commited            bool   `json:"commited"`
	CommitMessage       string `json:"commitMessage"`
	SupertaskName       string `json:"supertaskName"`
	SupertaskDesc       string `json:"supertaskDesc"`
	SupertaskObjectJson string `json:"supertaskObjectJson"`
}

type SaveSupertaskData struct {
	SupertaskId   int32 `json:"supertaskId"`
	VersionNumber int32 `json:"versionNumber"`
}

type SaveSupertaskResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    SaveSupertaskData `json:"data"`
}

func SaveSupertask(w http.ResponseWriter, r *http.Request) {

	var request SaveSupertaskRequest
	var response SaveSupertaskResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		response = SaveSupertaskResponse{
			Success: false,
			Message: "Body reading failed",
		}
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			response = SaveSupertaskResponse{
				Success: false,
				Message: "JSON decoding failed",
			}
		} else {

			success := true
			message := ""

			userRights, err := db.GetUserRights(userId)

			if err != nil {
				success = false
				message = err.Error()
			} else if !userRights.IsDeveloper {
				success = false
				message = "User does not have developer rights to create or edit supertask"
			}

			if success && request.SupertaskId != 0 {
				// Если сохраняется существующая задача - необходимо проверить права пользователя на эту задачу
				supertaskUserRight, err := db.GetSupertaskUserRight(request.SupertaskId, userId)

				if err != nil {
					success = false
					message = err.Error()
				} else if supertaskUserRight != 1 && supertaskUserRight != 2 {
					success = false
					message = "User does not have permition on the supertask"
				}
			}

			if success {
				supertaskVersion := db.SupertaskVersion{
					SupertaskId:         request.SupertaskId,
					VersionNumber:       0,
					AuthorUserId:        userId,
					ParentVersionNumber: request.ParentVersionNumber,
					Commited:            request.Commited,
					CommitMessage:       request.CommitMessage,
					SupertaskName:       request.SupertaskName,
					SupertaskDesc:       request.SupertaskDesc,
					SupertaskObjectJson: request.SupertaskObjectJson,
				}

				if err := db.SaveSupertaskVersion(&supertaskVersion); err != nil {
					response = SaveSupertaskResponse{
						Success: false,
						Message: err.Error(),
					}
				} else {
					data := SaveSupertaskData{
						SupertaskId:   supertaskVersion.SupertaskId,
						VersionNumber: supertaskVersion.VersionNumber,
					}
					response = SaveSupertaskResponse{
						Success: true,
						Message: "OK",
						Data:    data,
					}
				}
			} else {
				response = SaveSupertaskResponse{
					Success: false,
					Message: message,
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
