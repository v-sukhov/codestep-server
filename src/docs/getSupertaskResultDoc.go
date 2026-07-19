package docs

import (
	supertask "codestep/services/supertask"
)

// swagger:route POST /api/protected/get-supertask-result get-supertask-result idOfGetSupertaskResultEndpoint
// Получить результат по задаче
// responses:
//   200: getSupertaskResultResponse

// swagger:parameters idOfGetSupertaskResultEndpoint
type getSupertaskResultRequestWrapper struct {
	// in:body
	Body supertask.GetSupertaskResultRequest
}

// Get supertask result success
// swagger:response getSupertaskResultResponse
type getSupertaskResultWrapper struct {
	// in:body
	Body supertask.GetSupertaskResultResponse
}

