package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Session struct {
	ID           string
	User         User
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]Session
	path     string
}

func NewSessionStore(path string) (*SessionStore, error) {
	store := &SessionStore{
		sessions: make(map[string]Session),
		path:     path,
	}

	if path == "" {
		return store, nil
	}

	if err := store.loadFromDisk(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *SessionStore) Create(user User, token Token) (Session, error) {
	id, err := randomToken(32)
	if err != nil {
		return Session{}, err
	}

	expiresAt := token.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt = time.Now().Add(1 * time.Hour)
	}

	session := Session{
		ID:           id,
		User:         user,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    expiresAt,
	}

	s.mu.Lock()
	s.sessions[id] = session
	err = s.saveLocked()
	s.mu.Unlock()
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (s *SessionStore) Get(id string) (Session, bool) {
	s.mu.RLock()
	session, ok := s.sessions[id]
	s.mu.RUnlock()
	return session, ok
}

func (s *SessionStore) UpdateTokens(id string, token Token) (Session, bool, error) {
	s.mu.Lock()
	session, ok := s.sessions[id]
	if !ok {
		s.mu.Unlock()
		return Session{}, false, nil
	}

	session.AccessToken = token.AccessToken
	if token.RefreshToken != "" {
		session.RefreshToken = token.RefreshToken
	}
	if !token.ExpiresAt.IsZero() {
		session.ExpiresAt = token.ExpiresAt
	}

	s.sessions[id] = session
	err := s.saveLocked()
	s.mu.Unlock()

	return session, true, err
}

func (s *SessionStore) loadFromDisk() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var sessions map[string]Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		return err
	}

	s.mu.Lock()
	s.sessions = sessions
	s.mu.Unlock()
	return nil
}

func (s *SessionStore) saveLocked() error {
	if s.path == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.sessions, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		return err
	}

	return os.Rename(tmpPath, s.path)
}

func randomToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
