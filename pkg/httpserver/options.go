package httpserver

import (
	"fmt"
	"log"
	"net"
	"time"
)

// Option -.
type Option func(*Server)

// WithAddress -.
func WithAddress(host string, port int) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort(host, fmt.Sprint(port))
	}
}

// WithReadHeaderTimeout -.
func WithReadHeaderTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readHeaderTimeout = timeout
	}
}

// WithIdleTimeout -.
func WithIdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.idleTimeout = timeout
	}
}

// WithReadTimeout -.
func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WithMaxReadFrameSize -.
func WithMaxReadFrameSize(maxReadFrameSize uint32) Option {
	return func(s *Server) {
		s.maxReadFrameSize = maxReadFrameSize
	}
}

// WithMaxConcurrentStreams -.
func WithMaxConcurrentStreams(maxConcurrentStreams uint32) Option {
	return func(s *Server) {
		s.maxConcurrentStreams = maxConcurrentStreams
	}
}

// WithShutdownTimeout -.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

// WithErrorLog -.
func WithErrorLog(errorLog *log.Logger) Option {
	return func(s *Server) {
		s.errorLog = errorLog
	}
}

// OnBeforeStart -.
func OnBeforeStart(fn func(addr string)) Option {
	return func(s *Server) {
		s.OnBeforeStart = fn
	}
}

// OnBeforeShutdown -.
func OnBeforeShutdown(fn func(addr string)) Option {
	return func(s *Server) {
		s.OnBeforeShutdown = fn
	}
}

// OnAfterShutdown -.
func OnAfterShutdown(fn func(addr string)) Option {
	return func(s *Server) {
		s.OnAfterShutdown = fn
	}
}

// OnError -.
func OnError(fn func(err error)) Option {
	return func(s *Server) {
		s.OnError = fn
	}
}
