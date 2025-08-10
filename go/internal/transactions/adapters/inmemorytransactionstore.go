package adapters

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/transactions"
	"fmt"
	"sync"
)

type InMemoryTransactionStore struct {
	mu            sync.RWMutex
	tansByTanID   map[transactions.TransactionID]transactions.Transaction
	tansByAcctNum map[accounts.AccountNumber][]transactions.Transaction
}

func NewInMemoryTransactionStore() *InMemoryTransactionStore {
	return &InMemoryTransactionStore{
		tansByTanID:   make(map[transactions.TransactionID]transactions.Transaction),
		tansByAcctNum: make(map[accounts.AccountNumber][]transactions.Transaction),
	}
}

func (s *InMemoryTransactionStore) GetByTransactionID(tanID transactions.TransactionID) (transactions.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tan, ok := s.tansByTanID[tanID]
	if !ok {
		return transactions.Transaction{}, transactions.ErrTransactionNotFound
	}
	return tan, nil
}

func (s *InMemoryTransactionStore) GetByAccountNumber(acctNum accounts.AccountNumber) ([]transactions.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tans, ok := s.tansByAcctNum[acctNum]
	if !ok {
		return nil, transactions.ErrTransactionNotFound
	}
	result := make([]transactions.Transaction, len(tans))
	copy(result, tans)
	return result, nil
}

func (s *InMemoryTransactionStore) Put(tan transactions.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.tansByTanID[tan.ID]
	if exists {
		return fmt.Errorf("cannot modify transaction")
	}
	for _, t := range s.tansByAcctNum[tan.AccountNumber] {
		if t.ID == tan.ID {
			return fmt.Errorf("cannot modify transaction")
		}
	}
	s.tansByTanID[tan.ID] = tan
	s.tansByAcctNum[tan.AccountNumber] = append(s.tansByAcctNum[tan.AccountNumber], tan)
	return nil
}
