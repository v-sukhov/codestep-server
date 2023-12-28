package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-user-contest-list get-user-contest-list idOfGetUserContestListEndpoint
// Возвращает список контестов пользователя с указанием уровня прав
// responses:
//   200: getUserContestListResponse

// swagger:parameters idOfGetUserContestListEndpoint
type getUserContestListRequestWrapper struct {
	// in:body
	Body services.GetUserContestListRequest
}

// Get user contest list
// swagger:response getUserContestListResponse
type getUserContestListResponseWrapper struct {
	// in:body
	Body services.GetUserContestListResponse
}
