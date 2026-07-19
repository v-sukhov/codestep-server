package docs

import (
	contest_participation "codestep/services/contest_participation"
)

// swagger:route POST /api/protected/get-contest-participation-info get-contest-participation-info idOfGetContestParticipationInfoEndpoint
// Возвращает информацию об участии пользователя(ей) в контесте.
// Без userIdList возвращает информацию по вызвавшему пользователю.
// С userIdList возвращает информацию по указанным пользователям (только для владельца и администратора контеста).
// responses:
//   200: getContestParticipationInfoResponse

// swagger:parameters idOfGetContestParticipationInfoEndpoint
type getContestParticipationInfoRequestWrapper struct {
	// in:body
	Body contest_participation.GetContestParticipationInfoRequest
}

// swagger:response getContestParticipationInfoResponse
type getContestParticipationInfoResponseWrapper struct {
	// in:body
	Body contest_participation.GetContestParticipationInfoResponse
}
