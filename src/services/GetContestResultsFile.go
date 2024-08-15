package services

import (
	"codestep/db"
	"codestep/security"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetContestResultsFileRequest struct {
	ContestId int32 `json:"contestId"`
}

func GetContestResultsFile(w http.ResponseWriter, r *http.Request) {
	var request GetContestResultsRequest

	w.Header().Set("Content-Type", "application/json")

	userId := r.Context().Value(security.ContextUserIdKey).(int32)

	success := true
	message := ""
	var results db.ContestResults

	if body, err := io.ReadAll(r.Body); err != nil {
		success = false
		message = "Body reading failed"
	} else {
		if err := json.Unmarshal(body, &request); err != nil {
			success = false
			message = "JSON decoding failed"
		} else {
			userContestRights, err := db.GetContestUserRights(userId, request.ContestId)

			if err != nil {
				success = false
				message = err.Error()
			} else if userContestRights&7 == 0 { // 7 = 4 + 2 + 1 - права владельца, администратора или жюри
				success = false
				message = "User does not have owner, admin or jury rights to request contest result"
			} else {
				results, err = db.GetContestResults(request.ContestId)

				if err != nil {
					success = false
					message = err.Error()
				} else {
					success = true
				}
			}
		}
	}

	answerFilename := fmt.Sprintf("contest-results-%d.txt", request.ContestId)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+answerFilename+"\"")

	var bytes []byte

	if !success {
		bytes = []byte(message)
	} else {
		bytes = printContestResultsToCSV(results)
	}

	w.Write(bytes)

}

func appendContestResultRowToCSV(bytes []byte, result db.ContestUserResult) []byte {
	var res []byte

	res = fmt.Appendf(bytes, "%s", result.UserLogin)

	for i, st := range result.SupertaskScore {
		for j, t := range st {
			passedSign := ""
			if result.SupertaskPassed[i][j] {
				passedSign = "+"
			}
			res = fmt.Appendf(res, "\t%s%d(%d)", passedSign, t, result.SupertaskTries[i][j])
		}
	}

	res = fmt.Appendf(res, "\t%d\t%d\t%d\n", result.TotalScore, result.TotalPassed, result.TotalTries)

	return res
}

func printContestResultsToCSV(results db.ContestResults) []byte {
	var bytes []byte

	for _, e := range results.Errors {
		bytes = fmt.Appendf(bytes, "%s\n", e)
	}

	for _, name := range results.SupertaskNames {
		bytes = fmt.Appendf(bytes, "\t\"%s\"", name)
	}

	bytes = fmt.Appendf(bytes, "\n")

	bytes = fmt.Appendf(bytes, "login")

	for _, tasks := range results.MaxPossibleResult.SupertaskPassed {
		for i := range tasks {
			bytes = fmt.Appendf(bytes, "\t%d", i+1)
		}
	}

	bytes = fmt.Appendf(bytes, "\ttotal_score\ttotal_passed\ttotal_tries\n")

	bytes = appendContestResultRowToCSV(bytes, results.MaxPossibleResult)

	for _, ur := range results.UserResults {
		bytes = appendContestResultRowToCSV(bytes, ur)
	}

	return bytes
}
