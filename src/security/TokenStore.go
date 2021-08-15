package security

import (
	"sync"
)

const (
	REGULAR     int = 1
	PARTICIPANT int = 2
)

type UserInfoCache struct {
	Id       int
	UserType int // REGULAR or PARTICIPANT
}

/*
	****************************************
	userByTokenMap
	****************************************
*/

// thread-safe map for indexing user info by token
type userByTokenMap struct {
	mx sync.RWMutex
	m  map[string]UserInfoCache
}

var userTokenMap userByTokenMap

func getUserByToken(token string) (UserInfoCache, bool) {
	userTokenMap.mx.RLock()
	user, ok := userTokenMap.m[token]
	userTokenMap.mx.RUnlock()
	return user, ok
}

// if token already exists return false and do nothing
func addUserByToken(token string, user UserInfoCache) {
	userTokenMap.mx.Lock()
	userTokenMap.m[token] = user
	userTokenMap.mx.Unlock()
}

// delete token
func deleteToken(token string) {
	userTokenMap.mx.Lock()
	delete(userTokenMap.m, token)
	userTokenMap.mx.Unlock()
}

// init
func init() {
	userTokenMap.m = make(map[string]UserInfoCache)
}
