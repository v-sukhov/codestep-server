package docs

import (
	contest_participation "codestep/services/contest_participation"
)

// swagger:route POST /api/protected/upload-add-participation-time-file upload-add-participation-time-file idOfAddParticipationTimeFileEndpoint
// Устанавливает добавленное время участникам соревнования через импорт файла.
// Получает на вход multipart/form-data с полями contestId и participationTimeList.
// Структура файла (разделитель - табуляции или пробелы):
//
//	login	added_time_sec
//
// Операцию может выполнить только владелец или администратор контеста.
// responses:
//   200: addParticipationTimeFileResponse

// swagger:parameters idOfAddParticipationTimeFileEndpoint
type addParticipationTimeFileRequestWrapper struct {
	// in:body
}

// Возвращает список несуществующих логинов и ошибки разбора строк файла.
// swagger:response addParticipationTimeFileResponse
type addParticipationTimeFileResponseWrapper struct {
	// in:body
	Body contest_participation.AddParticipationTimeFileResponse
}
