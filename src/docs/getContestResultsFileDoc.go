package docs

import (
	contest "codestep/services/contest"
)

// swagger:route POST /api/protected/get-contest-results-file get-contest-results-file idOfGetContestResultsFileEndpoint
// Возвращает результаты участников контеста, отсортированные в порядке убывания баллов, решённых задач и возрастрания количества попыток в файле
// responses:
//   200: getContestResultsFileResponse

// swagger:parameters idOfGetContestResultsFileEndpoint
type getContestResultsFileRequestWrapper struct {
	// in:body
	Body contest.GetContestResultsFileRequest
}

// Возвращает результаты участников контеста, отсортированные в порядке убывания баллов, решённых задач и возрастрания количества попыток
// swagger:response getContestResultsFileResponse
type getContestResultsFileResponseWrapper struct {
	// in:body
}

