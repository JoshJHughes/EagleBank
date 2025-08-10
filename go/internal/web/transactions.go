package web

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/transactions"
	"eaglebank/internal/validation"
	"encoding/json"
	"errors"
	"net/http"
)

func handleCreateTransaction(tanSvc TransactionService, acctSvc AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err := validation.Get().Struct(req)
		if err != nil {
			writeBadRequestErrorResponse(w, err)
			return
		}

		acct, err := checkTransactionAuth(w, r, acctSvc)
		if err != nil {
			return
		}

		ref := ""
		if req.Reference != nil {
			ref = *req.Reference
		}
		domReq, err := transactions.NewCreateTransactionRequest(acct.AccountNumber, acct.UserID, req.Amount, accounts.Currency(req.Currency), transactions.TransactionType(req.Type), ref)
		if err != nil {
			writeBadRequestErrorResponse(w, err)
			return
		}

		tan, err := tanSvc.CreateTransaction(domReq)
		if err != nil {
			if errors.Is(err, accounts.ErrInsufficientFunds) {
				writeErrorResponse(w, http.StatusUnprocessableEntity, err)
				return
			}
			writeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		resp := newTransactionResponseFromDomain(tan)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

func handleListTransactions(svc TransactionService, acctSvc AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acct, err := checkTransactionAuth(w, r, acctSvc)
		if err != nil {
			return
		}

		tans, err := svc.ListTransactions(acct.AccountNumber)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, err)
		}

		tanResps := make([]TransactionResponse, 0, len(tans))
		for _, tan := range tans {
			tanResps = append(tanResps, newTransactionResponseFromDomain(tan))
		}

		resp := ListTransactionsResponse{Transactions: tanResps}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func checkTransactionAuth(w http.ResponseWriter, r *http.Request, acctSvc AccountService) (accounts.BankAccount, error) {
	acctNum, err := accounts.NewAccountNumber(r.PathValue("accountId"))
	if err != nil {
		writeBadRequestErrorResponse(w, err)
		return accounts.BankAccount{}, err
	}

	acct, err := acctSvc.FetchAccount(acctNum)
	if err != nil {
		if errors.Is(err, accounts.ErrAccountNotFound) {
			writeErrorResponse(w, http.StatusNotFound, err)
			return accounts.BankAccount{}, err
		}
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return accounts.BankAccount{}, err
	}
	userID := GetAuthenticatedUserID(r.Context())
	if acct.UserID.String() != userID {
		err = errors.New("forbidden")
		writeErrorResponse(w, http.StatusForbidden, err)
		return accounts.BankAccount{}, err
	}
	return acct, nil
}
