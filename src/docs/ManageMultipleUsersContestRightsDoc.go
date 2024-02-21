package docs

import (
	services "codestep/services"
)

// swagger:route POST /api/protected/upload-manage-multiple-users-contest-rights upload-manage-multiple-users-contest-rights idOfManageMultipleUsersContestRightsEndpoint
// Управление правами пользователей на контесты: позволяет выдавать и забирать права пользователей на контесты
// Получает на вход файл следующей структуры (разделитель - табуляции или пробелы):

//	login	contest_id	[revoke|admin|jury|participant]*

//	revoke - отнимаются все права, с опцией revoke недопустимо задавать другие варианты прав
//	По умолчанию (если права не перечислены, считается, что заданы права participant)
//	Если один и тот же логин встречается в файле несколько раз, то каждая комбинация логин-контест-тип прав учитывается один раз
//  При изменении прав права владельца контеста всегда остаются без изменений
// responses:
//   200: manageMultipleUsersContestRightsResponse

// swagger:parameters idOfManageMultipleUsersContestRightsEndpoint
type manageMultipleUsersContestRightsRequestWrapper struct {
	// in:body
}

// Возвращает результаты задания прав пользователей на контесты. Если файл удалось успешно обработать и не было ошибок БД возвращает Success = true
// Все проблемные строки, по которым не удалось выполнить операцию возвращает в ответе с указанием номера строки в файле и описанием проблемы.
// swagger:response manageMultipleUsersContestRightsResponse
type manageMultipleUsersContestRightsResponseWrapper struct {
	// in:body
	Body services.ManageMultipleUsersContestRightsResponse
}
