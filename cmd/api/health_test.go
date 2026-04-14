package main

import (
	"net/http"
	"testing"

	"github.com/felipeeguia03/vol7/internal/ratelimiter"
)

func TestHealthHandler(t *testing.T) {
	cfg := config{
		rateLimiter: ratelimiter.Config{
			Enabled:             false,
			RequestsPerTimeFrame: 20,
			TimeFrame:           0,
		},
	}
	app := newTestApplication(t, cfg)
	mux := app.mount()

	req, _ := http.NewRequest(http.MethodGet, "/v1/health", nil)
	rr := executeRequest(req, mux)

	checkResponseCode(t, http.StatusOK, rr.Code)
}
