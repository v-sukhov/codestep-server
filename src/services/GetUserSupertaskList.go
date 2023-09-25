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

type GetUserSupertaskListData struct {
	Supertasks []db.SupertaskLastVersionInfo `json:"supertasks"`
}

type GetUserSupertaskListResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    GetUserSupertaskListData `json:"data"`
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
		data := GetUserSupertaskListData{
			Supertasks: supertasks,
		}
		response = GetUserSupertaskListResponse{
			Success: true,
			Message: "",
			Data:    data,
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
