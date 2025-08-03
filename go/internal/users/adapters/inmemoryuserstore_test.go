package adapters

import (
	"eaglebank/internal/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewInMemoryUserStore(t *testing.T) {
	store := NewInMemoryUserStore()

	t.Run("should error getting user which does not exist", func(t *testing.T) {
		missingID := users.MustNewUserID("usr-missing")
		_, err := store.Get(missingID)
		assert.Error(t, err)
	})
	t.Run("should not error deleting user which does not exist", func(t *testing.T) {
		missingID := users.MustNewUserID("usr-missing")
		err := store.Delete(missingID)
		assert.NoError(t, err)
	})
	t.Run("should perform put-get-update-delete cycle without errors", func(t *testing.T) {
		usr := newTestUser(t)
		t.Run("should create user that does not exist in store", func(t *testing.T) {
			err := store.Put(usr)
			require.NoError(t, err)
		})
		t.Run("should get an existing user", func(t *testing.T) {
			gotUsr, err := store.Get(usr.ID)
			require.NoError(t, err)
			require.Equal(t, usr, gotUsr)
		})
		t.Run("should update existing user", func(t *testing.T) {
			updatedUsr := usr
			updatedUsr.Name = "new name"
			require.NotEqual(t, usr.Name, updatedUsr.Name)

			err := store.Put(updatedUsr)
			require.NoError(t, err)

			gotUsr, err := store.Get(usr.ID)
			require.NoError(t, err)
			require.Equal(t, updatedUsr, gotUsr)
		})
		t.Run("should delete existing user", func(t *testing.T) {
			err := store.Delete(usr.ID)
			require.NoError(t, err)

			require.Empty(t, store.store)
		})
	})

}

func newTestUser(t *testing.T) users.User {
	t.Helper()

	usrID := users.MustNewRandUserID()
	name := "Mr Foo"
	addr := users.MustNewAddress("address line1", "town", "county", "postcode")
	phone := users.MustNewPhoneNumber("+440000000000")
	email := users.MustNewEmail("foo@bar.com")
	return users.MustNewUser(usrID, name, addr, phone, email)
}
