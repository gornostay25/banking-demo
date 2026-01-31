//go:build integration

package services_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"banking/internal/testhelpers"

	"github.com/stretchr/testify/require"
)

func TestAccountsListSuccess(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	req := httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
	attachAuthCookie(req, token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var payload []map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&payload))
	require.Len(t, payload, 2)
}

func TestAccountsListUnauthorized(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/accounts", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAccountBalanceSuccess(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)
	db := testhelpers.RequireTestDB(t)

	userID := testhelpers.GetUserIDByEmail(t, db, "user1@test.com")
	accountID := testhelpers.GetAccountID(t, db, userID, "USD")

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	req := httptest.NewRequest(http.MethodGet, "/api/accounts/"+accountID.String()+"/balance", nil)
	attachAuthCookie(req, token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var payload map[string]any
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&payload))
	require.Equal(t, accountID.String(), payload["id"])
}

func TestAccountBalanceNotFound(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	req := httptest.NewRequest(http.MethodGet, "/api/accounts/550e8400-e29b-41d4-a716-446655440000/balance", nil)
	attachAuthCookie(req, token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAccountBalanceOtherUserNotFound(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)
	db := testhelpers.RequireTestDB(t)

	otherUserID := testhelpers.GetUserIDByEmail(t, db, "user2@test.com")
	otherAccountID := testhelpers.GetAccountID(t, db, otherUserID, "USD")

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	req := httptest.NewRequest(http.MethodGet, "/api/accounts/"+otherAccountID.String()+"/balance", nil)
	attachAuthCookie(req, token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAccountBalanceUnauthorized(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/accounts/550e8400-e29b-41d4-a716-446655440000/balance", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
