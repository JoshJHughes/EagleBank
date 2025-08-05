package adapters

import (
	"eaglebank/internal/users"
)

type InMemoryUserStore struct {
	store map[users.UserID]users.User
}

func NewInMemoryUserStore() InMemoryUserStore {
	return InMemoryUserStore{store: map[users.UserID]users.User{}}
}

func (s InMemoryUserStore) Get(id users.UserID) (users.User, error) {
	user, ok := s.store[id]
	if !ok {
		return users.User{}, users.ErrUserNotFound
	}
	return user, nil
}

func (s InMemoryUserStore) Put(user users.User) error {
	s.store[user.ID] = user
	return nil
}

func (s InMemoryUserStore) Delete(id users.UserID) error {
	delete(s.store, id)
	return nil
}
