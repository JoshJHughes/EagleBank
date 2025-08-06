package accounts

import (
	"eaglebank/internal/users"
	"errors"
	"fmt"
)

type AccountStore interface {
	GetByAcctNum(acctNum AccountNumber) (BankAccount, error)
	GetByUserID(userID users.UserID) ([]BankAccount, error)
	Put(acct BankAccount) error
	Delete(acctNum AccountNumber) error
}

type AccountService struct {
	accountStore AccountStore
}

func NewAccountService(acctStore AccountStore) *AccountService {
	return &AccountService{accountStore: acctStore}
}

func (svc *AccountService) CreateAccount(req CreateAccountRequest) (*BankAccount, error) {
	if !req.IsValid() {
		return nil, fmt.Errorf("invalid create account request %+v", req)
	}
	acctNum, err := NewRandAccountNumber()
	if err != nil {
		return nil, fmt.Errorf("error generating account number %w", err)
	}
	acct, err := NewBankAccount(
		req.UserID,
		acctNum,
		"10-10-10",
		req.Name,
		req.AccountType,
		GBP,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid bank account details")
	}
	err = svc.accountStore.Put(acct)
	if err != nil {
		return nil, fmt.Errorf("error creating bank account %w", err)
	}
	return &acct, nil
}

func (svc *AccountService) ListAccounts(id users.UserID) ([]BankAccount, error) {
	accts, err := svc.accountStore.GetByUserID(id)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			return []BankAccount{}, nil
		}
		return nil, fmt.Errorf("error listing bank accounts %w", err)
	}
	return accts, nil
}
