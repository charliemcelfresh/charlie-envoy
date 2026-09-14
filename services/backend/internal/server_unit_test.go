package internal

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_unit(t *testing.T) {
	testLogger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)

	testApp := NewApp(testLogger, "localhost:8080")

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	testApp.healthHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			resp.StatusCode,
		)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body["status"] != "ok" {
		t.Fatalf(
			"expected status %q, got %q",
			"ok",
			body["status"],
		)
	}
}
