//go:build integration

package services_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"banking/internal/testhelpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTransactionsTransferSuccess(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)
	db := testhelpers.RequireTestDB(t)

	user2ID := testhelpers.GetUserIDByEmail(t, db, "user2@test.com")
	recipientAccountID := testhelpers.GetAccountID(t, db, user2ID, "USD")

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	body := []byte(`{"to_account_id":"` + recipientAccountID.String() + `","amount":"10.00","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/transfer", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestTransactionsTransferInsufficientFunds(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)
	db := testhelpers.RequireTestDB(t)

	user2ID := testhelpers.GetUserIDByEmail(t, db, "user2@test.com")
	recipientAccountID := testhelpers.GetAccountID(t, db, user2ID, "USD")

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	body := []byte(`{"to_account_id":"` + recipientAccountID.String() + `","amount":"1000000.00","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/transfer", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTransactionsTransferRecipientNotFound(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	body := []byte(`{"to_account_id":"` + uuid.New().String() + `","amount":"10.00","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/transfer", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTransactionsTransferSelfTransfer(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)
	db := testhelpers.RequireTestDB(t)

	user1ID := testhelpers.GetUserIDByEmail(t, db, "user1@test.com")
	selfAccountID := testhelpers.GetAccountID(t, db, user1ID, "USD")

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	body := []byte(`{"to_account_id":"` + selfAccountID.String() + `","amount":"10.00","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/transfer", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTransactionsTransferCurrencyMismatch(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)
	db := testhelpers.RequireTestDB(t)

	user2ID := testhelpers.GetUserIDByEmail(t, db, "user2@test.com")
	recipientAccountID := testhelpers.GetAccountID(t, db, user2ID, "EUR")

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	body := []byte(`{"to_account_id":"` + recipientAccountID.String() + `","amount":"10.00","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/transfer", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTransactionsTransferUnauthorized(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	body := []byte(`{"to_account_id":"` + uuid.New().String() + `","amount":"10.00","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/transfer", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTransactionsExchangeSuccessUSDToEUR(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	body := []byte(`{"from_currency":"USD","to_currency":"EUR","amount":"10.00"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/exchange", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestTransactionsExchangeSuccessEURToUSD(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	body := []byte(`{"from_currency":"EUR","to_currency":"USD","amount":"10.00"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/exchange", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestTransactionsExchangeInsufficientFunds(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	body := []byte(`{"from_currency":"USD","to_currency":"EUR","amount":"1000000.00"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/exchange", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTransactionsExchangeSameCurrency(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")
	body := []byte(`{"from_currency":"USD","to_currency":"USD","amount":"10.00"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/exchange", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTransactionsExchangeUnauthorized(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	body := []byte(`{"from_currency":"USD","to_currency":"EUR","amount":"10.00"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/exchange", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTransactionsListSuccess(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)
	db := testhelpers.RequireTestDB(t)

	user2ID := testhelpers.GetUserIDByEmail(t, db, "user2@test.com")
	recipientAccountID := testhelpers.GetAccountID(t, db, user2ID, "USD")

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")

	transferBody := []byte(`{"to_account_id":"` + recipientAccountID.String() + `","amount":"10.00","currency":"USD"}`)
	transferReq := httptest.NewRequest(http.MethodPost, "/api/transactions/transfer", bytes.NewReader(transferBody))
	attachAuthCookie(transferReq, token)
	transferReq.Header.Set("Content-Type", "application/json")
	transferRec := httptest.NewRecorder()
	handler.ServeHTTP(transferRec, transferReq)
	require.Equal(t, http.StatusOK, transferRec.Code)

	exchangeBody := []byte(`{"from_currency":"USD","to_currency":"EUR","amount":"10.00"}`)
	exchangeReq := httptest.NewRequest(http.MethodPost, "/api/transactions/exchange", bytes.NewReader(exchangeBody))
	attachAuthCookie(exchangeReq, token)
	exchangeReq.Header.Set("Content-Type", "application/json")
	exchangeRec := httptest.NewRecorder()
	handler.ServeHTTP(exchangeRec, exchangeReq)
	require.Equal(t, http.StatusOK, exchangeRec.Code)

	listReq := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	attachAuthCookie(listReq, token)
	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, listReq)

	require.Equal(t, http.StatusOK, listRec.Code)
	var payload struct {
		Transactions []map[string]any `json:"transactions"`
	}
	require.NoError(t, json.NewDecoder(listRec.Body).Decode(&payload))
	require.NotEmpty(t, payload.Transactions)
}

func TestTransactionsListFilterAndPagination(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)
	db := testhelpers.RequireTestDB(t)

	user2ID := testhelpers.GetUserIDByEmail(t, db, "user2@test.com")
	recipientAccountID := testhelpers.GetAccountID(t, db, user2ID, "USD")

	token, _ := loginAndGetTokens(t, handler, "user1@test.com", "password")

	body := []byte(`{"to_account_id":"` + recipientAccountID.String() + `","amount":"10.00","currency":"USD"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/transactions/transfer", bytes.NewReader(body))
	attachAuthCookie(req, token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	listReq := httptest.NewRequest(http.MethodGet, "/api/transactions?type=transfer&limit=1&page=1", nil)
	attachAuthCookie(listReq, token)
	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, listReq)

	require.Equal(t, http.StatusOK, listRec.Code)
	var payload struct {
		Transactions []map[string]any `json:"transactions"`
		Limit        int              `json:"limit"`
		Page         int              `json:"page"`
	}
	require.NoError(t, json.NewDecoder(listRec.Body).Decode(&payload))
	require.Equal(t, 1, payload.Limit)
	require.Equal(t, 1, payload.Page)
	for _, item := range payload.Transactions {
		require.Equal(t, "transfer", item["type"])
	}
}

func TestTransactionsListUnauthorized(t *testing.T) {
	testhelpers.ResetTestData(t)
	handler := testhelpers.SetupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
