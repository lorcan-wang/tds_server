package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"tds_server/internal/config"
)

func TestPartnerTokenServiceCachesAndRefreshes(t *testing.T) {
	var tokenRequestCount int32
	var registerRequestCount int32
	var expectedAudience string

	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode json body: %v", err)
		}
		if got := body["grant_type"]; got != partnerGrantType {
			t.Fatalf("unexpected grant_type: %s", got)
		}
		if got := body["client_id"]; got != "client-id" {
			t.Fatalf("unexpected client_id: %s", got)
		}
		if got := body["client_secret"]; got != "client-secret" {
			t.Fatalf("unexpected client_secret: %s", got)
		}
		if got := body["audience"]; got != expectedAudience {
			t.Fatalf("unexpected audience: %s", got)
		}
		if got := body["scope"]; got != "scope" {
			t.Fatalf("unexpected scope: %s", got)
		}

		current := atomic.AddInt32(&tokenRequestCount, 1)
		token := fmt.Sprintf("token-%d", current)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"access_token":"%s","expires_in":120}`, token)
	}))
	defer tokenServer.Close()

	registerServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/1/partner_accounts" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			t.Fatalf("unexpected content type: %s", contentType)
		}
		expectedAuth := "Bearer token-1"
		if auth := r.Header.Get("Authorization"); auth != expectedAuth {
			t.Fatalf("unexpected authorization header: %s", auth)
		}

		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode register json body: %v", err)
		}
		if got := body["domain"]; got != "domain.com" {
			t.Fatalf("unexpected domain: %s", got)
		}

		if current := atomic.AddInt32(&registerRequestCount, 1); current > 1 {
			t.Fatalf("register endpoint called more than once")
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer registerServer.Close()

	expectedAudience = registerServer.URL

	cfg := &config.Config{
		TeslaClientID:        "client-id",
		TeslaClientSecret:    "client-secret",
		TeslaAPIURL:          registerServer.URL,
		TeslaPartnerTokenURL: tokenServer.URL,
		TeslaPartnerScope:    "scope",
		TeslaPartnerDomain:   "domain.com",
	}

	svc, err := NewPartnerTokenService(cfg)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	if got := atomic.LoadInt32(&tokenRequestCount); got != 1 {
		t.Fatalf("expected 1 request during initialization, got %d", got)
	}
	if got := atomic.LoadInt32(&registerRequestCount); got != 1 {
		t.Fatalf("expected register endpoint to be called once at startup, got %d", got)
	}

	token, err := svc.GetToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error getting token: %v", err)
	}
	if token != "token-1" {
		t.Fatalf("unexpected token: %s", token)
	}

	if got := atomic.LoadInt32(&tokenRequestCount); got != 1 {
		t.Fatalf("expected cached token to be used, total token requests %d", got)
	}
	if got := atomic.LoadInt32(&registerRequestCount); got != 1 {
		t.Fatalf("register endpoint should not be called again, got %d", got)
	}

	svc.mu.Lock()
	svc.expiresAt = time.Now().Add(-time.Minute)
	svc.mu.Unlock()

	refreshed, err := svc.GetToken(context.Background())
	if err != nil {
		t.Fatalf("unexpected error getting refreshed token: %v", err)
	}
	if refreshed != "token-2" {
		t.Fatalf("unexpected refreshed token: %s", refreshed)
	}

	if got := atomic.LoadInt32(&tokenRequestCount); got != 2 {
		t.Fatalf("expected refresh request to be made, total token requests %d", got)
	}
	if got := atomic.LoadInt32(&registerRequestCount); got != 1 {
		t.Fatalf("register endpoint should only be called once, got %d", got)
	}
}
