package users_test

import (
	"eaglebank/internal/users"
	"eaglebank/internal/users/adapters"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserService(t *testing.T) {
	store := adapters.NewInMemoryUserStore()
	svc := users.NewUserService(store)
	t.Run("create user", func(t *testing.T) {
		t.Run("should successfully create user", func(t *testing.T) {
			usr, err := svc.CreateUser(users.CreateUserRequest{
				Name: "name",
				Address: users.Address{
					Line1:    "line1",
					Town:     "town",
					County:   "county",
					Postcode: "postcode",
				},
				PhoneNumber: "+440000000000",
				Email:       "foo@bar.com",
			})
			require.NoError(t, err)

			retUsr, err := store.Get(usr.ID)
			require.NoError(t, err)
			assert.Equal(t, retUsr, usr)
		})
	})
	t.Run("should fail for invalid user", func(t *testing.T) {
		_, err := svc.CreateUser(users.CreateUserRequest{
			Name: "name",
			Address: users.Address{
				Line1:    "line1",
				Town:     "town",
				County:   "county",
				Postcode: "postcode",
			},
			Email: "foo@bar.com",
		})
		assert.Error(t, err)
	})
	t.Run("should fail if put fails", func(t *testing.T) {
		usrStore := newFailingUserStore(t)
		failSvc := users.NewUserService(usrStore)
		_, err := failSvc.CreateUser(users.CreateUserRequest{
			Name: "name",
			Address: users.Address{
				Line1:    "line1",
				Town:     "town",
				County:   "county",
				Postcode: "postcode",
			},
			PhoneNumber: "+440000000000",
			Email:       "foo@bar.com",
		})
		assert.Error(t, err)
	})
}

type failingUserStore struct{}

func newFailingUserStore(t *testing.T) failingUserStore {
	t.Helper()
	return failingUserStore{}
}

func (f failingUserStore) Get(id users.UserID) (users.User, error) {
	//TODO implement me
	panic("implement me")
}

func (f failingUserStore) Put(user users.User) error {
	return errors.New("error")
}

func (f failingUserStore) Delete(id users.UserID) error {
	//TODO implement me
	panic("implement me")
}
