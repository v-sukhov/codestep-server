package docs

import (
	admin "codestep/services/admin"
)

// swagger:route POST /api/protected/create-contest-users create-contest-users idOfCreateContestUsersEndpoint
// Массово создаёт пользователей и предоставляет им права на контест.
// Логины формируются как loginPrefix + номер с дополнением нулями слева.
// roles: ADMIN, JURY, PARTICIPANT. Не более 1000 пользователей за один запрос.
// Доступно только владельцам и администраторам контеста.
// responses:
//   200: createContestUsersResponse

// swagger:parameters idOfCreateContestUsersEndpoint
type createContestUsersRequestWrapper struct {
	// in:body
	Body admin.CreateContestUsersRequest
}

// Возвращает файл со списком логин/пароль и ошибками
// swagger:response createContestUsersResponse
type createContestUsersResponseWrapper struct {
	// out:body
}
