package bot

import (
	"sync"

	"github.com/vozisov/cargo-bot/order"
)

type State int

const (
	StateIdle State = iota
	StateFromCity
	StateToCity
	StateCargo
	StateDateSelect
	StateDateManual
	StateContact
)

// UserSession — состояние одного пользователя
type UserSession struct {
	State State
	Order order.Order // черновик заказа
}

// SessionStore — хранилище сессий в памяти.
// Позже заменим на Redis, когда будем масштабироваться.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[int64]*UserSession
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[int64]*UserSession),
	}
}

func (s *SessionStore) Get(userID int64) *UserSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if session, ok := s.sessions[userID]; ok {
		return session
	}
	return &UserSession{State: StateIdle}
}

func (s *SessionStore) Set(userID int64, session *UserSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[userID] = session
}

func (s *SessionStore) Reset(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, userID)
}
