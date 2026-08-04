package contest

/*
	Получение списка пользователей контеста с пагинацией, сортировкой и фильтрацией.
	Доступно только владельцам и администраторам контеста.
*/

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"codestep/db"
	"codestep/security"
)

type GetContestUserListFilter struct {
	LoginSubstring string   `json:"loginSubstring"`
	Roles          []string `json:"roles"`
}

type GetContestUserListRequest struct {
	ContestId     int32                    `json:"contestId"`
	PageSize      int32                    `json:"pageSize"`
	PageNumber    int32                    `json:"pageNumber"`
	SortField     string                   `json:"sortField"`
	SortDirection string                   `json:"sortDirection"`
	Filter        GetContestUserListFilter `json:"filter"`
}

type GetContestUserListResponse struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    db.GetContestUserListResult `json:"data"`
}

func GetContestUserList(w http.ResponseWriter, r *http.Request) {
	var request GetContestUserListRequest
	var response GetContestUserListResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response = GetContestUserListResponse{
			Success: false,
			Message: "Body reading failed",
		}
		writeGetContestUserListResponse(w, response)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		response = GetContestUserListResponse{
			Success: false,
			Message: "JSON decoding failed",
		}
		writeGetContestUserListResponse(w, response)
		return
	}

	if request.PageSize <= 0 {
		response = GetContestUserListResponse{
			Success: false,
			Message: "Page size must be greater than zero",
		}
		writeGetContestUserListResponse(w, response)
		return
	}

	if request.PageNumber <= 0 {
		response = GetContestUserListResponse{
			Success: false,
			Message: "Page number must be greater than zero",
		}
		writeGetContestUserListResponse(w, response)
		return
	}

	userContestRights, err := db.GetContestUserRights(userId, request.ContestId)
	if err != nil {
		response = GetContestUserListResponse{
			Success: false,
			Message: err.Error(),
		}
		writeGetContestUserListResponse(w, response)
		return
	}
	if userContestRights&3 == 0 {
		response = GetContestUserListResponse{
			Success: false,
			Message: "User does not have permition on the contest",
		}
		writeGetContestUserListResponse(w, response)
		return
	}

	sortField := strings.TrimSpace(request.SortField)
	if sortField == "" {
		sortField = "login"
	}

	sortDirection := strings.ToLower(strings.TrimSpace(request.SortDirection))
	if sortDirection == "" {
		sortDirection = "asc"
	} else if sortDirection != "asc" && sortDirection != "desc" {
		response = GetContestUserListResponse{
			Success: false,
			Message: "Sort direction must be asc or desc",
		}
		writeGetContestUserListResponse(w, response)
		return
	}

	data, err := db.GetContestUserList(db.GetContestUserListParams{
		ContestId:      request.ContestId,
		PageSize:       request.PageSize,
		PageNumber:     request.PageNumber,
		SortField:      sortField,
		SortDirection:  sortDirection,
		LoginSubstring: strings.TrimSpace(request.Filter.LoginSubstring),
		Roles:          request.Filter.Roles,
	})
	if err != nil {
		response = GetContestUserListResponse{
			Success: false,
			Message: err.Error(),
		}
		writeGetContestUserListResponse(w, response)
		return
	}

	response = GetContestUserListResponse{
		Success: true,
		Message: "OK",
		Data:    data,
	}
	writeGetContestUserListResponse(w, response)
}

func writeGetContestUserListResponse(w http.ResponseWriter, response GetContestUserListResponse) {
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
