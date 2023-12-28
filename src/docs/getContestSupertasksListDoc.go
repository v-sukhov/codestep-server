package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-contest-supertasks-list get-contest-supertasks-list idOfGetContestSupertasksListEndpoint
// Возвращает список суперзадач контеста
// responses:
//   200: getContestSupertasksListResponse

// swagger:parameters idOfGetContestSupertasksListEndpoint
type getContestSupertasksListRequestWrapper struct {
	// in:body
	Body services.GetContestSupertasksListRequest
}

// Get contest supertask list
// swagger:response getContestSupertasksListResponse
type getContestSupertasksListResponseWrapper struct {
	// in:body
	Body services.GetContestSupertasksListResponse
}
