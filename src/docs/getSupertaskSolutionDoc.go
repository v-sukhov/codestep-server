package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-supertask-solution get-supertask-solution idOfGetSupertaskSolutionEndpoint
// Получить решение задачи
// responses:
//   200: getSupertaskSolutionResponse

// swagger:parameters idOfGetSupertaskSolutionEndpoint
type getSupertaskSolutionRequestWrapper struct {
	// in:body
	Body services.GetSupertaskSolutionRequest
}

// Get supertask solution success
// swagger:response getSupertaskSolutionResponse
type getSupertaskSolutionWrapper struct {
	// in:body
	Body services.GetSupertaskSolutionResponse
}
