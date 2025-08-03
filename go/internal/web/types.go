package web

import "time"

type Address struct {
	Line1    string  `json:"line1" validate:"required"`
	Line2    *string `json:"line2,omitempty"`
	Line3    *string `json:"line3,omitempty"`
	Town     string  `json:"town" validate:"required"`
	County   string  `json:"county" validate:"required"`
	Postcode string  `json:"postcode" validate:"required"`
}

type CreateBankAccountRequest struct {
	Name        string `json:"name" validate:"required"`
	AccountType string `json:"accountType" validate:"required,oneof=personal"`
}

type UpdateBankAccountRequest struct {
	Name        *string `json:"name,omitempty"`
	AccountType *string `json:"accountType,omitempty" validate:"omitempty,oneof=personal"`
}

type BankAccountResponse struct {
	AccountNumber    string    `json:"accountNumber" validate:"required,regexp=^01\\d{6}$"`
	SortCode         string    `json:"sortCode" validate:"required,eq=10-10-10"`
	Name             string    `json:"name" validate:"required"`
	AccountType      string    `json:"accountType" validate:"required,oneof=personal"`
	Balance          float64   `json:"balance" validate:"required,min=0,max=10000"`
	Currency         string    `json:"currency" validate:"required,oneof=GBP"`
	CreatedTimestamp time.Time `json:"createdTimestamp" validate:"required"`
	UpdatedTimestamp time.Time `json:"updatedTimestamp" validate:"required"`
}

type ListBankAccountsResponse struct {
	Accounts []BankAccountResponse `json:"accounts" validate:"required"`
}

type CreateTransactionRequest struct {
	Amount    float64 `json:"amount" validate:"required,min=0,max=10000"`
	Currency  string  `json:"currency" validate:"required,oneof=GBP"`
	Type      string  `json:"type" validate:"required,oneof=deposit withdrawal"`
	Reference *string `json:"reference,omitempty"`
}

type TransactionResponse struct {
	ID               string    `json:"id" validate:"required,regexp=^tan-[A-Za-z0-9]$"`
	Amount           float64   `json:"amount" validate:"required,min=0,max=10000"`
	Currency         string    `json:"currency" validate:"required,oneof=GBP"`
	Type             string    `json:"type" validate:"required,oneof=deposit withdrawal"`
	Reference        *string   `json:"reference,omitempty"`
	UserID           *string   `json:"userId,omitempty" validate:"omitempty,regexp=^usr-[A-Za-z0-9]+$"`
	CreatedTimestamp time.Time `json:"createdTimestamp" validate:"required"`
}

type ListTransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions" validate:"required"`
}

type CreateUserRequest struct {
	Name        string  `json:"name" validate:"required"`
	Address     Address `json:"address" validate:"required"`
	PhoneNumber string  `json:"phoneNumber" validate:"required,e164"`
	Email       string  `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name        *string  `json:"name,omitempty"`
	Address     *Address `json:"address,omitempty"`
	PhoneNumber *string  `json:"phoneNumber,omitempty" validate:"omitempty,e164"`
	Email       *string  `json:"email,omitempty" validate:"omitempty,email"`
}

type UserResponse struct {
	ID               string    `json:"id" validate:"required,regexp=^usr-[A-Za-z0-9]+$"`
	Name             string    `json:"name" validate:"required"`
	Address          Address   `json:"address" validate:"required"`
	PhoneNumber      string    `json:"phoneNumber" validate:"required,e164"`
	Email            string    `json:"email" validate:"required,email"`
	CreatedTimestamp time.Time `json:"createdTimestamp" validate:"required"`
	UpdatedTimestamp time.Time `json:"updatedTimestamp" validate:"required"`
}

type ErrorResponse struct {
	Message string `json:"message" validate:"required"`
}

type ValidationDetail struct {
	Field   string `json:"field" validate:"required"`
	Message string `json:"message" validate:"required"`
	Type    string `json:"type" validate:"required"`
}

type BadRequestErrorResponse struct {
	Message string             `json:"message" validate:"required"`
	Details []ValidationDetail `json:"details" validate:"required"`
}
