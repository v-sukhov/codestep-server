package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-supertask get-supertask idOfGetSupertaskEndpoint
// Get supertask
// responses:
//   200: getSupertaskResponse

// swagger:parameters idOfGetSupertaskEndpoint
type getSupertaskRequestWrapper struct {
	// in:body
	Body services.GetSupertaskRequest
}

// Save supertask success
// swagger:response getSupertaskResponse
type getSupertaskWrapper struct {
	// in:body
	Body services.GetSupertaskResponse
}
