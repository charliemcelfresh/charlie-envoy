package internal

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_healthHandler_integration(t *testing.T) {
	handler := slog.NewJSONHandler(
		io.Discard, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	testLogger := slog.New(handler).With(
		slog.String("service-name", "envoy-charlie"),
		slog.String("environment", "test"),
	)

	a := NewApp(testLogger, "localhost:8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.healthHandler)

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			resp.StatusCode,
		)
	}

	var got map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got["status"] != "ok" {
		t.Fatalf("expected ok, got %q", got["status"])
	}
}
