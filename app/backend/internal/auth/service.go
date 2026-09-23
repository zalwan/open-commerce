// Package auth is a mock authenticator for v0.1. NOT production security.
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// Role distinguishes admin operators from shoppers.
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
)

// Session is an opaque token. Opaque + in-memory by design for the stub.
type Session struct {
	Token     string
	Email     string
	Role      Role
	ExpiresAt time.Time
}

var ErrUnauthorized = errors.New("unauthorized")

// Service holds users (hardcoded dummy) and tokens.
type Service struct {
	mu     sync.Mutex
	tokens map[string]Session
}

func NewService() *Service { return &Service{tokens: make(map[string]Session)} }

// Login accepts dummy credentials only:
// admin@shop.test / admin123 -> admin, customer@shop.test / customer123 -> customer.
func (s *Service) Login(email, password string) (Session, error) {
	var role Role
	switch {
	case email == "admin@shop.test" && password == "admin123":
		role = RoleAdmin
	case email == "customer@shop.test" && password == "customer123":
		role = RoleCustomer
	default:
		return Session{}, ErrUnauthorized
	}
	tok := randomToken()
	sess := Session{Token: tok, Email: email, Role: role, ExpiresAt: time.Now().Add(24 * time.Hour)}
	s.mu.Lock()
	s.tokens[tok] = sess
	s.mu.Unlock()
	return sess, nil
}

func (s *Service) Authenticate(token string) (Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.tokens[token]
	if !ok {
		return Session{}, ErrUnauthorized
	}
	if time.Now().After(sess.ExpiresAt) {
		delete(s.tokens, token)
		return Session{}, ErrUnauthorized
	}
	return sess, nil
}

func randomToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
