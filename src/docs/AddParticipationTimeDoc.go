package docs

import (
	contest_participation "codestep/services/contest_participation"
)

// swagger:route POST /api/protected/add-participation-time add-participation-time idOfAddParticipationTimeEndpoint
// Устанавливает добавленное время участникам соревнования по userId.
// Операцию может выполнить только владелец или администратор контеста.
// responses:
//   200: addParticipationTimeResponse

// swagger:parameters idOfAddParticipationTimeEndpoint
type addParticipationTimeRequestWrapper struct {
	// in:body
	Body contest_participation.AddParticipationTimeRequest
}

// Возвращает результат установки добавленного времени участникам соревнования.
// swagger:response addParticipationTimeResponse
type addParticipationTimeResponseWrapper struct {
	// in:body
	Body contest_participation.AddParticipationTimeResponse
}
