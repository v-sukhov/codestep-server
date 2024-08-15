package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-supertask-versions get-supertask-versions idOfGetSupertaskVersionsEndpoint
// Возвращает список версий суперзадачи
// responses:
//   200: getSupertaskVersionsResponse

// swagger:parameters idOfGetSupertaskVersionsEndpoint
type getSupertaskVersionsRequestWrapper struct {
	// in:body
	Body services.GetSupertaskVersionsRequest
}

// Response
// swagger:response getSupertaskVersionsResponse
type getSupertaskVerionsWrapper struct {
	// in:body
	Body services.GetSupertaskVersionsResponse
}
