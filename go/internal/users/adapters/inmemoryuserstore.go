package adapters

import (
	"eaglebank/internal/users"
	"sync"
)

type InMemoryUserStore struct {
	mu    sync.RWMutex
	store map[users.UserID]users.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{store: map[users.UserID]users.User{}}
}

func (s *InMemoryUserStore) Get(id users.UserID) (users.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.store[id]
	if !ok {
		return users.User{}, users.ErrUserNotFound
	}
	return user, nil
}

func (s *InMemoryUserStore) Put(user users.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[user.ID] = user
	return nil
}

func (s *InMemoryUserStore) Delete(id users.UserID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, id)
	return nil
}
