package adapters

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/users"
	"sync"
)

type InMemoryAccountStore struct {
	mu            sync.RWMutex
	acctsByNumber map[accounts.AccountNumber]accounts.BankAccount
	acctsByUserID map[users.UserID][]accounts.BankAccount
}

func NewInMemoryAccountStore() *InMemoryAccountStore {
	return &InMemoryAccountStore{
		acctsByNumber: make(map[accounts.AccountNumber]accounts.BankAccount),
		acctsByUserID: make(map[users.UserID][]accounts.BankAccount),
	}
}

func (s *InMemoryAccountStore) GetByAcctNum(acctNum accounts.AccountNumber) (accounts.BankAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	acct, ok := s.acctsByNumber[acctNum]
	if !ok {
		return accounts.BankAccount{}, accounts.ErrAccountNotFound
	}
	return acct, nil
}

func (s *InMemoryAccountStore) GetByUserID(userID users.UserID) ([]accounts.BankAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	accts, ok := s.acctsByUserID[userID]
	if !ok {
		return nil, accounts.ErrAccountNotFound
	}
	result := make([]accounts.BankAccount, len(accts))
	copy(result, accts)
	return result, nil
}

func (s *InMemoryAccountStore) Put(acct accounts.BankAccount) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.acctsByNumber[acct.AccountNumber] = acct
	for i, a := range s.acctsByUserID[acct.UserID] {
		if a.AccountNumber == acct.AccountNumber {
			s.acctsByUserID[acct.UserID][i] = acct
			return nil
		}
	}
	s.acctsByUserID[acct.UserID] = append(s.acctsByUserID[acct.UserID], acct)
	return nil
}

func (s *InMemoryAccountStore) Delete(acctNum accounts.AccountNumber) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	acctToDel, ok := s.acctsByNumber[acctNum]
	if !ok {
		return accounts.ErrAccountNotFound
	}
	delete(s.acctsByNumber, acctNum)

	accts, ok := s.acctsByUserID[acctToDel.UserID]
	if !ok {
		return accounts.ErrAccountNotFound
	}
	for i, acct := range accts {
		if acct.AccountNumber == acctToDel.AccountNumber {
			copy(accts[i:], accts[i+1:])
			accts[len(accts)-1] = accounts.BankAccount{}
			accts = accts[:len(accts)-1]
			s.acctsByUserID[acctToDel.UserID] = accts
			break
		}
	}

	return nil
}
