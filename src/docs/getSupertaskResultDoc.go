package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-supertask-result get-supertask-result idOfGetSupertaskResultEndpoint
// Получить результат по задаче
// responses:
//   200: getSupertaskResultResponse

// swagger:parameters idOfGetSupertaskResultEndpoint
type getSupertaskResultRequestWrapper struct {
	// in:body
	Body services.GetSupertaskResultRequest
}

// Get supertask result success
// swagger:response getSupertaskResultResponse
type getSupertaskResultWrapper struct {
	// in:body
	Body services.GetSupertaskResultResponse
}
