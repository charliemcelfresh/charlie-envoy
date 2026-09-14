package internal

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type app struct {
	server server
	logger logger
	addr   string
}

type logger interface {
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
}

type server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func NewApp(logger logger, addr string) *app {
	a := &app{
		logger: logger,
		addr:   addr,
	}

	s := &http.Server{
		Addr:              addr,
		Handler:           a.BuildHandler(),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       90 * time.Second,
	}
	a.server = s
	return a
}

func (a *app) BuildHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /health", chain(http.HandlerFunc(a.healthHandler), a.jsonApplicationContentMiddleware))
	return mux
}

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func (a *app) Serve() {
	// make cancellation
	signalCtx, signalCleanup := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalCleanup()

	// idempotent shutdown
	var oneShutdown sync.Once
	shutdown := func() {
		oneShutdown.Do(
			func() {
				timeoutCtx, timeoutFunc := context.WithTimeout(context.Background(), 30*time.Second)
				defer timeoutFunc()
				if err := a.server.Shutdown(timeoutCtx); err != nil {
					a.logger.LogAttrs(
						timeoutCtx, slog.LevelError, "error shutdown",
						slog.String("caller", "internal.app.Serve"),
						slog.Any("error", err),
					)
				}
				a.logger.LogAttrs(
					timeoutCtx, slog.LevelInfo, "graceful shutdown complete",
					slog.String("caller", "internal.app.Serve"),
				)
			},
		)
	}

	var wg sync.WaitGroup
	a.logger.LogAttrs(
		context.Background(), slog.LevelInfo, "starting server",
		slog.String("caller", "internal.app.Serve"),
		slog.String("addr", a.addr),
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.LogAttrs(
				context.Background(), slog.LevelError, "Serve failed",
				slog.String("caller", "internal.app.Serve"),
				slog.Any("error", err),
			)
			shutdown()
		}
	}()

	<-signalCtx.Done()
	a.logger.LogAttrs(
		context.Background(), slog.LevelInfo, "shutting down",
		slog.String("caller", "internal.	app.Serve"),
	)
	shutdown()
	// allow server goroutine to complete
	wg.Wait()
}
