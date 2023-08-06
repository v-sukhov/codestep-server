package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-user-supertask-list get-user-supertask-list idOfGetUserSupertaskListEndpoint
// Get supertask list for current user
// responses:
//   200: getUserSupertaskListResponse

// swagger:parameters idOfGetUserSupertaskListEndpoint
type getUserSupertaskListRequestWrapper struct {
	// in:body
	Body services.GetUserSupertaskListRequest
}

// Save supertask success
// swagger:response getUserSupertaskListResponse
type getUserSupertaskListResponseWrapper struct {
	// in:body
	Body services.GetUserSupertaskListResponse
}
