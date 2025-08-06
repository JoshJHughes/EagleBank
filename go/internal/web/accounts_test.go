package web

import (
	"bytes"
	"eaglebank/internal/accounts"
	"eaglebank/internal/accounts/adapters"
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

func TestAccounts(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	acctStore := adapters.NewInMemoryAccountStore()
	acctSvc := accounts.NewAccountService(acctStore)
	srv := NewServer(ServerArgs{Logger: logger, AcctSvc: acctSvc})

	token := login(t, srv, "usr-testuser")

	t.Run("POST to /v1/accounts", func(t *testing.T) {
		t.Run("with required data should 201", func(t *testing.T) {
			rr := httptest.NewRecorder()
			reqObj := CreateBankAccountRequest{
				Name:        "Mr Foo",
				AccountType: accounts.PersonalAcct.String(),
			}
			req := createAccountRequest(t, reqObj, token)
			srv.ServeHTTP(rr, req)

			var resp BankAccountResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusCreated, rr.Code)
			assert.Equal(t, reqObj.Name, resp.Name)
			assert.Equal(t, reqObj.AccountType, resp.AccountType)
		})
		t.Run("with invalid data should 400", func(t *testing.T) {
			rr := httptest.NewRecorder()
			reqObj := CreateBankAccountRequest{
				Name:        "Mr Foo",
				AccountType: "invalid-account-type",
			}
			req := createAccountRequest(t, reqObj, token)
			srv.ServeHTTP(rr, req)

			var resp BadRequestErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
		t.Run("without authentication should 401", func(t *testing.T) {
			rr := httptest.NewRecorder()
			reqObj := CreateBankAccountRequest{
				Name:        "Mr Foo",
				AccountType: accounts.PersonalAcct.String(),
			}
			req := createAccountRequest(t, reqObj)
			srv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
		t.Run("unexpected error should 500", func(t *testing.T) {
			errAcctSvc := newErroringAccountService(t)
			errSrv := NewServer(ServerArgs{Logger: logger, AcctSvc: errAcctSvc})

			rr := httptest.NewRecorder()
			reqObj := CreateBankAccountRequest{
				Name:        "Mr Foo",
				AccountType: accounts.PersonalAcct.String(),
			}
			req := createAccountRequest(t, reqObj, token)
			errSrv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusInternalServerError, rr.Code)
		})
	})
}

func createAccountRequest(t *testing.T, reqObj CreateBankAccountRequest, token ...string) *http.Request {
	t.Helper()
	by, err := json.Marshal(reqObj)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewBuffer(by))
	if len(token) != 0 {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}
	return req
}

type erroringAccountService struct{}

func (e erroringAccountService) CreateAccount(req accounts.CreateAccountRequest) (*accounts.BankAccount, error) {
	return nil, errors.New("some error")
}

func newErroringAccountService(t *testing.T) erroringAccountService {
	t.Helper()
	return erroringAccountService{}
}
