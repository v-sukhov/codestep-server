package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/save-supertask save-supertask idOfSaveSupertaskEndpoint
// Save supertask
// responses:
//   200: saveSupertaskResponse

// swagger:parameters idOfSaveSupertaskEndpoint
type saveSupertaskRequestWrapper struct {
	// in:body
	Body services.SaveSupertaskRequest
}

// Save supertask success
// swagger:response saveSupertaskResponse
type saveSupertaskWrapper struct {
	// in:body
	Body services.SaveSupertaskResponse
}
