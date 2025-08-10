package transactions

import (
	"eaglebank/internal/accounts"
	"errors"
	"fmt"
)

type TransactionStore interface {
	GetByTransactionID(tanID TransactionID) (Transaction, error)
	GetByAccountNumber(acctNum accounts.AccountNumber) ([]Transaction, error)
	Put(tan Transaction) error
}

type accountStore interface {
	GetByAcctNum(acctNum accounts.AccountNumber) (accounts.BankAccount, error)
	Put(acct accounts.BankAccount) error
}

type TransactionService struct {
	transactionStore TransactionStore
	acctStore        accountStore
}

func NewTransactionService(tanStore TransactionStore, acctStore accountStore) *TransactionService {
	return &TransactionService{transactionStore: tanStore, acctStore: acctStore}
}

func (svc *TransactionService) CreateTransaction(req CreateTransactionRequest) (Transaction, error) {
	acct, err := svc.acctStore.GetByAcctNum(req.AccountNumber)
	if err != nil {
		if errors.Is(err, accounts.ErrAccountNotFound) {
			return Transaction{}, err
		}
		return Transaction{}, fmt.Errorf("error fetching account %w", err)
	}
	tanID, err := NewRandTransactionID()
	if err != nil {
		return Transaction{}, fmt.Errorf("error generating transactionID %w", err)
	}
	tan, err := NewTransaction(tanID, req.AccountNumber, req.UserID, req.Amount, req.Currency, req.Type, req.Reference)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid transaction details %w", err)
	}

	newAcct := acct
	switch tan.Type {
	case Deposit:
		newAcct, err = newAcct.Deposit(tan.Amount)
	case Withdrawal:
		newAcct, err = newAcct.Withdraw(tan.Amount)
	}
	if err != nil {
		return Transaction{}, fmt.Errorf("error processing transaction %w", err)
	}

	err = svc.transactionStore.Put(tan)
	if err != nil {
		return Transaction{}, fmt.Errorf("error processing transaction %w", err)
	}
	err = svc.acctStore.Put(newAcct)
	if err != nil {
		return Transaction{}, fmt.Errorf("error processing transaction %w", err)
	}
	return tan, nil
}

func (svc *TransactionService) ListTransactions(acctNum accounts.AccountNumber) ([]Transaction, error) {
	tans, err := svc.transactionStore.GetByAccountNumber(acctNum)
	if err != nil {
		if errors.Is(err, ErrTransactionNotFound) {
			return []Transaction{}, nil
		}
		return nil, fmt.Errorf("error listing transactions %w", err)
	}
	return tans, nil
}
