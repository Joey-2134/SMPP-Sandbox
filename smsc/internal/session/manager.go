package session

import (
	"net"
	"sync"
)

type Manager struct {
	mu       sync.RWMutex
	sessions map[net.Conn]*Session
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[net.Conn]*Session),
	}
}

func (m *Manager) AddSession(conn net.Conn) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	session := NewSession(conn)
	m.sessions[conn] = session
	return session
}

func (m *Manager) RemoveSession(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, conn)
}

func (m *Manager) GetSession(conn net.Conn) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[conn]
}

func (m *Manager) GetAllSessions() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sessions := make([]*Session, 0, len(m.sessions))
	for _, session := range m.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}
