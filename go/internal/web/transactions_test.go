package web

import (
	"bytes"
	"eaglebank/internal/accounts"
	"eaglebank/internal/accounts/adapters"
	"eaglebank/internal/transactions"
	adapters2 "eaglebank/internal/transactions/adapters"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	acctStore := adapters.NewInMemoryAccountStore()
	acctSvc := accounts.NewAccountService(acctStore)
	tanStore := adapters2.NewInMemoryTransactionStore()
	tanSvc := transactions.NewTransactionService(tanStore, acctStore)
	srv := NewServer(ServerArgs{Logger: logger, TanSvc: tanSvc, AcctSvc: acctSvc})

	token := login(t, srv, "usr-testuser")

	t.Run("POST to /v1/accounts/{accountId}/transactions", func(t *testing.T) {
		validAcct := createAccount(t, token, srv)
		t.Run("valid deposit request should 201", func(t *testing.T) {
			rr := httptest.NewRecorder()

			ref := "valid request"
			reqObj := CreateTransactionRequest{
				Amount:    100,
				Currency:  accounts.GBP.String(),
				Type:      transactions.Deposit.String(),
				Reference: &ref,
			}
			req := createTransactionRequest(t, reqObj, validAcct.AccountNumber, token)
			srv.ServeHTTP(rr, req)

			var resp TransactionResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rr.Code)
			assert.Equal(t, reqObj.Amount, resp.Amount)
			assert.Equal(t, reqObj.Type, resp.Type)
			assert.Equal(t, reqObj.Currency, resp.Currency)
			assert.Equal(t, *reqObj.Reference, *resp.Reference)
		})
		t.Run("valid withdrawal request should 201", func(t *testing.T) {
			rr := httptest.NewRecorder()

			ref := "valid request"
			reqObj := CreateTransactionRequest{
				Amount:    100,
				Currency:  accounts.GBP.String(),
				Type:      transactions.Withdrawal.String(),
				Reference: &ref,
			}
			req := createTransactionRequest(t, reqObj, validAcct.AccountNumber, token)
			srv.ServeHTTP(rr, req)

			var resp TransactionResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rr.Code)
			assert.Equal(t, reqObj.Amount, resp.Amount)
			assert.Equal(t, reqObj.Type, resp.Type)
			assert.Equal(t, reqObj.Currency, resp.Currency)
			assert.Equal(t, *reqObj.Reference, *resp.Reference)
		})
		t.Run("with invalid data should 400", func(t *testing.T) {
			rr := httptest.NewRecorder()

			ref := "invalid request"
			reqObj := CreateTransactionRequest{
				Amount:    100,
				Currency:  accounts.GBP.String(),
				Type:      "invalid type",
				Reference: &ref,
			}
			req := createTransactionRequest(t, reqObj, validAcct.AccountNumber, token)
			srv.ServeHTTP(rr, req)

			var resp BadRequestErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
		t.Run("without authentication should 401", func(t *testing.T) {
			rr := httptest.NewRecorder()

			ref := "valid request"
			reqObj := CreateTransactionRequest{
				Amount:    100,
				Currency:  accounts.GBP.String(),
				Type:      transactions.Deposit.String(),
				Reference: &ref,
			}
			req := createTransactionRequest(t, reqObj, validAcct.AccountNumber)
			srv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
		t.Run("account not found should 404", func(t *testing.T) {
			rr := httptest.NewRecorder()

			ref := "valid request"
			reqObj := CreateTransactionRequest{
				Amount:    100,
				Currency:  accounts.GBP.String(),
				Type:      transactions.Deposit.String(),
				Reference: &ref,
			}
			req := createTransactionRequest(t, reqObj, "01111111")
			srv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
		t.Run("insufficient funds should 422", func(t *testing.T) {
			rr := httptest.NewRecorder()

			ref := "valid request"
			reqObj := CreateTransactionRequest{
				Amount:    100,
				Currency:  accounts.GBP.String(),
				Type:      transactions.Withdrawal.String(),
				Reference: &ref,
			}
			req := createTransactionRequest(t, reqObj, validAcct.AccountNumber, token)
			srv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		})
		t.Run("unexpected error should 500", func(t *testing.T) {
			errTanSvc := newErroringTransactionService(t)
			errSrv := NewServer(ServerArgs{Logger: logger, TanSvc: errTanSvc, AcctSvc: acctSvc})

			rr := httptest.NewRecorder()

			ref := "valid request"
			reqObj := CreateTransactionRequest{
				Amount:    100,
				Currency:  accounts.GBP.String(),
				Type:      transactions.Withdrawal.String(),
				Reference: &ref,
			}
			req := createTransactionRequest(t, reqObj, validAcct.AccountNumber, token)
			errSrv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusInternalServerError, rr.Code)
		})
	})
}

func createTransactionRequest(t *testing.T, reqObj CreateTransactionRequest, acctNum string, token ...string) *http.Request {
	t.Helper()
	by, err := json.Marshal(reqObj)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/v1/accounts/"+acctNum+"/transactions", bytes.NewBuffer(by))
	if len(token) != 0 {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}
	return req
}

func createAccount(t *testing.T, token string, srv http.Handler) BankAccountResponse {
	t.Helper()

	acctRR := httptest.NewRecorder()
	req := createAccountRequest(t, CreateBankAccountRequest{
		Name:        "Mr Foo's Account",
		AccountType: accounts.PersonalAcct.String(),
	}, token)
	srv.ServeHTTP(acctRR, req)

	var resp BankAccountResponse
	err := json.NewDecoder(acctRR.Body).Decode(&resp)
	require.NoError(t, err)
	return resp
}

type erroringTransactionService struct{}

func (e erroringTransactionService) CreateTransaction(req transactions.CreateTransactionRequest) (transactions.Transaction, error) {
	return transactions.Transaction{}, errors.New("some error")
}

func newErroringTransactionService(t *testing.T) erroringTransactionService {
	t.Helper()
	return erroringTransactionService{}
}
