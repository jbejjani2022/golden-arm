package internal

import (
	"sync"
	"time"
)

var sessionStore = struct {
	sync.RWMutex
	data map[string]SessionData
}{
	data: make(map[string]SessionData),
}

type SessionData struct {
	User      string
	ExpiresAt time.Time
}

func StoreSession(token string, user string, expiresAt time.Time) {
	sessionStore.Lock()
	defer sessionStore.Unlock()
	sessionStore.data[token] = SessionData{
		User:      user,
		ExpiresAt: expiresAt,
	}
}

func ValidateSession(token string) bool {
	sessionStore.RLock()
	defer sessionStore.RUnlock()
	session, exists := sessionStore.data[token]
	if !exists || time.Now().After(session.ExpiresAt) {
		return false
	}
	return true
}

func DeleteSession(token string) {
	sessionStore.Lock()
	defer sessionStore.Unlock()
	delete(sessionStore.data, token)
}
