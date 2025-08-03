package users

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestUserID(t *testing.T) {
	t.Run("NewUserID", func(t *testing.T) {
		t.Run("should create valid user ID", func(t *testing.T) {
			validID := "usr-abc123"
			id, err := NewUserID(validID)
			require.NoError(t, err)
			assert.Equal(t, UserID(validID), id)
			assert.Equal(t, validID, id.String())
		})

		t.Run("should error on empty string", func(t *testing.T) {
			_, err := NewUserID("")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "user ID cannot be empty")
		})

		t.Run("should error on invalid format", func(t *testing.T) {
			invalidIDs := []string{
				"invalid",
				"usr-",
				"user-abc123",
				"usr-abc@123",
				"usr-abc 123",
			}
			for _, invalidID := range invalidIDs {
				t.Run(invalidID, func(t *testing.T) {
					_, err := NewUserID(invalidID)
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "invalid user ID format")
				})
			}
		})
	})

	t.Run("MustNewUserID", func(t *testing.T) {
		t.Run("should create valid user ID", func(t *testing.T) {
			validID := "usr-abc123"
			id := MustNewUserID(validID)
			assert.Equal(t, UserID(validID), id)
		})

		t.Run("should panic on invalid input", func(t *testing.T) {
			assert.Panics(t, func() {
				MustNewUserID("")
			})
		})
	})

	t.Run("NewRandUserID", func(t *testing.T) {
		t.Run("should create random user ID", func(t *testing.T) {
			id, err := NewRandUserID()
			require.NoError(t, err)
			assert.True(t, strings.HasPrefix(id.String(), "usr-"))
			assert.True(t, len(id.String()) > 4)
		})

		t.Run("should create different IDs on multiple calls", func(t *testing.T) {
			id1, err := NewRandUserID()
			require.NoError(t, err)
			id2, err := NewRandUserID()
			require.NoError(t, err)
			assert.NotEqual(t, id1, id2)
		})
	})

	t.Run("MustNewRandUserID", func(t *testing.T) {
		t.Run("should create random user ID", func(t *testing.T) {
			id := MustNewRandUserID()
			assert.True(t, strings.HasPrefix(id.String(), "usr-"))
			assert.True(t, len(id.String()) > 4)
		})
	})
}

func TestEmail(t *testing.T) {
	t.Run("NewEmail", func(t *testing.T) {
		t.Run("should create valid email", func(t *testing.T) {
			validEmails := []string{
				"test@example.com",
				"user.name@domain.co.uk",
				"user+tag@example.org",
			}
			for _, validEmail := range validEmails {
				t.Run(validEmail, func(t *testing.T) {
					email, err := NewEmail(validEmail)
					require.NoError(t, err)
					assert.Equal(t, Email(validEmail), email)
					assert.Equal(t, validEmail, email.String())
				})
			}
		})

		t.Run("should error on empty string", func(t *testing.T) {
			_, err := NewEmail("")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "email cannot be empty")
		})

		t.Run("should error on invalid format", func(t *testing.T) {
			invalidEmails := []string{
				"invalid",
				"@example.com",
				"test@",
				"test.example.com",
				"test @example.com",
			}
			for _, invalidEmail := range invalidEmails {
				t.Run(invalidEmail, func(t *testing.T) {
					_, err := NewEmail(invalidEmail)
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "invalid email format")
				})
			}
		})
	})

	t.Run("MustNewEmail", func(t *testing.T) {
		t.Run("should create valid email", func(t *testing.T) {
			validEmail := "test@example.com"
			email := MustNewEmail(validEmail)
			assert.Equal(t, Email(validEmail), email)
		})

		t.Run("should panic on invalid input", func(t *testing.T) {
			assert.Panics(t, func() {
				MustNewEmail("")
			})
		})
	})
}

