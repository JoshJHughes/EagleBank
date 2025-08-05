package web

import (
	"eaglebank/internal/users"
	"eaglebank/internal/validation"
	"encoding/json"
	"errors"
	"net/http"
)

func handleCreateUser(usrSvc UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err := validation.Get().Struct(req)
		if err != nil {
			writeBadRequestErrorResponse(w, err)
			return
		}

		usrReq, err := req.toDomain()
		if err != nil {
			writeBadRequestErrorResponse(w, err)
			return
		}
		usr, err := usrSvc.CreateUser(usrReq)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		resp := newUserResponseFromDomain(usr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

func handleGetUser(usrSvc UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := users.NewUserID(r.PathValue("userId"))
		if err != nil {
			writeBadRequestErrorResponse(w, err)
			return
		}

		authenticatedUserID := GetAuthenticatedUserID(r.Context())
		if authenticatedUserID != userID.String() {
			writeErrorResponse(w, http.StatusForbidden, errors.New("forbidden"))
		}

		usr, err := usrSvc.GetUser(userID)
		if err != nil {
			if errors.Is(err, users.ErrUserNotFound) {
				writeErrorResponse(w, http.StatusNotFound, err)
				return
			}
			writeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		resp := newUserResponseFromDomain(usr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
