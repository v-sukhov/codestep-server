package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/add-supertask-to-contest add-supertask-to-contest idOfAddSupertaskToContestEndpoint
// Добавить суперзадачу в контест
// responses:
//   200: addSupertaskToContestResponse

// swagger:parameters idOfAddSupertaskToContestEndpoint
type addSupertaskToContestRequestWrapper struct {
	// in:body
	Body services.AddSupertaskToContestRequest
}

// Add supertask to contest
// swagger:response addSupertaskToContestResponse
type addSupertaskToContestResponseWrapper struct {
	// in:body
	Body services.AddSupertaskToContestResponse
}
