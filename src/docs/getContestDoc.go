package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-contest get-contest idOfGetContestEndpoint
// Возвращает данные контеста
// responses:
//   200: getContestResponse

// swagger:parameters idOfGetContestEndpoint
type getContestRequestWrapper struct {
	// in:body
	Body services.GetContestRequest
}

// Возвращает данные контеста
// swagger:response getContestResponse
type getContestResponseWrapper struct {
	// in:body
	Body services.GetContestResponse
}
