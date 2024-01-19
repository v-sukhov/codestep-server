package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-supertask-in-contest-with-results get-supertask-in-contest-with-results idOfGetSupertaskInContestWithResultsEndpoint
// Получить данные актуальной версии суперзадачи с результатами запрашивающего пользователя
// responses:
//   200: getSupertaskInContestWithResultsResponse

// swagger:parameters idOfGetSupertaskInContestWithResultsEndpoint
type getSupertaskInContestWithResultsRequestWrapper struct {
	// in:body
	Body services.GetSupertaskInContestWithResultsRequest
}

// Get supertask all tasks results success
// swagger:response getSupertaskInContestWithResultsResponse
type getSupertaskInContestWithResultsWrapper struct {
	// in:body
	Body services.GetSupertaskInContestWithResultsResponse
}
