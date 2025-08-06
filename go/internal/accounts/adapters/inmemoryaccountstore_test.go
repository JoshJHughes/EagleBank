package adapters

import (
	"eaglebank/internal/accounts"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemoryAccountStore(t *testing.T) {
	store := NewInMemoryAccountStore()

	t.Run("should error getting account which does not exist", func(t *testing.T) {
		missingID := accounts.AccountNumber("0100000")
		_, err := store.GetByAcctNum(missingID)
		assert.Error(t, err)
	})
	t.Run("should error not found deleting account which does not exist", func(t *testing.T) {
		missingID := accounts.AccountNumber("0100000")
		err := store.Delete(missingID)
		assert.ErrorIs(t, err, accounts.ErrAccountNotFound)
	})
	t.Run("should perform put-get-update-delete cycle without errors", func(t *testing.T) {
		acct1 := newTestAccount(t)
		acct2 := newTestAccount(t)
		t.Run("should create account that does not exist in store", func(t *testing.T) {
			err := store.Put(acct1)
			require.NoError(t, err)
		})
		t.Run("should create second account for same user", func(t *testing.T) {
			err := store.Put(acct2)
			require.NoError(t, err)
		})
		t.Run("should get an existing account by acct num", func(t *testing.T) {
			gotAcct, err := store.GetByAcctNum(acct1.AccountNumber)
			require.NoError(t, err)
			require.Equal(t, acct1, gotAcct)
		})
		t.Run("should get both accounts by userID", func(t *testing.T) {
			gotAccts, err := store.GetByUserID(acct1.UserID)
			require.NoError(t, err)
			require.Len(t, gotAccts, 2)
		})
		t.Run("should update existing account", func(t *testing.T) {
			updatedAcct := acct1
			updatedAcct.Name = "new name"
			require.NotEqual(t, acct1.Name, updatedAcct.Name)

			err := store.Put(updatedAcct)
			require.NoError(t, err)

			gotAcct, err := store.GetByAcctNum(acct1.AccountNumber)
			require.NoError(t, err)
			require.Equal(t, updatedAcct, gotAcct)

			gotAccts, err := store.GetByUserID(acct1.UserID)
			require.NoError(t, err)
			require.Len(t, gotAccts, 2)
			require.Contains(t, gotAccts, updatedAcct)
		})
		t.Run("should delete existing account", func(t *testing.T) {
			err := store.Delete(acct1.AccountNumber)
			require.NoError(t, err)

			require.NotContains(t, store.acctsByNumber, acct1.AccountNumber)
			require.Contains(t, store.acctsByNumber, acct2.AccountNumber)
			accts := store.acctsByUserID[acct1.UserID]
			require.Len(t, accts, 1)
			require.NotContains(t, accts, acct1)
			require.Contains(t, accts, acct2)
		})
	})

}

func newTestAccount(t *testing.T) accounts.BankAccount {
	t.Helper()

	acctNum, err := accounts.NewRandAccountNumber()
	require.NoError(t, err)
	now := time.Now()

	return accounts.BankAccount{
		UserID:           "usr-123",
		AccountNumber:    acctNum,
		SortCode:         "01-01-01",
		Name:             "Mr Foo",
		AccountType:      accounts.PersonalAcct,
		Currency:         accounts.GBP,
		CreatedTimestamp: now,
		UpdatedTimestamp: now,
	}
}
