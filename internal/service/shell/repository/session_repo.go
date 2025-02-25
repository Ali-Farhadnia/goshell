package repository

import (
	"sync"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell"
)

type SessionRepository struct {
	mu      sync.Mutex
	session *shell.Session
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

func (r *SessionRepository) GetSession() (shell.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return *r.session, nil
}

func (r *SessionRepository) SetSession(session shell.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.session = &session

	return nil
}
