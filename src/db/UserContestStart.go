package db

/**
Реализует взаимодействие с таблицей T_USER_CONTEST_START

1. SetUserContestStart - записывает время начала участия пользователя в соревновании
2. GetUserContestStart - получает время начала участия пользователя в соревновании
3. StartContestForUser - стартует соревнование для пользователя: создаёт запись в таблице T_USER_CONTEST_START с текущим временем фактического старта
	и отсчётным временем, определяемым следующей логикой:
	- если соревнование не допускает виртуальное участие (плавающий старт),
	то отсчётное время участия в соревновании определяется по полю CONTEST_START_TIME таблицы T_CONTEST.
	- иначе отсчётное время совпадает с фактическим временем старта.
4. GetUserContestParticipationStatus - проверяет статус участия пользователя в соревновании, работает в соответствии со следующей логикой:
	- если соревнование не началось, то возвращает ContestNotStarted
	- если соревнование закончено, но пользователь не стартовал его, то возвращает ContestFinished
	- если соревнование началось, ещё не закончилось, но пользователь не стартовал его,
	то возвращает UserNotStartedStartAllowed или UserNotStartedStartNotAllowed:
			- Если соревнование не разрешает плавающий старт и оно ещё не закончено, то UserNotStartedStartAllowed
			- Если соревнование разрешает плавающий старт и текущее время не превышает LIMIT_VIRTUAL_START, то UserNotStartedStartAllowed
			- Иначе UserNotStartedStartNotAllowed
	- если соревнование началось и пользователь стартовал его, то возвращает UserParticipating
	- если соревнование началось и пользователь закончил его, то возвращает UserFinished
	Время старта соревнования определяется по полю CONTEST_START_TIME таблицы T_CONTEST.
	Время окончания соревнования для которого не разрешён плавающий старт определяется по полям
	CONTEST_START_TIME + CONTEST_DURATION (длительность соревнования в минутах)
	и NO_END_TIME таблицы T_CONTEST. Если NO_END_TIME равно TRUE, то соревнование бесконечно.
	Время окончания соревнования для которого разрешён плаваюший старта определяется по полям LIMIT_VIRTUAL_START + CONTEST_DURATION.
	Время начала участия пользователя в соревновании определяется по полю USER_START_TIME таблицы T_USER_CONTEST_START.
	Если запись для данного пользователя и соревнования не существует, то пользователь не стартовал соревнование.
	Время окончания участия пользователя в соревновании определяется по полям USER_START_TIME таблицы T_USER_CONTEST_START
	и CONTEST_DURATION таблицы T_CONTEST.

	Возвращает вычисленный статус.

5. TryUserContestAction - пытается выполнить действие над соревнованием для пользователя, работает в соответствии со следующей логикой:
	- Проверяет статус участия пользователя в соревновании
	- Если статус равен ContestNotStarted, ContestFinished, UserFinished, UserNotStartedStartNotAllowed или UserParticipating
	просто возвращает этот статус
	- Иначе выполняет StartContestForUser для данного контеста и данного пользователя и возвращает статус UserParticipating
*/
import (
	"database/sql"
	"fmt"
	"time"
)

type UserContestParticipationStatus int // статус участия пользователя в соревновании

const (
	ContestNotStarted             UserContestParticipationStatus = iota // соревнование не началось
	ContestFinished                                                     // соревнование закончено, но пользователь не стартовал его
	UserNotStartedStartAllowed                                          // соревнование началось, но пользователь не стартовал его, старт разрешён
	UserNotStartedStartNotAllowed                                       // соревнование началось, но пользователь не стартовал его, старт не разрешён
	UserParticipating                                                   // пользователь участвует в соревновании
	UserFinished                                                        // пользователь закончил соревнование
)

func (s UserContestParticipationStatus) String() string {
	switch s {
	case ContestNotStarted:
		return "ContestNotStarted"
	case ContestFinished:
		return "ContestFinished"
	case UserNotStartedStartAllowed:
		return "UserNotStartedStartAllowed"
	case UserNotStartedStartNotAllowed:
		return "UserNotStartedStartNotAllowed"
	case UserParticipating:
		return "UserParticipating"
	case UserFinished:
		return "UserFinished"
	default:
		return fmt.Sprintf("UnknownStatus(%d)", int(s))
	}
}

