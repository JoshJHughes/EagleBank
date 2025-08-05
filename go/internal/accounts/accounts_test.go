package accounts_test

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/accounts/adapters"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccountService(t *testing.T) {
	store := adapters.NewInMemoryAccountStore()
	svc := accounts.NewAccountService(store)

	t.Run("create account", func(t *testing.T) {
		t.Run("should successfully create account", func(t *testing.T) {
			req := accounts.CreateAccountRequest{
				Name:        "Mr Foo",
				AccountType: accounts.PersonalAcct,
			}
			acct, err := svc.CreateAccount(req)
			require.NoError(t, err)

			retAcct, err := store.Get(acct.AccountNumber)
			require.NoError(t, err)
			assert.Equal(t, *acct, *retAcct)
		})
		t.Run("should fail for invalid user", func(t *testing.T) {
			req := accounts.CreateAccountRequest{
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
				Name:        "Mr Foo",
				AccountType: accounts.PersonalAcct,
			}
			_, err := failSvc.CreateAccount(req)
			assert.Error(t, err)
		})
	})
}

type failingAccountStore struct{}

func (f failingAccountStore) Get(acctNum accounts.AccountNumber) (*accounts.BankAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (f failingAccountStore) Put(acct *accounts.BankAccount) error {
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
