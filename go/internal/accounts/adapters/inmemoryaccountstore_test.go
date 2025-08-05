package adapters

import (
	"eaglebank/internal/accounts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewInMemoryAccountStore(t *testing.T) {
	store := NewInMemoryAccountStore()

	t.Run("should error getting account which does not exist", func(t *testing.T) {
		missingID := accounts.AccountNumber("0100000")
		_, err := store.Get(missingID)
		assert.Error(t, err)
	})
	t.Run("should not error deleting account which does not exist", func(t *testing.T) {
		missingID := accounts.AccountNumber("0100000")
		err := store.Delete(missingID)
		assert.NoError(t, err)
	})
	t.Run("should perform put-get-update-delete cycle without errors", func(t *testing.T) {
		acct := newTestAccount(t)
		t.Run("should create account that does not exist in store", func(t *testing.T) {
			err := store.Put(acct)
			require.NoError(t, err)
		})
		t.Run("should get an existing account", func(t *testing.T) {
			gotAcct, err := store.Get(acct.AccountNumber)
			require.NoError(t, err)
			require.Equal(t, acct, gotAcct)
		})
		t.Run("should update existing account", func(t *testing.T) {
			acctVal := *acct
			updatedAcct := &(acctVal)
			updatedAcct.Name = "new name"
			require.NotEqual(t, acct.Name, updatedAcct.Name)

			err := store.Put(updatedAcct)
			require.NoError(t, err)

			gotAcct, err := store.Get(acct.AccountNumber)
			require.NoError(t, err)
			require.Equal(t, updatedAcct, gotAcct)
		})
		t.Run("should delete existing account", func(t *testing.T) {
			err := store.Delete(acct.AccountNumber)
			require.NoError(t, err)

			require.Empty(t, store.store)
		})
	})

}

func newTestAccount(t *testing.T) *accounts.BankAccount {
	t.Helper()

	acctNum, err := accounts.NewRandAccountNumber()
	require.NoError(t, err)
	now := time.Now()

	return &accounts.BankAccount{
		AccountNumber:    acctNum,
		SortCode:         "01-01-01",
		Name:             "Mr Foo",
		AccountType:      accounts.PersonalAcct,
		Currency:         accounts.GBP,
		CreatedTimestamp: now,
		UpdatedTimestamp: now,
	}
}