type UserContestParticipationInfo struct {
	// Статус участия пользователя в соревновании
	Status UserContestParticipationStatus `json:"status"`

	// Текущее время
	CurrentTime time.Time `json:"currentTime"`

	// Время отсчётного старта пользователя (user_start_time)
	UserStartTime time.Time `json:"userStartTime"`

	// Фактическое время старта пользователя (user_fact_start_time)
	UserFactStartTime time.Time `json:"userFactStartTime"`

	// Время старта соревнования (contest_start_time)
	ContestStartTime time.Time `json:"contestStartTime"`

	// Длительность соревнования в минутах (contest_duration)
	ContestDuration int32 `json:"contestDuration"`

	// Признак бесконечного соревнования (no_end_time)
	NoEndTime bool `json:"noEndTime"`

	// Признак виртуального участия (virtual_participation)
	VirtualParticipation bool `json:"virtualParticipation"`

	// Признак ограничения окна виртуального старта (limit_virtual_start)
	LimitVirtualStart bool `json:"limitVirtualStart"`

	// Время окончания окна виртуального старта (virtual_start_end_time)
	VirtualStartEndTime time.Time `json:"virtualStartEndTime"`
}

// StartContestForUser стартует соревнование для пользователя, не проверяет ограничения на время старта
func StartContestForUser(
	userId int32,
	contestId int32,
) error {
	now := time.Now().UTC()
	var userStartTime time.Time
	var contestStartTime sql.NullTime
	var virtualParticipation bool
	err := db.QueryRow(`
		SELECT contest_start_time, virtual_participation
		FROM t_contest
		WHERE contest_id = $1
	`, contestId).Scan(&contestStartTime, &virtualParticipation)
	if err != nil {
		return err
	}
	if !virtualParticipation {
		// Если не разрешён плавающий старт, отсчётное время = время старта соревнования
		if !contestStartTime.Valid {
			return fmt.Errorf("Contest start not allowed: contest start time is not set for contest_id=%d", contestId)
		}
		userStartTime = contestStartTime.Time
	} else {
		// Иначе отсчётное время = фактическое время старта
		userStartTime = now
	}

	// Вставляем запись о старте пользователя
	_, err = db.Exec(`
		INSERT INTO t_user_contest_start (user_id, contest_id, user_start_time, user_fact_start_time)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, contest_id) DO NOTHING
	`, userId, contestId, userStartTime, now)
	return err
}

// GetUserContestParticipationStatus возвращает статус участия пользователя в соревновании
// Функция делает 2 запроса к БД
// Предполагается, что она должна редко вызываться, а в основном при работе с контестом информация должна браться из кэша
// Кэша нет, его ещё надо сделать
func GetUserContestParticipationStatus(userId int32, contestId int32) (info UserContestParticipationInfo, err error) {
	// Получаем информацию о соревновании
	var contestStartTime sql.NullTime
	var contestDuration int32
	var noEndTime bool
	var virtualParticipation bool
	var limitVirtualStart bool
	var virtualStartEndTime sql.NullTime

	err = db.QueryRow(`
		SELECT contest_start_time, contest_duration, no_end_time, virtual_participation, limit_virtual_start, virtual_start_end_time
		FROM t_contest
		WHERE contest_id = $1
	`, contestId).Scan(&contestStartTime, &contestDuration, &noEndTime, &virtualParticipation, &limitVirtualStart, &virtualStartEndTime)
	if err != nil {
		return
	}

	now := time.Now().UTC()

	// Определяем время начала соревнования, окончания соревнования и время окончания возможности старта соревнования
	var contestBegin, contestEnd, contestStartEnd time.Time
	if contestStartTime.Valid {
		contestBegin = contestStartTime.Time
		if noEndTime {
			contestEnd = time.Time{} // - нет времени завершения
			contestStartEnd = time.Time{}
		} else if !virtualParticipation { // не разрешён плавающий старт
			contestEnd = contestBegin.Add(time.Duration(contestDuration) * time.Minute)
			contestStartEnd = time.Time{}
		} else if !limitVirtualStart || !virtualStartEndTime.Valid { // разрешён плавающий старт и он не ограничен по времени
			contestEnd = time.Time{}      // нет времени завершения соревнования
			contestStartEnd = time.Time{} // не ограничено время самостоятельного старта
		} else { // разрешён плавающий старт и он ограничен по времени
			contestStartEnd = virtualStartEndTime.Time
			contestEnd = contestStartEnd.Add(time.Duration(contestDuration) * time.Minute)
		}
	} else {
		// Если не задано время старта, считаем что не началось - время начало пустое
		contestBegin = time.Time{}
		contestEnd = time.Time{}
		contestStartEnd = time.Time{}
	}

	// Получаем запись о старте пользователя
	var userStartTime sql.NullTime
	var userFactStartTime sql.NullTime
	err = db.QueryRow(`
		SELECT user_start_time, user_fact_start_time
		FROM t_user_contest_start
		WHERE user_id = $1 AND contest_id = $2
	`, userId, contestId).Scan(&userStartTime, &userFactStartTime)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	userHasStarted := (err == nil && userStartTime.Valid)
	err = nil

	info.CurrentTime = now
	info.ContestStartTime = contestStartTime.Time
	info.ContestDuration = contestDuration
	info.NoEndTime = noEndTime
	info.VirtualParticipation = virtualParticipation
	info.LimitVirtualStart = limitVirtualStart
	info.VirtualStartEndTime = virtualStartEndTime.Time
	info.UserStartTime = userStartTime.Time
	info.UserFactStartTime = userFactStartTime.Time

	if contestBegin.IsZero() || now.Before(contestBegin) {
		// 1. Если соревнование не началось или не задано время старта - возвращаем ContestNotStarted
		info.Status = ContestNotStarted
		return
	} else if !contestEnd.IsZero() && contestEnd.Before(now) {
		// 2. Если определено время завершения контеста и оно в прошлом - контест закончился
		info.Status = ContestFinished
		return
	} else if !userHasStarted {
		// 3. Соревнование началось, ещё не закончилось и пользователь ещё не стартовал
		if !virtualParticipation {
			// Если запрещён плавающий старт - можно стартовать
			info.Status = UserNotStartedStartAllowed
			return
		} else {
			// Если разрешён плавающий старт - тогда можно стартовать, если либо время завершения плавающего старта не ограничено,
			// либо оно ещё не прошло
			if !limitVirtualStart || !virtualStartEndTime.Time.Before(now) {
				info.Status = UserNotStartedStartAllowed
				return
			} else {
				info.Status = UserNotStartedStartNotAllowed
				return
			}
		}
	} else {
		// 4. Соревнование началось, ещё не закончилось и пользователь его стартовал
		// Пользователь участвует в соревновании, если оно бесконечное или ещё не вышло время соревнования с момента его
		// отсчётного старта - логика одинаково работает как для обычных соревнований, так и для соревнований с плавающим стартом
		if noEndTime || now.Before(userStartTime.Time.Add(time.Duration(contestDuration)*time.Minute)) {
			info.Status = UserParticipating
			return
		} else {
			info.Status = UserFinished
			return
		}
	}
}

