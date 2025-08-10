package transactions

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/users"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const Deposit TransactionType = "deposit"
const Withdrawal TransactionType = "withdrawal"

func (t TransactionType) String() string { return string(t) }

func (t TransactionType) IsValid() bool {
	switch t {
	case Deposit, Withdrawal:
		return true
	default:
		return false
	}
}

type TransactionID string

var transactionIDRegex = regexp.MustCompile(`^tan-[A-Za-z0-9]+$`)

func (id TransactionID) String() string { return string(id) }

func (id TransactionID) IsValid() bool { return transactionIDRegex.MatchString(id.String()) }

func NewTransactionID(s string) (TransactionID, error) {
	id := TransactionID(s)
	if !id.IsValid() {
		return "", fmt.Errorf("invalid transaction ID %q", s)
	}
	return id, nil
}

func NewRandTransactionID() (TransactionID, error) {
	id := uuid.New()
	clean := strings.ReplaceAll(id.String(), "-", "")
	return NewTransactionID("tan-" + clean)
}

const TransactionMax float64 = 10000
const TransactionMin float64 = 0

type Transaction struct {
	ID               TransactionID
	AccountNumber    accounts.AccountNumber
	UserID           users.UserID
	Amount           float64
	Currency         accounts.Currency
	Type             TransactionType
	Reference        string
	CreatedTimestamp time.Time
}

func (t Transaction) IsValid() bool {
	if t.Amount < TransactionMin || t.Amount > TransactionMax {
		return false
	}
	if !t.ID.IsValid() {
		return false
	}
	if !t.AccountNumber.IsValid() {
		return false
	}
	if !t.UserID.IsValid() {
		return false
	}
	if !t.Currency.IsValid() {
		return false
	}
	if !t.Type.IsValid() {
		return false
	}
	return true
}

func NewTransaction(id TransactionID, acctNum accounts.AccountNumber, userID users.UserID, amt float64, curr accounts.Currency, tanType TransactionType, ref string) (Transaction, error) {
	now := time.Now()
	tan := Transaction{
		ID:               id,
		AccountNumber:    acctNum,
		UserID:           userID,
		Amount:           amt,
		Currency:         curr,
		Type:             tanType,
		Reference:        ref,
		CreatedTimestamp: now,
	}
	if !tan.IsValid() {
		return Transaction{}, fmt.Errorf("invalid transaction %+v", tan)
	}
	return tan, nil
}

type CreateTransactionRequest struct {
	AccountNumber accounts.AccountNumber
	UserID        users.UserID
	Amount        float64
	Currency      accounts.Currency
	Type          TransactionType
	Reference     string
}

func (r CreateTransactionRequest) IsValid() bool {
	if r.Amount < TransactionMin || r.Amount > TransactionMax {
		return false
	}
	if r.AccountNumber.IsValid() {
		return false
	}
	if r.UserID.IsValid() {
		return false
	}
	if !r.Currency.IsValid() {
		return false
	}
	if !r.Type.IsValid() {
		return false
	}
	return true
}

func NewCreateTransactionRequest(acctNum accounts.AccountNumber, userID users.UserID, amt float64, curr accounts.Currency, tanType TransactionType, ref string) (CreateTransactionRequest, error) {
	req := CreateTransactionRequest{
		AccountNumber: acctNum,
		UserID:        userID,
		Amount:        amt,
		Currency:      curr,
		Type:          tanType,
		Reference:     ref,
	}
	if !req.IsValid() {
		return CreateTransactionRequest{}, fmt.Errorf("invalid create transaction request %+v", req)
	}
	return req, nil
}
