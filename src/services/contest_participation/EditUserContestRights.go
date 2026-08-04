package contest_participation

/*
	Изменение прав пользователей на контест.
	Права ADMIN, JURY, PARTICIPANT задаются полной заменой (как ManageUserContestRights).
	OWNER нельзя назначить или отозвать — он сохраняется у пользователя независимо от переданного списка ролей.
	Доступно только владельцу и администратору контеста.
*/

import (	"codestep/db"
	"codestep/security"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type EditUserContestRightsItem struct {
	UserId int32    `json:"userId"`
	Roles  []string `json:"roles"`
}

type EditUserContestRightsRequest struct {
	ContestId int32                       `json:"contestId"`
	Rows      []EditUserContestRightsItem `json:"rows"`
}

type EditUserContestRightsRowError struct {
	RowNumber    int    `json:"rowNumber"`
	ErrorMessage string `json:"errorMessage"`
}

type EditUserContestRightsData struct {
	RowError []EditUserContestRightsRowError `json:"rowError"`
}

type EditUserContestRightsResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    EditUserContestRightsData `json:"data"`
}

func parseEditUserContestRoles(roles []string) ([]int16, string) {
	rightIds := make([]int16, 0, len(roles))
	seen := make(map[int16]bool)

	for _, role := range roles {
		var rightId int16
		switch strings.ToUpper(strings.TrimSpace(role)) {
		case "OWNER":
			return nil, "Role OWNER cannot be assigned or revoked"
		case "ADMIN":
			rightId = 2
		case "JURY":
			rightId = 3
		case "PARTICIPANT":
			rightId = 4
		default:
			return nil, "Unknown role: " + role
		}
		if !seen[rightId] {
			seen[rightId] = true
			rightIds = append(rightIds, rightId)
		}
	}

	if len(rightIds) == 0 {
		return nil, "At least one role must be specified"
	}

	return rightIds, ""
}

func EditUserContestRights(w http.ResponseWriter, r *http.Request) {
	var request EditUserContestRightsRequest
	var response EditUserContestRightsResponse

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response = EditUserContestRightsResponse{
			Success: false,
			Message: "Body reading failed",
		}
		writeEditUserContestRightsResponse(w, response)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		response = EditUserContestRightsResponse{
			Success: false,
			Message: "JSON decoding failed",
		}
		writeEditUserContestRightsResponse(w, response)
		return
	}

	response = processEditUserContestRights(userId, request.ContestId, request.Rows)
	writeEditUserContestRightsResponse(w, response)
}

func processEditUserContestRights(
	requestUserId int32,
	contestId int32,
	items []EditUserContestRightsItem,
) EditUserContestRightsResponse {
	if contestId <= 0 {
		return EditUserContestRightsResponse{
			Success: false,
			Message: "Contest Id is required",
		}
	}

	userContestRights, err := db.GetContestUserRights(requestUserId, contestId)
	if err != nil {
		return EditUserContestRightsResponse{
			Success: false,
			Message: err.Error(),
		}
	}
	if userContestRights&3 == 0 {
		return EditUserContestRightsResponse{
			Success: false,
			Message: "User does not have rights on operation",
		}
	}

	rightsRows := make([]db.ManageUserContestRightsRow, len(items))
	userIdsToLookup := make([]int32, 0, len(items))

	for i, item := range items {
		if item.UserId <= 0 {
			rightsRows[i].IsError = true
			rightsRows[i].ErrorMessage = "User Id is required"
			continue
		}

		rightIds, roleErr := parseEditUserContestRoles(item.Roles)
		if roleErr != "" {
			rightsRows[i].IsError = true
			rightsRows[i].ErrorMessage = roleErr
			continue
		}

		userIdsToLookup = append(userIdsToLookup, item.UserId)
		rightsRows[i] = db.ManageUserContestRightsRow{
			Login:     "",
			ContestId: contestId,
			Rights:    rightIds,
		}
	}

	if len(userIdsToLookup) > 0 {
		loginsByUserId, err := db.GetUserLoginsByIds(userIdsToLookup)
		if err != nil {
			return EditUserContestRightsResponse{
				Success: false,
				Message: err.Error(),
			}
		}

		for i, item := range items {
			if rightsRows[i].IsError {
				continue
			}

			login, ok := loginsByUserId[item.UserId]
			if !ok {
				rightsRows[i].IsError = true
				rightsRows[i].ErrorMessage = fmt.Sprintf("Unknown user Id: %d", item.UserId)
				continue
			}

			rightsRows[i].Login = login
		}
	}

	hasValidRows := false
	for _, row := range rightsRows {
		if !row.IsError {
			hasValidRows = true
			break
		}
	}

	if hasValidRows {
		if err := db.ManageUserContestRights(rightsRows); err != nil {
			return EditUserContestRightsResponse{
				Success: false,
				Message: err.Error(),
			}
		}
	}

	rowError := make([]EditUserContestRightsRowError, 0)
	for i, row := range rightsRows {
		if row.IsError {
			rowError = append(rowError, EditUserContestRightsRowError{
				RowNumber:    i + 1,
				ErrorMessage: row.ErrorMessage,
			})
		}
	}

	return EditUserContestRightsResponse{
		Success: true,
		Message: "",
		Data: EditUserContestRightsData{
			RowError: rowError,
		},
	}
}

func writeEditUserContestRightsResponse(w http.ResponseWriter, response EditUserContestRightsResponse) {
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
