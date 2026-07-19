package contest_participation

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GetContestParticipationInfoRequest struct {
	ContestId  int32   `json:"contestId"`
	UserIdList []int32 `json:"userIdList"`
}

type ContestParticipationInfoByUser struct {
	UserId int32                           `json:"userId"`
	Info   db.UserContestParticipationInfo `json:"userContestParticipationInfo"`
}

type GetContestParticipationInfoData struct {
	Info           db.UserContestParticipationInfo  `json:"userContestParticipationInfo,omitempty"`
	InfoByUserList []ContestParticipationInfoByUser `json:"userContestParticipationInfoByUserList,omitempty"`
}

type GetContestParticipationInfoResponse struct {
	Success bool                            `json:"success"`
	Message string                          `json:"message"`
	Data    GetContestParticipationInfoData `json:"data"`
}

func GetContestParticipationInfo(w http.ResponseWriter, r *http.Request) {
	var request GetContestParticipationInfoRequest
	var response GetContestParticipationInfoResponse

	w.Header().Set("Content-Type", "application/json")

	callerUserId := r.Context().Value(security.ContextUserIdKey).(int32)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response = GetContestParticipationInfoResponse{
			Success: false,
			Message: "Body reading failed",
		}
		writeGetContestParticipationInfoResponse(w, response)
		return
	}

	if err := json.Unmarshal(body, &request); err != nil {
		response = GetContestParticipationInfoResponse{
			Success: false,
			Message: "JSON decoding failed",
		}
		writeGetContestParticipationInfoResponse(w, response)
		return
	}

	if len(request.UserIdList) == 0 {
		callerContestRights, err := db.GetContestUserRights(callerUserId, request.ContestId)
		if err != nil {
			response = GetContestParticipationInfoResponse{
				Success: false,
				Message: err.Error(),
			}
			writeGetContestParticipationInfoResponse(w, response)
			return
		}
		if callerContestRights == 0 {
			response = GetContestParticipationInfoResponse{
				Success: false,
				Message: "User does not have rights on contest",
			}
			writeGetContestParticipationInfoResponse(w, response)
			return
		}

		info, err := db.GetUserContestParticipationInfo(callerUserId, request.ContestId)
		if err != nil {
			response = GetContestParticipationInfoResponse{
				Success: false,
				Message: err.Error(),
			}
			writeGetContestParticipationInfoResponse(w, response)
			return
		}

		response = GetContestParticipationInfoResponse{
			Success: true,
			Message: "OK",
			Data: GetContestParticipationInfoData{
				Info: info,
			},
		}
		writeGetContestParticipationInfoResponse(w, response)
		return
	}

	callerContestRights, err := db.GetContestUserRights(callerUserId, request.ContestId)
	if err != nil {
		response = GetContestParticipationInfoResponse{
			Success: false,
			Message: err.Error(),
		}
		writeGetContestParticipationInfoResponse(w, response)
		return
	}
	if callerContestRights&3 == 0 {
		response = GetContestParticipationInfoResponse{
			Success: false,
			Message: "User does not have rights on operation",
		}
		writeGetContestParticipationInfoResponse(w, response)
		return
	}

	infoByUserList := make([]ContestParticipationInfoByUser, 0, len(request.UserIdList))
	for _, targetUserId := range request.UserIdList {
		info, err := db.GetUserContestParticipationInfo(targetUserId, request.ContestId)
		if err != nil {
			response = GetContestParticipationInfoResponse{
				Success: false,
				Message: err.Error(),
			}
			writeGetContestParticipationInfoResponse(w, response)
			return
		}
		infoByUserList = append(infoByUserList, ContestParticipationInfoByUser{
			UserId: targetUserId,
			Info:   info,
		})
	}

	response = GetContestParticipationInfoResponse{
		Success: true,
		Message: "OK",
		Data: GetContestParticipationInfoData{
			InfoByUserList: infoByUserList,
		},
	}
	writeGetContestParticipationInfoResponse(w, response)
}

func writeGetContestParticipationInfoResponse(w http.ResponseWriter, response GetContestParticipationInfoResponse) {
	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
