package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Manager -.
type Manager struct {
	secret      string
	tokenExpiry time.Duration
	now         func() time.Time
}

// New -.
func New(secret string, opts ...Option) (*Manager, error) {
	m := Manager{
		secret:      secret,
		tokenExpiry: 1 * time.Hour,
		now:         time.Now,
	}
	for _, opt := range opts {
		opt(&m)
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

// Generate -.
func (m *Manager) Generate(id string) (token string, err error) {
	token, err = jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.RegisteredClaims{
		Subject:   id,
		ExpiresAt: jwtlib.NewNumericDate(m.now().Add(m.tokenExpiry)),
	}).SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return token, nil
}

func (m *Manager) validate() error {
	var errs []error
	if m.secret == "" {
		errs = append(errs, errors.New("jwt secret is required"))
	}
	if m.tokenExpiry <= 0 {
		errs = append(errs, errors.New("jwt token expiry must be > 0"))
	}
	if m.now == nil {
		errs = append(errs, errors.New("now func is required"))
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
