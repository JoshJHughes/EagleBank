package accounts_test

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/accounts/adapters"
	"eaglebank/internal/users"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountService(t *testing.T) {
	store := adapters.NewInMemoryAccountStore()
	svc := accounts.NewAccountService(store)

	t.Run("create account", func(t *testing.T) {
		t.Run("should successfully create account", func(t *testing.T) {
			req := accounts.CreateAccountRequest{
				UserID:      "usr-123",
				Name:        "Mr Foo",
				AccountType: accounts.PersonalAcct,
			}
			acct, err := svc.CreateAccount(req)
			require.NoError(t, err)

			retAcct, err := store.GetByAcctNum(acct.AccountNumber)
			require.NoError(t, err)
			assert.Equal(t, *acct, retAcct)
		})
		t.Run("should fail for invalid user", func(t *testing.T) {
			req := accounts.CreateAccountRequest{
				UserID:      "usr-123",
				Name:        "Mr Foo",
				AccountType: "invalid account type",
			}
			_, err := svc.CreateAccount(req)
			assert.Error(t, err)
		})
		t.Run("should fail if put fails", func(t *testing.T) {
			failStore := newFailingAccountStore(t)
			failSvc := accounts.NewAccountService(failStore)
			req := accounts.CreateAccountRequest{
				UserID:      "usr-123",
				Name:        "Mr Foo",
				AccountType: accounts.PersonalAcct,
			}
			_, err := failSvc.CreateAccount(req)
			assert.Error(t, err)
		})
	})
	t.Run("list accounts", func(t *testing.T) {
		userID := users.MustNewUserID("usr-123")
		acct1, err := svc.CreateAccount(accounts.CreateAccountRequest{
			UserID:      userID,
			Name:        "Mr Foo",
			AccountType: accounts.PersonalAcct,
		})
		require.NoError(t, err)
		acct2, err := svc.CreateAccount(accounts.CreateAccountRequest{
			UserID:      userID,
			Name:        "Mr Foo",
			AccountType: accounts.PersonalAcct,
		})
		require.NoError(t, err)
		t.Run("should list all accounts", func(t *testing.T) {
			accts, err := svc.ListAccounts(userID)
			require.NoError(t, err)
			assert.Len(t, accts, 2)
			assert.Contains(t, accts, *acct1)
			assert.Contains(t, accts, *acct2)
		})
		t.Run("should return empty list if user has no accounts", func(t *testing.T) {
			userIDNoAccs := users.MustNewUserID("usr-1234")
			accts, err := svc.ListAccounts(userIDNoAccs)
			require.NoError(t, err)
			assert.Empty(t, accts)
		})
		t.Run("should error if store errors for other reason", func(t *testing.T) {
			failStore := newFailingAccountStore(t)
			failSvc := accounts.NewAccountService(failStore)
			_, err = failSvc.ListAccounts(userID)
			assert.Error(t, err)
		})
	})
}

type failingAccountStore struct{}

func (f failingAccountStore) GetByUserID(userID users.UserID) ([]accounts.BankAccount, error) {
	return nil, errors.New("some error")
}

func (f failingAccountStore) GetByAcctNum(acctNum accounts.AccountNumber) (accounts.BankAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (f failingAccountStore) Put(acct accounts.BankAccount) error {
	return errors.New("error")
}

func (f failingAccountStore) Delete(acctNum accounts.AccountNumber) error {
	//TODO implement me
	panic("implement me")
}

func newFailingAccountStore(t *testing.T) *failingAccountStore {
	t.Helper()
	return &failingAccountStore{}
}
