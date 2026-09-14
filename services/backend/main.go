package main

import (
	"log/slog"
	"os"

	"github.com/charliemcelfresh/envoy-charlie/services/backend/internal"
)

const serverAddr = "0.0.0.0:8080"

func main() {
	handler := slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	logger := slog.New(handler).With(
		slog.String("service-name", "envoy-charlie"),
		slog.String("environment", "production"),
	)

	slog.SetDefault(logger)

	slog.Info("server started")

	a := internal.NewApp(logger, serverAddr)
	a.Serve()
}
