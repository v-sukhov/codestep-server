package docs

import (
	contest "codestep/services/contest"
)

// swagger:route POST /api/protected/add-supertask-to-contest add-supertask-to-contest idOfAddSupertaskToContestEndpoint
// НЕ ИСПОЛЬЗУЕТСЯ!!! Добавить суперзадачу в контест
// responses:
//   200: addSupertaskToContestResponse

// swagger:parameters idOfAddSupertaskToContestEndpoint
type addSupertaskToContestRequestWrapper struct {
	// in:body
	Body contest.AddSupertaskToContestRequest
}

// Add supertask to contest
// swagger:response addSupertaskToContestResponse
type addSupertaskToContestResponseWrapper struct {
	// in:body
	Body contest.AddSupertaskToContestResponse
}

