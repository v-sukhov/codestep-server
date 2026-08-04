package docs

import (
	contest "codestep/services/contest"
)

// swagger:route POST /api/protected/get-contest-user-list get-contest-user-list idOfGetContestUserListEndpoint
// Возвращает список пользователей контеста с пагинацией, сортировкой и фильтрацией.
// sortField: login, role, userStartTime, userFactStartTime, addedTime.
// sortDirection: asc или desc.
// Доступно только владельцам и администраторам контеста.
// responses:
//   200: getContestUserListResponse

// swagger:parameters idOfGetContestUserListEndpoint
type getContestUserListRequestWrapper struct {
	// in:body
	Body contest.GetContestUserListRequest
}

// swagger:response getContestUserListResponse
type getContestUserListResponseWrapper struct {
	// in:body
	Body contest.GetContestUserListResponse
}
