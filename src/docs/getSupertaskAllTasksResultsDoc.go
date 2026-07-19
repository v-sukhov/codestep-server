package docs

import (
	supertask "codestep/services/supertask"
)

// swagger:route POST /api/protected/get-supertask-all-tasks-results get-supertask-all-tasks-results idOfGetSupertaskAllTasksResultsEndpoint
// Получить результаты по всем задачам суперзадачи
// responses:
//   200: getSupertaskAllTasksResultsResponse

// swagger:parameters idOfGetSupertaskAllTasksResultsEndpoint
type getSupertaskAllTasksResultsRequestWrapper struct {
	// in:body
	Body supertask.GetSupertaskAllTasksResultsRequest
}

// Get supertask all tasks results success
// swagger:response getSupertaskAllTasksResultsResponse
type getSupertaskAllTasksResultsWrapper struct {
	// in:body
	Body supertask.GetSupertaskAllTasksResultsResponse
}

