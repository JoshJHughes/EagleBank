package web

import (
	"eaglebank/internal/users"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {
	t.Run("ErrorResponse", func(t *testing.T) {
		t.Run("should return correctly specified ErrorResponse", func(t *testing.T) {
			err := errors.New("some error")
			resp := newErrorResponse(err)
			assert.Equal(t, "some error", resp.Message)
		})
		t.Run("should return default for missing error message", func(t *testing.T) {
			err := errors.New("")
			resp := newErrorResponse(err)
			assert.Equal(t, "unspecified error", resp.Message)
		})
	})
	t.Run("BadRequestErrorResponse", func(t *testing.T) {
		t.Run("should return BadRequestErrorResponse for variable validation", func(t *testing.T) {
			// force validation error
			_, err := users.NewUserID("invalid id")
			resp := newBadRequestErrorResponse(err)
			assert.NotEmpty(t, resp.Message)
			assert.NotEmpty(t, resp.Details)
		})
		t.Run("should return BadRequestErrorResponse for struct validation", func(t *testing.T) {
			// force validation error
			_, err := users.NewUser("invalid id", "", users.Address{}, "not-a-phone", "not-an-email")
			resp := newBadRequestErrorResponse(err)
			assert.NotEmpty(t, resp.Message)
			assert.NotEmpty(t, resp.Details)
		})
	})
}
