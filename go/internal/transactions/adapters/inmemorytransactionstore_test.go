package adapters

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/transactions"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemoryTransactionStore(t *testing.T) {
	store := NewInMemoryTransactionStore()

	t.Run("should error getting transaction which does not exist", func(t *testing.T) {
		missingID := transactions.TransactionID("0100000")
		_, err := store.GetByTransactionID(missingID)
		assert.Error(t, err)
	})
	t.Run("should perform put-get without errors and fail to update", func(t *testing.T) {
		tan1 := newTestTransaction(t, transactions.Deposit, 150)
		tan2 := newTestTransaction(t, transactions.Deposit, 200)
		t.Run("should create transaction that does not exist in store", func(t *testing.T) {
			err := store.Put(tan1)
			require.NoError(t, err)
		})
		t.Run("should create second transaction for same user", func(t *testing.T) {
			err := store.Put(tan2)
			require.NoError(t, err)
		})
		t.Run("should get an existing transaction by ID", func(t *testing.T) {
			gotTan, err := store.GetByTransactionID(tan1.ID)
			require.NoError(t, err)
			require.Equal(t, tan1, gotTan)
		})
		t.Run("should get both transactions by acctNum", func(t *testing.T) {
			gotTans, err := store.GetByAccountNumber(tan1.AccountNumber)
			require.NoError(t, err)
			require.Len(t, gotTans, 2)
			require.Contains(t, gotTans, tan1)
			require.Contains(t, gotTans, tan2)
		})
		t.Run("should fail updating existing transaction", func(t *testing.T) {
			updatedTan := tan1
			updatedTan.Amount = 9000
			require.NotEqual(t, tan1.Amount, updatedTan.Amount)

			err := store.Put(updatedTan)
			require.Error(t, err)
		})
	})

}

func newTestTransaction(t *testing.T, tanType transactions.TransactionType, amt float64) transactions.Transaction {
	t.Helper()

	tanID, err := transactions.NewRandTransactionID()
	require.NoError(t, err)
	now := time.Now()

	return transactions.Transaction{
		ID:               tanID,
		AccountNumber:    "01000000",
		UserID:           "usr-123",
		Amount:           amt,
		Currency:         accounts.GBP,
		Type:             tanType,
		Reference:        "foo",
		CreatedTimestamp: now,
	}
}
