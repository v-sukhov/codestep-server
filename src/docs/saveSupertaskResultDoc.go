package docs

import (
	contest_participation "codestep/services/contest_participation"
)

// swagger:route POST /api/protected/save-supertask-result save-supertask-result idOfSaveSupertaskResultEndpoint
// Сохранить результат тестирования задачи
// responses:
//   200: saveSupertaskResultResponse

// swagger:parameters idOfSaveSupertaskResultEndpoint
type saveSupertaskResultRequestWrapper struct {
	// in:body
	Body contest_participation.SaveSupertaskResultRequest
}

// Save supertask result success
// swagger:response saveSupertaskResultResponse
type saveSupertaskResultWrapper struct {
	// in:body
	Body contest_participation.SaveSupertaskResultResponse
}
