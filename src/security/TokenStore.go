package security

/*
	НЕ ИСПОЛЬЗУЕТСЯ
*/

import (
	"crypto/rand"
	"encoding/hex"
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
	Generate token
	****************************************
*/

func _generateToken() string {
	len := 256
	b := make([]byte, len)

	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
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

func _getUserByToken(token string) (UserInfoCache, bool) {
	userTokenMap.mx.RLock()
	user, ok := userTokenMap.m[token]
	userTokenMap.mx.RUnlock()
	return user, ok
}

// if token already exists return false and do nothing
func _addUserByToken(token string, user UserInfoCache) bool {
	userTokenMap.mx.Lock()

	var res bool
	if _, ok := userTokenMap.m[token]; ok {
		res = false
	} else {
		userTokenMap.m[token] = user
		res = true
	}

	userTokenMap.mx.Unlock()
	return res
}

// add user and return new token
func _addUser(user UserInfoCache) string {

	var token string
	for token = ""; token == ""; {
		token = _generateToken()
		if ok := _addUserByToken(token, user); !ok {
			token = ""
		}
	}

	return token
}

// delete token
func _deleteToken(token string) {
	userTokenMap.mx.Lock()
	delete(userTokenMap.m, token)
	userTokenMap.mx.Unlock()
}

// init
func init() {
	userTokenMap.m = make(map[string]UserInfoCache)
}