func TestPhoneNumber(t *testing.T) {
	t.Run("NewPhoneNumber", func(t *testing.T) {
		t.Run("should create valid phone number", func(t *testing.T) {
			validNumbers := []string{
				"+1234567890123",
				"+44123456789",
				"+331234567890",
			}
			for _, validNumber := range validNumbers {
				t.Run(validNumber, func(t *testing.T) {
					phone, err := NewPhoneNumber(validNumber)
					require.NoError(t, err)
					assert.Equal(t, PhoneNumber(validNumber), phone)
					assert.Equal(t, validNumber, phone.String())
				})
			}
		})

		t.Run("should error on empty string", func(t *testing.T) {
			_, err := NewPhoneNumber("")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "phone number cannot be empty")
		})

		t.Run("should error on invalid format", func(t *testing.T) {
			invalidNumbers := []string{
				"1234567890",          // missing +
				"+0123456789",         // starts with 0
				"+123456789012345678", // too long
				"+1",                  // too short
				"+12a3456789",         // contains letter
				"+12 3456789",         // contains space
				"+12-3456789",         // contains dash
			}
			for _, invalidNumber := range invalidNumbers {
				t.Run(invalidNumber, func(t *testing.T) {
					_, err := NewPhoneNumber(invalidNumber)
					assert.Error(t, err)
					assert.Contains(t, err.Error(), "invalid phone number format")
				})
			}
		})
	})

	t.Run("MustNewPhoneNumber", func(t *testing.T) {
		t.Run("should create valid phone number", func(t *testing.T) {
			validNumber := "+1234567890"
			phone := MustNewPhoneNumber(validNumber)
			assert.Equal(t, PhoneNumber(validNumber), phone)
		})

		t.Run("should panic on invalid input", func(t *testing.T) {
			assert.Panics(t, func() {
				MustNewPhoneNumber("")
			})
		})
	})
}

func TestAddress(t *testing.T) {
	t.Run("NewAddress", func(t *testing.T) {
		t.Run("should create valid address with required fields", func(t *testing.T) {
			line1 := "123 Main St"
			town := "Springfield"
			county := "Simpson County"
			postcode := "12345"

			addr, err := NewAddress(line1, town, county, postcode)
			require.NoError(t, err)
			assert.Equal(t, line1, addr.Line1)
			assert.Equal(t, town, addr.Town)
			assert.Equal(t, county, addr.County)
			assert.Equal(t, postcode, addr.Postcode)
			assert.Empty(t, addr.Line2)
			assert.Empty(t, addr.Line3)
		})

		t.Run("should create valid address with optional fields", func(t *testing.T) {
			line1 := "123 Main St"
			line2 := "Apt 4B"
			line3 := "Building C"
			town := "Springfield"
			county := "Simpson County"
			postcode := "12345"

			addr, err := NewAddress(line1, town, county, postcode, WithLine2(line2), WithLine3(line3))
			require.NoError(t, err)
			assert.Equal(t, line1, addr.Line1)
			assert.Equal(t, line2, addr.Line2)
			assert.Equal(t, line3, addr.Line3)
			assert.Equal(t, town, addr.Town)
			assert.Equal(t, county, addr.County)
			assert.Equal(t, postcode, addr.Postcode)
		})

		t.Run("should error on missing required fields", func(t *testing.T) {
			validLine1 := "123 Main St"
			validTown := "Springfield"
			validCounty := "Simpson County"
			validPostcode := "12345"

			testCases := []struct {
				name     string
				line1    string
				town     string
				county   string
				postcode string
				errorMsg string
			}{
				{"empty line1", "", validTown, validCounty, validPostcode, "address line 1 is required"},
				{"empty town", validLine1, "", validCounty, validPostcode, "town is required"},
				{"empty county", validLine1, validTown, "", validPostcode, "county is required"},
				{"empty postcode", validLine1, validTown, validCounty, "", "postcode is required"},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					_, err := NewAddress(tc.line1, tc.town, tc.county, tc.postcode)
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tc.errorMsg)
				})
			}
		})
	})

	t.Run("MustNewAddress", func(t *testing.T) {
		t.Run("should create valid address", func(t *testing.T) {
			line1 := "123 Main St"
			town := "Springfield"
			county := "Simpson County"
			postcode := "12345"

			addr := MustNewAddress(line1, town, county, postcode)
			assert.Equal(t, line1, addr.Line1)
			assert.Equal(t, town, addr.Town)
			assert.Equal(t, county, addr.County)
			assert.Equal(t, postcode, addr.Postcode)
		})

		t.Run("should panic on invalid input", func(t *testing.T) {
			assert.Panics(t, func() {
				MustNewAddress("", "town", "county", "postcode")
			})
		})
	})

	t.Run("AddressOptions", func(t *testing.T) {
		t.Run("WithLine2", func(t *testing.T) {
			line2 := "Apt 4B"
			addr := Address{}
			opt := WithLine2(line2)
			opt(&addr)
			assert.Equal(t, line2, addr.Line2)
		})

		t.Run("WithLine3", func(t *testing.T) {
			line3 := "Building C"
			addr := Address{}
			opt := WithLine3(line3)
			opt(&addr)
			assert.Equal(t, line3, addr.Line3)
		})
	})
}

