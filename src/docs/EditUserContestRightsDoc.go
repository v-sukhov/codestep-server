package docs

import (
	contest_participation "codestep/services/contest_participation"
)

// swagger:route POST /api/protected/edit-user-contest-rights edit-user-contest-rights idOfEditUserContestRightsEndpoint
// Полная замена прав пользователей на контест по userId.
// roles: ADMIN, JURY, PARTICIPANT. OWNER назначить или отозвать нельзя — он сохраняется у пользователя.
// Операцию может выполнить только владелец или администратор контеста.
// responses:
//   200: editUserContestRightsResponse

// swagger:parameters idOfEditUserContestRightsEndpoint
type editUserContestRightsRequestWrapper struct {
	// in:body
	Body contest_participation.EditUserContestRightsRequest
}

// Возвращает результат изменения прав пользователей на контест.
// swagger:response editUserContestRightsResponse
type editUserContestRightsResponseWrapper struct {
	// in:body
	Body contest_participation.EditUserContestRightsResponse
}
