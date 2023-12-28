package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/remove-supertask-from-contest remove-supertask-from-contest idOfRemoveSupertaskFromContestEndpoint
// Удаляет суперзадачу из контеста
// responses:
//   200: removeSupertaskFromContestResponse

// swagger:parameters idOfRemoveSupertaskFromContestEndpoint
type removeSupertaskFromContestRequestWrapper struct {
	// in:body
	Body services.RemoveSupertaskFromContestRequest
}

// Remove supertask from contest
// swagger:response removeSupertaskFromContestResponse
type removeSupertaskFromContestResponseWrapper struct {
	// in:body
	Body services.RemoveSupertaskFromContestResponse
}
