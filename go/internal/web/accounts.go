package web

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/users"
	"eaglebank/internal/validation"
	"encoding/json"
	"net/http"
)

func handleCreateAccount(svc AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateBankAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err := validation.Get().Struct(req)
		if err != nil {
			writeBadRequestErrorResponse(w, err)
			return
		}

		userID := GetAuthenticatedUserID(r.Context())
		domReq, err := accounts.NewCreateAccountRequest(users.UserID(userID), req.Name, accounts.AccountType(req.AccountType))
		if err != nil {
			writeBadRequestErrorResponse(w, err)
			return
		}

		acct, err := svc.CreateAccount(domReq)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, err)
		}

		resp := newBankAccountResponseFromDomain(*acct)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

func handleListAccounts(svc AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetAuthenticatedUserID(r.Context())
		accts, err := svc.ListAccounts(users.UserID(userID))
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, err)
		}

		acctResps := make([]BankAccountResponse, 0, len(accts))
		for _, acct := range accts {
			acctResps = append(acctResps, newBankAccountResponseFromDomain(acct))
		}

		resp := ListBankAccountsResponse{Accounts: acctResps}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
