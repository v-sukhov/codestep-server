package docs

import (
	supertask "codestep/services/supertask"
)

// swagger:route POST /api/protected/get-user-supertask-list get-user-supertask-list idOfGetUserSupertaskListEndpoint
// Get supertask list for current user
// responses:
//   200: getUserSupertaskListResponse

// swagger:parameters idOfGetUserSupertaskListEndpoint
type getUserSupertaskListRequestWrapper struct {
	// in:body
	Body supertask.GetUserSupertaskListRequest
}

// Save supertask success
// swagger:response getUserSupertaskListResponse
type getUserSupertaskListResponseWrapper struct {
	// in:body
	Body supertask.GetUserSupertaskListResponse
}

