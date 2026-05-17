package health

import "time"

type Option func(*health)

func WithStatus(status string) Option {
	return func(h *health) {
		h.response.Status = status
	}
}

func WithApp(app string) Option {
	return func(h *health) {
		h.response.App = app
	}
}

func WithEnv(env string) Option {
	return func(h *health) {
		h.response.Env = env
	}
}

func WithVersion(version string) Option {
	return func(h *health) {
		h.response.Version = version
	}
}

func WithStartedAt(startedAt time.Time) Option {
	return func(h *health) {
		h.response.StartedAt = startedAt
	}
}
