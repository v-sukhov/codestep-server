package docs

// swagger:route POST /api/protected/upload-create-user-list upload-create-user-list idOfUploadCreateUserListEndpoint
// Позволяет загрузить файл со списком логинов для создания пользователей. Каждый логин - в отдельной строке. Логины состоят из символов a-z, 0-9, - и _
// responses:
//   200: uploadCreateUserListResponse

// swagger:parameters idOfUploadCreateUserListEndpoint
type uploadCreateUserListRequestWrapper struct {
	// in:body
}

// Возвращает файл со списком логин/пароль
// swagger:response uploadCreateUserListResponse
type uploadCreateUserListResponseWrapper struct {
	// out:body
}
