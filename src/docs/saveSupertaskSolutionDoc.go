package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/save-supertask-solution save-supertask-solution idOfSaveSupertaskSolutionEndpoint
// Сохранить решение задачи
// responses:
//   200: saveSupertaskSolutionResponse

// swagger:parameters idOfSaveSupertaskSolutionEndpoint
type saveSupertaskSolutionRequestWrapper struct {
	// in:body
	Body services.SaveSupertaskSolutionRequest
}

// Save supertask solution success
// swagger:response saveSupertaskSolutionResponse
type saveSupertaskSolutionWrapper struct {
	// in:body
	Body services.SaveSupertaskSolutionResponse
}
