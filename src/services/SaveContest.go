package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type SupertaskIdAndVersionNumber struct {
	SupertaskId   int32 `json:"supertaskId"`
	VersionNumber int32 `json:"versionNumber"`
}

type SaveContestRequest struct {
	ContestId            int32                         `json:"contestId"`
	ContestName          string                        `json:"contestName"`
	ContestDesc          string                        `json:"contestDesc"`
	ContestLogoHref      string                        `json:"contestLogoHref"`
	ContestStartTime     string                        `json:"contestStartTime"`
	GmtOffset            int32                         `json:"gmtOffset"`
	ContestDuration      int32                         `json:"contestDuration"`
	NoEndTime            bool                          `json:"noEndTime"`
	VirtualParticipation bool                          `json:"virtualParticipation"`
	LimitVirtualStart    bool                          `json:"limitVirtualStart"`
	VirtualStartEndTime  string                        `json:"virtualStartEndTime"`
	SupertaskList        []SupertaskIdAndVersionNumber `json:"supertaskList"`
}

type SaveContestData struct {
	ContestId int32 `json:"contestId"`
}

type SaveContestResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    SaveContestData `json:"data"`
}

// validateDateTimeFormat checks if the string is either empty or a valid datetime in ISO 8601 format
// Empty start time means that contest cannot be started
// Supports: YYYY-MM-DDThh:mm:ss[.sss][Z|(+|-)hh:mm]
func validateDateTimeFormat(dateTimeStr string) (*time.Time, error) {
	if dateTimeStr == "" {
		return nil, nil
	}

	// Try parsing with different ISO 8601 formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00",     // Full ISO with timezone
		"2006-01-02T15:04:05.000Z07:00", // With milliseconds and timezone
		"2006-01-02T15:04:05+07:00",     // With milliseconds and timezone
		"2006-01-02T15:04:05.000+07:00", // With milliseconds and timezone
		"2006-01-02T15:04:05.000Z",      // With milliseconds, UTC
		"2006-01-02T15:04:05Z",          // UTC timezone
		"2006-01-02T15:04:05.000",       // With milliseconds, no timezone
		"2006-01-02T15:04:05",           // Basic format
	}

	for _, format := range formats {
		if parsedTime, err := time.Parse(format, dateTimeStr); err == nil {
			return &parsedTime, nil
		}
	}

	return nil, fmt.Errorf("invalid datetime format: %s. Expected ISO 8601 format: YYYY-MM-DDThh:mm:ss[.sss][Z|(+|-)hh:mm]", dateTimeStr)
}

// validateContestData performs all validation checks for contest data
func validateContestData(request *SaveContestRequest) error {
	// 1. Validate contestStartTime format
	if _, err := validateDateTimeFormat(request.ContestStartTime); err != nil {
		return fmt.Errorf("contestStartTime: %v", err)
	}

	// 2. Validate gmtOffset in HOURS (-24..24)
	if request.GmtOffset < -24 || request.GmtOffset > 24 {
		return fmt.Errorf("gmtOffset must be between -24 and 24 hours, got: %d", request.GmtOffset)
	}

	// 3. If noEndTime == false, contestDuration must be > 0
	if !request.NoEndTime && request.ContestDuration <= 0 {
		return errors.New("contestDuration must be greater than 0 when noEndTime is false")
	}

	// 4. If virtualParticipation == true and limitVirtualStart == true, virtualStartEndTime must be valid
	if request.VirtualParticipation && request.LimitVirtualStart {
		if _, err := validateDateTimeFormat(request.VirtualStartEndTime); err != nil {
			return fmt.Errorf("virtualStartEndTime: %v", err)
		}
	}

	return nil
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
			// Validate contest data
			if err := validateContestData(&request); err != nil {
				response = SaveContestResponse{
					Success: false,
					Message: err.Error(),
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

					// Parse datetime strings to time.Time pointers
					var contestStartTime *time.Time
					if request.ContestStartTime != "" {
						if parsed, err := validateDateTimeFormat(request.ContestStartTime); err == nil {
							contestStartTime = parsed
						}
					}

					var virtualStartEndTime *time.Time
					if request.VirtualStartEndTime != "" {
						if parsed, err := validateDateTimeFormat(request.VirtualStartEndTime); err == nil {
							virtualStartEndTime = parsed
						}
					}

					contest := db.Contest{
						ContestId:            request.ContestId,
						ContestName:          request.ContestName,
						ContestDesc:          request.ContestDesc,
						ContestLogoHref:      request.ContestLogoHref,
						ContestStartTime:     contestStartTime,
						GmtOffset:            request.GmtOffset,
						ContestDuration:      request.ContestDuration,
						NoEndTime:            request.NoEndTime,
						VirtualParticipation: request.VirtualParticipation,
						LimitVirtualStart:    request.LimitVirtualStart,
						VirtualStartEndTime:  virtualStartEndTime,
						SupertaskList:        supertaskList,
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
	}

	if byteArr, err := json.Marshal(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Response marshal failed"))
	} else {
		w.Write(byteArr)
	}
}
