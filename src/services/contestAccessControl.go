package services

import (
	"codestep/db"
)

/*
CheckContestTimeAccess проверяет доступ пользователя к контесту по времени.

Логика проверки:
 1. Для владельца и администратора (права 1 и 2) - доступ всегда разрешен
 2. Для жюри (право 4) - доступ разрешен только после старта контеста
 3. Для участника (право 8) - доступ разрешен только во время участия,
    при необходимости автоматически стартует контест для пользователя

Возвращает:
- access: true если доступ разрешен, false если запрещен
- info: информация о статусе участия пользователя в контесте
- err: ошибка при проверке прав или статуса участия
*/
func CheckContestTimeAccess(userId int32, contestId int32) (access bool, info db.UserContestParticipationInfo, err error) {
	userContestRights, err := db.GetContestUserRights(userId, contestId)
	if err != nil {
		return false, info, err
	}

	if userContestRights == 0 {
		return false, info, nil
	}

	info, err = db.GetUserContestParticipationStatus(userId, contestId)
	if err != nil {
		return false, info, err
	}

	// Владелец (1) и администратор (2) - доступ всегда разрешен
	if userContestRights&3 > 0 {
		access = true
		return
	}

	// Жюри (4) - доступ только после старта контеста
	if userContestRights&4 > 0 {
		switch info.Status {
		case db.ContestNotStarted:
			access = false
		default:
			access = true
		}
		return
	}

	// Участник (8) - доступ только во время участия
	if userContestRights&8 > 0 {
		switch info.Status {
		case db.UserParticipating:
			access = true
		case db.UserNotStartedStartAllowed:
			// Автоматически стартуем контест для пользователя
			if err := db.StartContestForUser(userId, contestId); err != nil {
				return false, info, err
			}
			info.Status = db.UserParticipating
			access = true
		default:
			access = false
		}
		return
	}

	// Если пользователь не имеет ни одного из перечисленных прав
	access = false
	return
}