// TryUserContestAction пытается выполнить действие над соревнованием для пользователя
// Если пользователь может участвовать в соревновании, но ещё не начал его - фиксирует время начала участия
// Возвращает статус участия пользователя
func TryUserContestAction(userId int32, contestId int32) (UserContestParticipationStatus, error) {
	info, err := GetUserContestParticipationStatus(userId, contestId)
	if err != nil {
		return info.Status, err
	}
	switch info.Status {
	case ContestNotStarted, ContestFinished, UserFinished, UserNotStartedStartNotAllowed, UserParticipating:
		return info.Status, nil
	case UserNotStartedStartAllowed:
		// Стартуем соревнование для пользователя
		err := StartContestForUser(userId, contestId)
		if err != nil {
			return info.Status, err
		}
		return UserParticipating, nil
	default:
		return info.Status, nil
	}
}

// SetUserContestStart записывает время начала участия пользователя в соревновании.
// Если запись уже существует, обновляем время отсчёта старта и фактического старта.
func SetUserContestStart(userId int32, contestId int32, startTime time.Time) error {
	// Попробуем вставить запись, если она уже есть - обновляем время отсчёта старта и фактического старта данным значением параметра
	_, err := db.Exec(`
		INSERT INTO t_user_contest_start (user_id, contest_id, user_start_time, user_fact_start_time)
		VALUES ($1, $2, $3, $3)
		ON CONFLICT (user_id, contest_id) DO UPDATE
		SET user_start_time = EXCLUDED.user_start_time,
		    user_fact_start_time = EXCLUDED.user_fact_start_time
	`, userId, contestId, startTime)
	return err
}

// GetUserContestStart получает время начала участия пользователя в соревновании.
// Если записи нет, возвращает (nil, nil).
// GetUserContestStart получает время начала участия пользователя в соревновании, а также фактическое время начала.
// Если записи нет, возвращает (nil, nil, nil).
func GetUserContestStart(userId int32, contestId int32) (*time.Time, *time.Time, error) {
	var startTime sql.NullTime
	var factStartTime sql.NullTime
	err := db.QueryRow(`
		SELECT user_start_time, user_fact_start_time
		FROM t_user_contest_start
		WHERE user_id = $1 AND contest_id = $2
	`, userId, contestId).Scan(&startTime, &factStartTime)
	if err == sql.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	var startPtr, factPtr *time.Time
	if startTime.Valid {
		startPtr = &startTime.Time
	}
	if factStartTime.Valid {
		factPtr = &factStartTime.Time
	}
	return startPtr, factPtr, nil
}