func TestUser(t *testing.T) {
	t.Run("NewUser", func(t *testing.T) {
		t.Run("should create valid user", func(t *testing.T) {
			id := MustNewUserID("usr-test123")
			name := "John Doe"
			addr := MustNewAddress("123 Main St", "Springfield", "Simpson County", "12345")
			phone := MustNewPhoneNumber("+1234567890")
			email := MustNewEmail("john@example.com")

			beforeTime := time.Now()
			user, err := NewUser(id, name, addr, phone, email)
			afterTime := time.Now()

			require.NoError(t, err)
			assert.Equal(t, id, user.ID)
			assert.Equal(t, name, user.Name)
			assert.Equal(t, addr, user.Address)
			assert.Equal(t, phone, user.PhoneNumber)
			assert.Equal(t, email, user.Email)
			assert.True(t, user.Created.After(beforeTime) || user.Created.Equal(beforeTime))
			assert.True(t, user.Created.Before(afterTime) || user.Created.Equal(afterTime))
			assert.True(t, user.Updated.After(beforeTime) || user.Updated.Equal(beforeTime))
			assert.True(t, user.Updated.Before(afterTime) || user.Updated.Equal(afterTime))
		})

		t.Run("should error on empty name", func(t *testing.T) {
			id := MustNewUserID("usr-test123")
			addr := MustNewAddress("123 Main St", "Springfield", "Simpson County", "12345")
			phone := MustNewPhoneNumber("+1234567890")
			email := MustNewEmail("john@example.com")

			_, err := NewUser(id, "", addr, phone, email)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "name is required")
		})
	})

	t.Run("MustNewUser", func(t *testing.T) {
		t.Run("should create valid user", func(t *testing.T) {
			id := MustNewUserID("usr-test123")
			name := "John Doe"
			addr := MustNewAddress("123 Main St", "Springfield", "Simpson County", "12345")
			phone := MustNewPhoneNumber("+1234567890")
			email := MustNewEmail("john@example.com")

			user := MustNewUser(id, name, addr, phone, email)
			assert.Equal(t, id, user.ID)
			assert.Equal(t, name, user.Name)
			assert.Equal(t, addr, user.Address)
			assert.Equal(t, phone, user.PhoneNumber)
			assert.Equal(t, email, user.Email)
		})

		t.Run("should panic on invalid input", func(t *testing.T) {
			id := MustNewUserID("usr-test123")
			addr := MustNewAddress("123 Main St", "Springfield", "Simpson County", "12345")
			phone := MustNewPhoneNumber("+1234567890")
			email := MustNewEmail("john@example.com")

			assert.Panics(t, func() {
				MustNewUser(id, "", addr, phone, email)
			})
		})
	})
}
