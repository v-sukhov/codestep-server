package docs

import (
	contest_participation "codestep/services/contest_participation"
)

// swagger:route POST /api/protected/get-supertask-in-contest-with-results get-supertask-in-contest-with-results idOfGetSupertaskInContestWithResultsEndpoint
// Получить данные актуальной версии суперзадачи с результатами запрашивающего пользователя
// responses:
//   200: getSupertaskInContestWithResultsResponse

// swagger:parameters idOfGetSupertaskInContestWithResultsEndpoint
type getSupertaskInContestWithResultsRequestWrapper struct {
	// in:body
	Body contest_participation.GetSupertaskInContestWithResultsRequest
}

// Get supertask all tasks results success
// swagger:response getSupertaskInContestWithResultsResponse
type getSupertaskInContestWithResultsWrapper struct {
	// in:body
	Body contest_participation.GetSupertaskInContestWithResultsResponse
}
