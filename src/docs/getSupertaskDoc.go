package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/get-supertask get-supertask idOfGetSupertaskEndpoint
// Возвращает объект версии суперзадачи
// Если versionNumber присутствует и != 0 - возвращает соответствующую закоммиченную версию
// Иначе если authorUserId присутствует и != 0 - возвращает рабочую версию соответствующего пользователя
// Иначе возвращает рабочую версию запрашивающего пользователя
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
