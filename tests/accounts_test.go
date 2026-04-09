package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"pismo-assignment/db"
	"pismo-assignment/models"
)

func TestAccounts_Create(t *testing.T) {
	r := setupRouter(t)

	const documentNumber = "12345678900"
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"`+documentNumber+`"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp models.CreateAccountResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, documentNumber, resp.DocumentNumber)
	require.NotZero(t, resp.AccountID)

	var row models.Account
	require.NoError(t, db.DB.First(&row, resp.AccountID).Error)
	require.Equal(t, documentNumber, row.DocumentNumber)
	require.Equal(t, int64(0), row.BalanceInPaisa)
}
