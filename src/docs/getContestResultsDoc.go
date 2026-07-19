package docs

import (
	contest "codestep/services/contest"
)

// swagger:route POST /api/protected/get-contest-results get-contest-results idOfGetContestResultsEndpoint
// Возвращает результаты участников контеста, отсортированные в порядке убывания баллов, решённых задач и возрастрания количества попыток
// responses:
//   200: getContestResultsResponse

// swagger:parameters idOfGetContestResultsEndpoint
type getContestResultsRequestWrapper struct {
	// in:body
	Body contest.GetContestResultsRequest
}

// Возвращает результаты участников контеста, отсортированные в порядке убывания баллов, решённых задач и возрастрания количества попыток
// swagger:response getContestResultsResponse
type getContestResultsResponseWrapper struct {
	// in:body
	Body contest.GetContestResultsResponse
}

