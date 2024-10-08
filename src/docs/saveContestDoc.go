package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/save-contest save-contest idOfSaveContestEndpoint
// Сохранить контест (перезаписывает данные контеста). Требует прав владельца или администратора контеста.
// Если contestId = 0, значит создаёт новый
// responses:
//   200: saveContest

// swagger:parameters idOfSaveContestEndpoint
type saveContesttWrapper struct {
	// in:body
	Body services.SaveContestRequest
}

// Save contest success
// swagger:response saveContestResponse
type saveContestWrapper struct {
	// in:body
	Body services.SaveContestResponse
}
