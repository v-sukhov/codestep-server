package services

/*
	***********************************************************

	Получение списка суперзадач, доступных пользователю

	***********************************************************
*/

import (
	"encoding/json"
	"log"
	"net/http"

	"codestep/db"
	"codestep/security"
)

type GetUserSupertaskListRequest struct {
}

type GetUserSupertaskListResponse struct {
	Success    bool                          `json:"success"`
	Message    string                        `json:"message"`
	Supertasks []db.SupertaskLastVersionInfo `json:"supertasks"`
}

func GetUserSupertaskList(w http.ResponseWriter, r *http.Request) {
	var response GetUserSupertaskListResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	supertasks, err := db.GetUserSupertaskList(userId)

	if err != nil {
		response = GetUserSupertaskListResponse{
			Success: false,
			Message: err.Error(),
		}
	} else {
		response = GetUserSupertaskListResponse{
			Success:    true,
			Message:    "",
			Supertasks: supertasks,
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
