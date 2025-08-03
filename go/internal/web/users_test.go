package web

import (
	"bytes"
	"eaglebank/internal/users"
	"eaglebank/internal/users/adapters"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestUsers(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	validate := validator.New(validator.WithRequiredStructEnabled())

	usrStore := adapters.NewInMemoryUserStore()
	usrSvc := users.NewUserService(usrStore)

	srv := NewServer(logger, validate, usrSvc)
	t.Run("POST to /v1/users", func(t *testing.T) {
		t.Run("with all required data should create user", func(t *testing.T) {
			reqObj := CreateUserRequest{
				Name: "Mr Foo",
				Address: Address{
					Line1:    "line 1",
					Town:     "town",
					County:   "county",
					Postcode: "postcode",
				},
				PhoneNumber: "+440000000000",
				Email:       "foo@bar.com",
			}
			by, err := json.Marshal(reqObj)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(by))
			srv.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusCreated, rr.Code)

			var usrResp UserResponse
			err = json.NewDecoder(rr.Body).Decode(&usrResp)
			require.NoError(t, err)

			assert.Equal(t, reqObj.Name, usrResp.Name)
			assert.Equal(t, reqObj.Address, usrResp.Address)
			assert.Equal(t, reqObj.PhoneNumber, usrResp.PhoneNumber)
			assert.Equal(t, reqObj.Email, usrResp.Email)

			storeUsr, err := usrStore.Get(users.MustNewUserID(usrResp.ID))
			require.NoError(t, err)
			storeUsrResp := newUserResponseFromDomain(storeUsr)
			assertUserResponseEqual(t, usrResp, storeUsrResp)
		})
		t.Run("without all required data should return bad request", func(t *testing.T) {
			reqObj := CreateUserRequest{
				Name:        "Mr Foo",
				PhoneNumber: "+440000000000",
				Email:       "foo@bar.com",
			}
			by, err := json.Marshal(reqObj)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(by))
			srv.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusBadRequest, rr.Code)

			var errResp string
			err = json.NewDecoder(rr.Body).Decode(&errResp)
			require.NoError(t, err)
			assert.NotEmpty(t, errResp)
		})
		t.Run("unexpected error should return internal server error", func(t *testing.T) {
			errUsrSvc := NewErroringUserService(t)
			errSrv := NewServer(logger, validate, errUsrSvc)
			reqObj := CreateUserRequest{
				Name: "Mr Foo",
				Address: Address{
					Line1:    "line 1",
					Town:     "town",
					County:   "county",
					Postcode: "postcode",
				},
				PhoneNumber: "+440000000000",
				Email:       "foo@bar.com",
			}
			by, err := json.Marshal(reqObj)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(by))
			errSrv.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusInternalServerError, rr.Code)

			var errResp ErrorResponse
			err = json.NewDecoder(rr.Body).Decode(&errResp)
			require.NoError(t, err)
			assert.NotEmpty(t, errResp.Message)

		})
	})

}

func assertUserResponseEqual(t *testing.T, expected UserResponse, actual UserResponse) {
	t.Helper()

	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Address, actual.Address)
	assert.Equal(t, expected.PhoneNumber, actual.PhoneNumber)
	assert.Equal(t, expected.Email, actual.Email)
	assert.WithinDuration(t, expected.CreatedTimestamp, actual.CreatedTimestamp, time.Millisecond*100)
	assert.WithinDuration(t, expected.UpdatedTimestamp, actual.UpdatedTimestamp, time.Millisecond*100)
}

type ErroringUserService struct{}

func NewErroringUserService(t *testing.T) ErroringUserService {
	t.Helper()
	return ErroringUserService{}
}

func (e ErroringUserService) CreateUser(_ users.CreateUserRequest) (users.User, error) {
	return users.User{}, errors.New("some error")
}
