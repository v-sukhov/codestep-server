package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-contest-supertask-list get-contest-supertask-list idOfGetContestSupertaskListEndpoint
// Возвращает список суперзадач контеста, отсортированный в порядке order_number
// responses:
//   200: getContestSupertaskListResponse

// swagger:parameters idOfGetContestSupertaskListEndpoint
type getContestSupertaskListRequestWrapper struct {
	// in:body
	Body services.GetContestSupertaskListRequest
}

// Get contest supertask list sorted by order_number
// swagger:response getContestSupertaskListResponse
type getContestSupertaskListResponseWrapper struct {
	// in:body
	Body services.GetContestSupertaskListResponse
}
