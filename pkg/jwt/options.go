package jwt

import (
	"time"
)

// Option -.
type Option func(*Manager)

// WithNowFunc -.
func WithNowFunc(nowFunc func() time.Time) Option {
	return func(m *Manager) {
		m.now = nowFunc
	}
}

// WithTokenExpiry -.
func WithTokenExpiry(tokenExpiry time.Duration) Option {
	return func(m *Manager) {
		m.tokenExpiry = tokenExpiry
	}
}
