package docs

import (
	supertask "codestep/services/supertask"
)

// swagger:route POST /api/protected/save-supertask save-supertask idOfSaveSupertaskEndpoint
// Save supertask
// Если SupertaskId = 0 - создаётся новый supertask
// Если Commited = true создаётся закоммиченная неизменяемая версия supertask.
// Ей присвивается новый номер версии. При этом указывается ParentVersionNumber - от какой версии произведена.
// Если ParentVersionNumber = 0 - то это считается первоначальной версией, то есть которая ни от чего не произведена.
// Если Commited = false создаётся незакоммиченная версия данного пользователя.
// Незакоммиченная версия у данного ползователя может быть только одна.
// responses:
//   200: saveSupertaskResponse

// swagger:parameters idOfSaveSupertaskEndpoint
type saveSupertaskRequestWrapper struct {
	// in:body
	Body supertask.SaveSupertaskRequest
}

// Save supertask success
// swagger:response saveSupertaskResponse
type saveSupertaskWrapper struct {
	// in:body
	Body supertask.SaveSupertaskResponse
}

