package internal

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (a *app) healthHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{
		"status": "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		a.logger.LogAttrs(
			r.Context(), slog.LevelError, "healthHandler error",
			slog.String("caller", "internal.healthHandler"),
			slog.Any("error", err),
		)
	}
}
