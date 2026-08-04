package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type GetContestUserListParams struct {
	ContestId      int32
	PageSize       int32
	PageNumber     int32
	SortField      string
	SortDirection  string
	LoginSubstring string
	Roles          []string
}

type ContestUserListItem struct {
	UserId            int32             `json:"userId"`
	Login             string            `json:"login"`
	UserContestRights UserContestRights `json:"userContestRights"`
	UserStartTime     *time.Time        `json:"userStartTime"`
	UserFactStartTime *time.Time        `json:"userFactStartTime"`
	AddedTimeSec      int32             `json:"addedTimeSec"`
}

type GetContestUserListResult struct {
	Users      []ContestUserListItem `json:"users"`
	TotalCount int32                 `json:"totalCount"`
}

const contestUsersBaseQuery = `
	FROM t_user u
	JOIN t_contest_user_right cur ON u.user_id = cur.user_id AND cur.contest_id = $1
	LEFT JOIN t_user_contest_start ucs ON u.user_id = ucs.user_id AND ucs.contest_id = $1
	WHERE ($2 = '' OR u.login ILIKE $2)
	GROUP BY u.user_id, u.login, ucs.user_start_time, ucs.user_fact_start_time, ucs.added_time_sec
`

func GetContestUserList(params GetContestUserListParams) (result GetContestUserListResult, err error) {
	loginPattern := ""
	if params.LoginSubstring != "" {
		loginPattern = "%" + params.LoginSubstring + "%"
	}

	roleFilter := buildContestUserRoleFilter(params.Roles)
	sortColumn, sortNulls := contestUserListSort(params.SortField)
	sortDirection := contestUserListSortDirection(params.SortDirection)
	offset := (params.PageNumber - 1) * params.PageSize

	countQuery := `
		SELECT COUNT(*)
		FROM (
			SELECT u.user_id,
				SUM(1 << (cur.contest_right_type_id - 1)) AS user_contest_rights_bitmask
		` + contestUsersBaseQuery + `
		) contest_users
		WHERE 1 = 1` + roleFilter

	var totalCount int32
	if err = db.QueryRow(countQuery, params.ContestId, loginPattern).Scan(&totalCount); err != nil {
		return
	}

	listQuery := fmt.Sprintf(`
		SELECT user_id, login, user_contest_rights_bitmask, user_start_time, user_fact_start_time, added_time_sec
		FROM (
			SELECT u.user_id,
				u.login,
				SUM(1 << (cur.contest_right_type_id - 1)) AS user_contest_rights_bitmask,
				ucs.user_start_time,
				ucs.user_fact_start_time,
				COALESCE(ucs.added_time_sec, 0) AS added_time_sec
			%s
		) contest_users
		WHERE 1 = 1%s
		ORDER BY %s %s%s
		LIMIT $3 OFFSET $4
	`, contestUsersBaseQuery, roleFilter, sortColumn, sortDirection, sortNulls)

	rows, err := db.Query(listQuery, params.ContestId, loginPattern, params.PageSize, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	result.Users = make([]ContestUserListItem, 0)
	for rows.Next() {
		var item ContestUserListItem
		var rightsBitmask int32
		var userStartTime, userFactStartTime sql.NullTime

		if err = rows.Scan(
			&item.UserId,
			&item.Login,
			&rightsBitmask,
			&userStartTime,
			&userFactStartTime,
			&item.AddedTimeSec,
		); err != nil {
			return
		}

		item.UserContestRights = NewUserContestRights(rightsBitmask)
		item.UserStartTime = nullTimePtr(userStartTime)
		item.UserFactStartTime = nullTimePtr(userFactStartTime)
		result.Users = append(result.Users, item)
	}

	if err = rows.Err(); err != nil {
		return
	}

	result.TotalCount = totalCount
	return
}

func buildContestUserRoleFilter(roles []string) string {
	if len(roles) == 0 {
		return ""
	}

	roleBits := map[string]int32{
		"OWNER":       1,
		"ADMIN":       2,
		"JURY":        4,
		"PARTICIPANT": 8,
	}

	conditions := make([]string, 0, len(roles))
	for _, role := range roles {
		if bit, ok := roleBits[strings.ToUpper(strings.TrimSpace(role))]; ok {
			conditions = append(conditions, fmt.Sprintf("(user_contest_rights_bitmask & %d) > 0", bit))
		}
	}

	if len(conditions) == 0 {
		return ""
	}

	return " AND (" + strings.Join(conditions, " OR ") + ")"
}

func contestUserListSort(sortField string) (column string, nullsOrder string) {
	switch strings.ToLower(sortField) {
	case "role", "roles":
		return "user_contest_rights_bitmask", ""
	case "userstarttime":
		return "user_start_time", " NULLS LAST"
	case "userfactstarttime":
		return "user_fact_start_time", " NULLS LAST"
	case "addedtime":
		return "added_time_sec", ""
	default:
		return "login", ""
	}
}

func contestUserListSortDirection(sortDirection string) string {
	if strings.ToLower(sortDirection) == "desc" {
		return "DESC"
	}
	return "ASC"
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}
