package httpserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultAddr                 = ":8080"
	_defaultReadHeaderTimeout    = 5 * time.Second  // Slowloris攻撃対策 (5秒以内にリクエストヘッダーを送信)
	_defaultIdleTimeout          = 60 * time.Second // ALBのタイムアウト時間を考慮 (60秒以内にリクエストがない場合)
	_defaultReadTimeout          = 10 * time.Second // リクエストの読み込み時間 (10秒以内にリクエストを読み込む)
	_defaultMaxConcurrentStreams = 250              // リソース保護と性能のバランスを考慮 (250 concurrent streams)
	_defaultMaxReadFrameSize     = 16 * 1024        // HTTP/2のフレームサイズの最大値 (16KB)
	_defaultShutdownTimeout      = 30 * time.Second // シャットダウン時間 (30秒以内にシャットダウン)
)

// Server -.
type Server struct {
	engine               http.Handler
	address              string
	readHeaderTimeout    time.Duration
	idleTimeout          time.Duration
	readTimeout          time.Duration
	maxConcurrentStreams uint32
	maxReadFrameSize     uint32
	shutdownTimeout      time.Duration
	errorLog             *log.Logger

	// Hooks
	OnBeforeStart    func(addr string)
	OnBeforeShutdown func(addr string)
	OnAfterShutdown  func(addr string)
	OnError          func(err error)
}

// New -.
func New(engine http.Handler, opts ...Option) *Server {
	// setup server
	s := Server{
		engine:               engine,
		address:              _defaultAddr,
		readHeaderTimeout:    _defaultReadHeaderTimeout,
		idleTimeout:          _defaultIdleTimeout,
		readTimeout:          _defaultReadTimeout,
		maxConcurrentStreams: _defaultMaxConcurrentStreams,
		maxReadFrameSize:     _defaultMaxReadFrameSize,
		shutdownTimeout:      _defaultShutdownTimeout,
	}

	// bind options
	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

// ListenAndServe -.
func (s *Server) ListenAndServe() error {
	defer func() {
		if s.OnAfterShutdown != nil {
			s.OnAfterShutdown(s.address)
		}
	}()

	// setup protocols
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true) // Use h2c so we can serve HTTP/2 without TLS.

	// setup server
	server := http.Server{
		Addr: s.address,
		Handler: h2c.NewHandler(s.engine, &http2.Server{
			MaxConcurrentStreams: s.maxConcurrentStreams,
			MaxReadFrameSize:     s.maxReadFrameSize,
		}),
		Protocols:         protocols,
		ReadHeaderTimeout: s.readHeaderTimeout,
		ErrorLog:          s.errorLog,
	}
	h2c.NewHandler(s.engine, &http2.Server{})

	sigCtx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,    // ctrl+c
		syscall.SIGTERM, // kill command
		syscall.SIGHUP,  // hung up
	)
	defer stop()

	g, ctx := errgroup.WithContext(sigCtx)
	g.SetLimit(2)

	g.Go(func() error {
		if s.OnBeforeStart != nil {
			s.OnBeforeStart(s.address)
		}

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		if s.OnBeforeShutdown != nil {
			s.OnBeforeShutdown(s.address)
		}

		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		if s.OnError != nil {
			s.OnError(err)
		}
		return err
	}
	return nil
}
