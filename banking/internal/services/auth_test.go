//go:build integration

package services_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"banking/internal/testhelpers"

	"github.com/stretchr/testify/require"
)

func TestAuthLoginSuccess(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	require.NotEmpty(t, token)
}

func TestAuthLoginInvalidEmail(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	status := loginStatus(t, handler, "not-an-email", "password")
	require.Equal(t, http.StatusUnauthorized, status)
}

func TestAuthLoginInvalidPassword(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	status := loginStatus(t, handler, "user1@test.com", "wrong-password")
	require.Equal(t, http.StatusUnauthorized, status)
}

func TestAuthLoginUserNotFound(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	status := loginStatus(t, handler, "missing@test.com", "password")
	require.Equal(t, http.StatusUnauthorized, status)
}

func TestAuthRefreshSuccess(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	_, refreshToken := loginAndGetTokens(t, handler, "user1@test.com", "password")
	require.NotEmpty(t, refreshToken)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var payload struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&payload))
	token := payload.Token
	if token == "" {
		token, _ = extractTokensFromCookies(rec.Result().Cookies())
	}
	require.NotEmpty(t, token)
}

func TestAuthRefreshInvalidToken(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "bad-token"})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthLogoutSuccess(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	attachAuthCookie(req, token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMeSuccess(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	attachAuthCookie(req, token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var payload struct {
		Email string `json:"email"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&payload))
	require.Equal(t, "user1@test.com", payload.Email)
}

func TestAuthMeUnauthorized(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func loginStatus(t *testing.T, handler http.Handler, email, password string) int {
	t.Helper()
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec.Code
}

func loginAndGetTokens(t *testing.T, handler http.Handler, email, password string) (string, string) {
	t.Helper()
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var payloadResp struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&payloadResp))

	token, refreshToken := extractTokensFromCookies(rec.Result().Cookies())
	if payloadResp.Token != "" {
		token = payloadResp.Token
	}

	return token, refreshToken
}

func extractTokensFromCookies(cookies []*http.Cookie) (string, string) {
	var token string
	var refreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			token = cookie.Value
		}
		if cookie.Name == "refresh_token" {
			refreshToken = cookie.Value
		}
	}
	return token, refreshToken
}

func attachAuthCookie(req *http.Request, token string) {
	if token == "" {
		return
	}
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
}
