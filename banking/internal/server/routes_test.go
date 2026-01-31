//go:build integration

package server_test

import (
	"banking/internal/testhelpers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["status"] != "up" {
		t.Fatalf("expected status up, got %s", payload["status"])
	}
}
