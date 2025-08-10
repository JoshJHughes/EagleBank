package transactions_test

import (
	"eaglebank/internal/accounts"
	adapters2 "eaglebank/internal/accounts/adapters"
	"eaglebank/internal/transactions"
	"eaglebank/internal/transactions/adapters"
	"eaglebank/internal/users"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	acctStore := adapters2.NewInMemoryAccountStore()
	acctSvc := accounts.NewAccountService(acctStore)

	tanStore := adapters.NewInMemoryTransactionStore()
	tanSvc := transactions.NewTransactionService(tanStore, acctStore)

	userID := users.MustNewUserID("usr-123")
	acct, err := acctSvc.CreateAccount(accounts.CreateAccountRequest{
		UserID:      userID,
		Name:        "Mr Foo",
		AccountType: accounts.PersonalAcct,
	})
	require.NoError(t, err)
	t.Run("should successfully create transaction and update balance", func(t *testing.T) {
		amt := 10.0
		tan, err := tanSvc.CreateTransaction(transactions.CreateTransactionRequest{
			AccountNumber: acct.AccountNumber,
			UserID:        userID,
			Amount:        amt,
			Currency:      accounts.GBP,
			Type:          transactions.Deposit,
			Reference:     "valid",
		})
		require.NoError(t, err)

		gotTan, err := tanStore.GetByTransactionID(tan.ID)
		require.NoError(t, err)
		assert.Equal(t, tan, gotTan)

		gotAcct, err := acctSvc.FetchAccount(acct.AccountNumber)
		require.NoError(t, err)
		assert.Equal(t, acct.Balance()+amt, gotAcct.Balance())
	})
	t.Run("should fail for invalid transaction", func(t *testing.T) {
		_, err := tanSvc.CreateTransaction(transactions.CreateTransactionRequest{
			AccountNumber: acct.AccountNumber,
			UserID:        userID,
			Amount:        -10,
			Currency:      accounts.GBP,
			Type:          transactions.Deposit,
			Reference:     "invalid",
		})
		assert.Error(t, err)
	})
	t.Run("should fail if account doesn't exist", func(t *testing.T) {
		_, err := tanSvc.CreateTransaction(transactions.CreateTransactionRequest{
			AccountNumber: "01000000",
			UserID:        userID,
			Amount:        10,
			Currency:      accounts.GBP,
			Type:          transactions.Deposit,
			Reference:     "acct missing",
		})
		assert.ErrorIs(t, err, accounts.ErrAccountNotFound)
	})
	t.Run("should fail withdraw with insufficient funds", func(t *testing.T) {
		preWithdrawAcct, err := acctSvc.FetchAccount(acct.AccountNumber)
		require.NoError(t, err)

		_, err = tanSvc.CreateTransaction(transactions.CreateTransactionRequest{
			AccountNumber: preWithdrawAcct.AccountNumber,
			UserID:        preWithdrawAcct.UserID,
			Amount:        preWithdrawAcct.Balance() * 2,
			Currency:      accounts.GBP,
			Type:          transactions.Withdrawal,
			Reference:     "overdrawn",
		})
		assert.ErrorIs(t, err, accounts.ErrInsufficientFunds)

		postWithdrawAcct, err := acctSvc.FetchAccount(acct.AccountNumber)
		require.NoError(t, err)
		assert.Equal(t, preWithdrawAcct, postWithdrawAcct)
	})
	t.Run("should fail deposit if account above limit", func(t *testing.T) {
		preDepositAcct, err := acctSvc.FetchAccount(acct.AccountNumber)
		require.NoError(t, err)
		require.Greater(t, preDepositAcct.Balance()+accounts.BalanceMax, accounts.BalanceMax)

		_, err = tanSvc.CreateTransaction(transactions.CreateTransactionRequest{
			AccountNumber: preDepositAcct.AccountNumber,
			UserID:        preDepositAcct.UserID,
			Amount:        accounts.BalanceMax,
			Currency:      accounts.GBP,
			Type:          transactions.Deposit,
			Reference:     "too much money",
		})
		assert.ErrorIs(t, err, accounts.ErrTooManyFunds)

		postDepositAcct, err := acctSvc.FetchAccount(acct.AccountNumber)
		require.NoError(t, err)
		assert.Equal(t, preDepositAcct, postDepositAcct)
	})
}

func TestListTransaction(t *testing.T) {
	acctStore := adapters2.NewInMemoryAccountStore()
	acctSvc := accounts.NewAccountService(acctStore)

	tanStore := adapters.NewInMemoryTransactionStore()
	tanSvc := transactions.NewTransactionService(tanStore, acctStore)

	userID := users.MustNewUserID("usr-123")
	acct, err := acctSvc.CreateAccount(accounts.CreateAccountRequest{
		UserID:      userID,
		Name:        "Mr Foo",
		AccountType: accounts.PersonalAcct,
	})
	require.NoError(t, err)

	tan1, err := tanSvc.CreateTransaction(transactions.CreateTransactionRequest{
		AccountNumber: acct.AccountNumber,
		UserID:        userID,
		Amount:        100,
		Currency:      accounts.GBP,
		Type:          transactions.Deposit,
	})
	require.NoError(t, err)
	tan2, err := tanSvc.CreateTransaction(transactions.CreateTransactionRequest{
		AccountNumber: acct.AccountNumber,
		UserID:        userID,
		Amount:        100,
		Currency:      accounts.GBP,
		Type:          transactions.Deposit,
	})
	require.NoError(t, err)

	t.Run("should list all accounts", func(t *testing.T) {
		accts, err := tanSvc.ListTransactions(acct.AccountNumber)
		require.NoError(t, err)
		assert.Len(t, accts, 2)
		assert.Contains(t, accts, tan1)
		assert.Contains(t, accts, tan2)
	})
	t.Run("should return empty list if user has no accounts", func(t *testing.T) {
		acctNumNoTans, err := accounts.NewRandAccountNumber()
		require.NoError(t, err)
		accts, err := tanSvc.ListTransactions(acctNumNoTans)
		require.NoError(t, err)
		assert.Empty(t, accts)
	})
}
