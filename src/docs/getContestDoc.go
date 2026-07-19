package docs

import (
	contest "codestep/services/contest"
)

// swagger:route POST /api/protected/get-contest get-contest idOfGetContestEndpoint
// Возвращает данные контеста
// responses:
//   200: getContestResponse

// swagger:parameters idOfGetContestEndpoint
type getContestRequestWrapper struct {
	// in:body
	Body contest.GetContestRequest
}

// Возвращает данные контеста
// swagger:response getContestResponse
type getContestResponseWrapper struct {
	// in:body
	Body contest.GetContestResponse
}

