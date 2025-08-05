package adapters

import (
	"eaglebank/internal/accounts"
)

type InMemoryAccountStore struct {
	store map[accounts.AccountNumber]accounts.BankAccount
}

func NewInMemoryAccountStore() InMemoryAccountStore {
	return InMemoryAccountStore{store: map[accounts.AccountNumber]accounts.BankAccount{}}
}

func (s InMemoryAccountStore) Get(acctNum accounts.AccountNumber) (*accounts.BankAccount, error) {
	acct, ok := s.store[acctNum]
	if !ok {
		return &accounts.BankAccount{}, accounts.ErrAccountNotFound
	}
	return &acct, nil
}

func (s InMemoryAccountStore) Put(acct *accounts.BankAccount) error {
	s.store[acct.AccountNumber] = *acct
	return nil
}

func (s InMemoryAccountStore) Delete(acctNum accounts.AccountNumber) error {
	delete(s.store, acctNum)
	return nil
}
