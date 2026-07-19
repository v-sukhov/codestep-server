package docs

import (
	supertask "codestep/services/supertask"
)

// swagger:route POST /api/protected/get-supertask-versions get-supertask-versions idOfGetSupertaskVersionsEndpoint
// Возвращает список версий суперзадачи
// responses:
//   200: getSupertaskVersionsResponse

// swagger:parameters idOfGetSupertaskVersionsEndpoint
type getSupertaskVersionsRequestWrapper struct {
	// in:body
	Body supertask.GetSupertaskVersionsRequest
}

// Response
// swagger:response getSupertaskVersionsResponse
type getSupertaskVerionsWrapper struct {
	// in:body
	Body supertask.GetSupertaskVersionsResponse
}

