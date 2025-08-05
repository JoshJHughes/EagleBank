package web

import (
	"bytes"
	"eaglebank/internal/users"
	"eaglebank/internal/users/adapters"
	"encoding/json"
	"errors"
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

	usrStore := adapters.NewInMemoryUserStore()
	usrSvc := users.NewUserService(usrStore)

	srv := NewServer(logger, usrSvc)
	t.Run("POST to /v1/users", func(t *testing.T) {
		t.Run("with all required data should create user", func(t *testing.T) {
			rr := httptest.NewRecorder()
			reqObj := validUserRequest
			req := createUserReq(t, reqObj)
			srv.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusCreated, rr.Code)

			var usrResp UserResponse
			err := json.NewDecoder(rr.Body).Decode(&usrResp)
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
			rr := httptest.NewRecorder()
			reqObj := validUserRequest
			reqObj.Email = "invalid-email"
			req := createUserReq(t, reqObj)
			srv.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusBadRequest, rr.Code)

			var errResp BadRequestErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&errResp)
			require.NoError(t, err)
			assert.NotEmpty(t, errResp)
		})
		t.Run("unexpected error should return internal server error", func(t *testing.T) {
			errUsrSvc := NewErroringUserService(t)
			errSrv := NewServer(logger, errUsrSvc)

			rr := httptest.NewRecorder()
			reqObj := validUserRequest
			req := createUserReq(t, reqObj)
			errSrv.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusInternalServerError, rr.Code)

			var errResp ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&errResp)
			require.NoError(t, err)
			assert.NotEmpty(t, errResp.Message)

		})
	})
	t.Run("GET /v1/users/{userId}", func(t *testing.T) {
		createRR := httptest.NewRecorder()
		validUserReq := validUserRequest
		req := createUserReq(t, validUserReq)
		srv.ServeHTTP(createRR, req)
		var user UserResponse
		err := json.NewDecoder(createRR.Body).Decode(&user)
		require.NoError(t, err)

		token := login(t, srv, user.ID)

		t.Run("200 on valid request", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req = getUserReq(t, user.ID, token)
			srv.ServeHTTP(rr, req)

			var resp UserResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, user, resp)
		})
		t.Run("400 on invalid request", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req = getUserReq(t, "not-a-userid", token)
			srv.ServeHTTP(rr, req)

			var resp BadRequestErrorResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
		})
		t.Run("401 on unauthorized", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req = getUserReq(t, user.ID)
			srv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusUnauthorized, rr.Code)
		})
		t.Run("403 on forbidden", func(t *testing.T) {
			rr := httptest.NewRecorder()
			forbiddenUserID := users.MustNewUserID("usr-forbidden")
			req = getUserReq(t, forbiddenUserID.String(), token)
			srv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusForbidden, rr.Code)
		})
		t.Run("404 on user not found", func(t *testing.T) {
			rr := httptest.NewRecorder()
			missingUserID := users.MustNewUserID("usr-missing")
			missingUserToken := login(t, srv, missingUserID.String())

			req = getUserReq(t, missingUserID.String(), missingUserToken)
			srv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusNotFound, rr.Code)
		})
		t.Run("500 on unexpected error", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req = getUserReq(t, user.ID, token)

			errUsrSvc := NewErroringUserService(t)
			errSrv := NewServer(logger, errUsrSvc)
			errSrv.ServeHTTP(rr, req)

			var resp ErrorResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			assert.Equal(t, http.StatusInternalServerError, rr.Code)
		})
	})
}

func createUserReq(t *testing.T, reqObj CreateUserRequest) *http.Request {
	t.Helper()
	by, err := json.Marshal(reqObj)
	require.NoError(t, err)
	return httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(by))
}

func getUserReq(t *testing.T, userID string, token ...string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+userID, nil)
	if len(token) != 0 {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}
	return req
}

func login(t *testing.T, srv http.Handler, userID string) string {
	t.Helper()
	loginBody := LoginRequest{
		UserID:       userID,
		PasswordHash: "123",
	}
	by, err := json.Marshal(loginBody)
	require.NoError(t, err)
	loginReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(by))
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	srv.ServeHTTP(loginRR, loginReq)

	var loginResp LoginResponse
	err = json.NewDecoder(loginRR.Body).Decode(&loginResp)
	require.NoError(t, err)

	return loginResp.Token
}

var validUserRequest = CreateUserRequest{
	Name: "name",
	Address: Address{
		Line1:    "line1",
		Town:     "town",
		County:   "county",
		Postcode: "postcode",
	},
	PhoneNumber: "+440000000000",
	Email:       "foo@bar.com",
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

func (e ErroringUserService) GetUser(userID users.UserID) (users.User, error) {
	return users.User{}, errors.New("some error")
}
