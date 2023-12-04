package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/save-supertask-result save-supertask-result idOfSaveSupertaskResultEndpoint
// Сохранить результат тестирования задачи
// responses:
//   200: saveSupertaskResultResponse

// swagger:parameters idOfSaveSupertaskResultEndpoint
type saveSupertaskResultRequestWrapper struct {
	// in:body
	Body services.SaveSupertaskResultRequest
}

// Save supertask result success
// swagger:response saveSupertaskResultResponse
type saveSupertaskResultWrapper struct {
	// in:body
	Body services.SaveSupertaskResultResponse
}
