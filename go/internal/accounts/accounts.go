package accounts

import (
	"fmt"
	"time"
)

type AccountStore interface {
	Get(acctNum AccountNumber) (*BankAccount, error)
	Put(acct *BankAccount) error
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
	now := time.Now()
	acct := BankAccount{
		AccountNumber:    acctNum,
		SortCode:         "01-01-01",
		Name:             req.Name,
		AccountType:      req.AccountType,
		balance:          0,
		Currency:         GBP,
		CreatedTimestamp: now,
		UpdatedTimestamp: now,
	}
	err = svc.accountStore.Put(&acct)
	if err != nil {
		return nil, fmt.Errorf("error creating bank account %w", err)
	}
	return &acct, nil
}
