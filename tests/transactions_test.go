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

func TestTransactions_Create(t *testing.T) {
	r := setupRouter(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"document_number":"99999999999"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var accResp models.CreateAccountResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &accResp))
	accountID := accResp.AccountID

	var accRow models.Account
	require.NoError(t, db.DB.First(&accRow, accountID).Error)
	require.Equal(t, int64(0), accRow.BalanceInPaisa)

	tx := map[string]any{"account_id": accountID, "operation_type_id": 4, "amount": 200.0}
	b, err := json.Marshal(tx)
	require.NoError(t, err)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	require.Equal(t, http.StatusCreated, w2.Code)

	var txResp models.CreateTransactionResponse
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &txResp))
	require.Equal(t, accountID, txResp.AccountID)
	require.Equal(t, 4, txResp.OperationTypeID)
	require.InDelta(t, 200.0, txResp.Amount, 0.001)
	require.NotZero(t, txResp.TransactionID)

	var txRow models.Transaction
	require.NoError(t, db.DB.First(&txRow, txResp.TransactionID).Error)
	require.Equal(t, accountID, txRow.AccountId)
	require.Equal(t, models.CreditVoucher, txRow.OperationTypeId)
	require.Equal(t, int64(20_000), txRow.AmountInPaisa)

	require.NoError(t, db.DB.First(&accRow, accountID).Error)
	require.Equal(t, int64(20_000), accRow.BalanceInPaisa)
}
