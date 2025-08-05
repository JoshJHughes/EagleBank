package web

import (
	"eaglebank/internal/validation"
	"encoding/json"
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
