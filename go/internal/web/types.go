package web

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/transactions"
	"eaglebank/internal/users"
	"eaglebank/internal/validation"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

type Address struct {
	Line1    string  `json:"line1" validate:"required"`
	Line2    *string `json:"line2,omitempty"`
	Line3    *string `json:"line3,omitempty"`
	Town     string  `json:"town" validate:"required"`
	County   string  `json:"county" validate:"required"`
	Postcode string  `json:"postcode" validate:"required"`
}

func (a Address) toDomain() (users.Address, error) {
	var opts []users.AddressOption
	if a.Line2 != nil {
		opts = append(opts, users.WithLine2(*a.Line2))
	}
	if a.Line3 != nil {
		opts = append(opts, users.WithLine3(*a.Line3))
	}
	return users.NewAddress(a.Line1, a.Town, a.County, a.Postcode, opts...)
}

func newAddressFromDomain(adr users.Address) Address {
	a := Address{
		Line1:    adr.Line1,
		Town:     adr.Town,
		County:   adr.County,
		Postcode: adr.Postcode,
	}
	if adr.Line2 != "" {
		line2 := adr.Line2
		a.Line2 = &line2
	}
	if adr.Line3 != "" {
		line3 := adr.Line3
		a.Line3 = &line3
	}
	return a
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
	AccountNumber    string    `json:"accountNumber" validate:"required,acctNum"`
	SortCode         string    `json:"sortCode" validate:"required,eq=10-10-10"`
	Name             string    `json:"name" validate:"required"`
	AccountType      string    `json:"accountType" validate:"required,oneof=personal"`
	Balance          float64   `json:"balance" validate:"required,min=0,max=10000"`
	Currency         string    `json:"currency" validate:"required,oneof=GBP"`
	CreatedTimestamp time.Time `json:"createdTimestamp" validate:"required"`
	UpdatedTimestamp time.Time `json:"updatedTimestamp" validate:"required"`
}

func newBankAccountResponseFromDomain(acct accounts.BankAccount) BankAccountResponse {
	return BankAccountResponse{
		AccountNumber:    acct.AccountNumber.String(),
		SortCode:         acct.SortCode.String(),
		Name:             acct.Name,
		AccountType:      acct.AccountType.String(),
		Balance:          acct.Balance(),
		Currency:         acct.Currency.String(),
		CreatedTimestamp: acct.CreatedTimestamp,
		UpdatedTimestamp: acct.UpdatedTimestamp,
	}
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
	ID               string    `json:"id" validate:"required,tanID"`
	Amount           float64   `json:"amount" validate:"required,min=0,max=10000"`
	Currency         string    `json:"currency" validate:"required,oneof=GBP"`
	Type             string    `json:"type" validate:"required,oneof=deposit withdrawal"`
	Reference        *string   `json:"reference,omitempty"`
	UserID           *string   `json:"userId,omitempty" validate:"omitempty,userID"`
	CreatedTimestamp time.Time `json:"createdTimestamp" validate:"required"`
}

func newTransactionResponseFromDomain(tan transactions.Transaction) TransactionResponse {
	userID := tan.UserID.String()
	resp := TransactionResponse{
		ID:               tan.ID.String(),
		Amount:           tan.Amount,
		Currency:         tan.Currency.String(),
		Type:             tan.Type.String(),
		UserID:           &userID,
		CreatedTimestamp: tan.CreatedTimestamp,
	}
	if tan.Reference != "" {
		ref := &tan.Reference
		resp.Reference = ref
	}
	return resp
}

type ListTransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions" validate:"required"`
}

type CreateUserRequest struct {
	Name        string  `json:"name" validate:"required"`
	Address     Address `json:"address" validate:"required"`
	PhoneNumber string  `json:"phoneNumber" validate:"required,phone"`
	Email       string  `json:"email" validate:"required,email"`
}

func (r CreateUserRequest) toDomain() (users.CreateUserRequest, error) {
	name := r.Name
	address, err := r.Address.toDomain()
	if err != nil {
		return users.CreateUserRequest{}, err
	}
	number, err := users.NewPhoneNumber(r.PhoneNumber)
	if err != nil {
		return users.CreateUserRequest{}, err
	}
	email, err := users.NewEmail(r.Email)
	if err != nil {
		return users.CreateUserRequest{}, err
	}
	return users.NewCreateUserRequest(name, address, number, email)
}

type UpdateUserRequest struct {
	Name        *string  `json:"name,omitempty"`
	Address     *Address `json:"address,omitempty"`
	PhoneNumber *string  `json:"phoneNumber,omitempty" validate:"omitempty,phone"`
	Email       *string  `json:"email,omitempty" validate:"omitempty,email"`
}

type UserResponse struct {
	ID               string    `json:"id" validate:"required,userID"`
	Name             string    `json:"name" validate:"required"`
	Address          Address   `json:"address" validate:"required"`
	PhoneNumber      string    `json:"phoneNumber" validate:"required,phone"`
	Email            string    `json:"email" validate:"required,email"`
	CreatedTimestamp time.Time `json:"createdTimestamp" validate:"required"`
	UpdatedTimestamp time.Time `json:"updatedTimestamp" validate:"required"`
}

func newUserResponseFromDomain(user users.User) UserResponse {
	return UserResponse{
		ID:               user.ID.String(),
		Name:             user.Name,
		Address:          newAddressFromDomain(user.Address),
		PhoneNumber:      user.PhoneNumber.String(),
		Email:            user.Email.String(),
		CreatedTimestamp: user.Created,
		UpdatedTimestamp: user.Updated,
	}
}

type ErrorResponse struct {
	Message string `json:"message" validate:"required"`
}

func newErrorResponse(err error) ErrorResponse {
	resp := ErrorResponse{Message: err.Error()}
	err = validation.Get().Struct(resp)
	if err != nil {
		resp.Message = "unspecified error"
	}
	return resp
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

func newBadRequestErrorResponse(err error) BadRequestErrorResponse {
	var details []ValidationDetail
	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		for _, e := range validateErrs {
			details = append(details, ValidationDetail{
				Field:   e.Field(),
				Message: e.Error(),
				Type:    e.Type().String(),
			})
		}
	}
	resp := BadRequestErrorResponse{
		Message: err.Error(),
		Details: details,
	}
	err = validation.Get().Struct(resp)
	if err != nil {
		if resp.Message == "" {
			resp.Message = "validation error"
		}
		if resp.Details == nil {
			resp.Details = []ValidationDetail{
				{
					Field:   "unknown",
					Message: "unknown",
					Type:    "unknown",
				},
			}
		}
	}
	return resp
}

type LoginRequest struct {
	UserID       string `json:"userID" validate:"required,userID"`
	PasswordHash string `json:"passwordhash" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token" validate:"required"`
}
