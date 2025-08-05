package accounts

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"regexp"
	"time"
)

type AccountType string

const PersonalAcct AccountType = "personal"

func (a AccountType) IsValid() bool {
	switch a {
	case PersonalAcct:
		return true
	default:
		return false
	}
}

func (a AccountType) String() string {
	return string(a)
}

func NewAccountType(s string) (AccountType, error) {
	acctType := AccountType(s)
	if !acctType.IsValid() {
		return "", fmt.Errorf("invalid account type %q", s)
	}
	return acctType, nil
}

type AccountNumber string

var accountNumberRegex = regexp.MustCompile(`^01\d{6}$`)

func (a AccountNumber) IsValid() bool {
	return accountNumberRegex.MatchString(string(a))
}

func (a AccountNumber) String() string {
	return string(a)
}

func NewAccountNumber(s string) (AccountNumber, error) {
	acct := AccountNumber(s)
	if !acct.IsValid() {
		return "", fmt.Errorf("invalid account number %q: must match format 01XXXXXX", s)
	}
	return acct, nil
}

func NewRandAccountNumber() (AccountNumber, error) {
	return NewAccountNumber(fmt.Sprintf("01%06d", rand.IntN(1000000)))
}

type SortCode string

func (c SortCode) IsValid() bool {
	return c == "10-10-10"
}

func (c SortCode) String() string {
	return string(c)
}

func NewSortCode(s string) (SortCode, error) {
	code := SortCode(s)
	if !code.IsValid() {
		return "", fmt.Errorf("invalid sort code %q: must be 10-10-10", s)
	}
	return code, nil
}

type Currency string

const GBP Currency = "GBP"

func (c Currency) IsValid() bool {
	switch c {
	case GBP:
		return true
	default:
		return false
	}
}

func (c Currency) String() string {
	return string(c)
}

func NewCurrency(s string) (Currency, error) {
	currency := Currency(s)
	if !currency.IsValid() {
		return "", fmt.Errorf("invalid currency %q", s)
	}
	return currency, nil
}

const balanceMax float64 = 10000
const balanceMin float64 = 0

type BankAccount struct {
	AccountNumber    AccountNumber
	SortCode         SortCode
	Name             string
	AccountType      AccountType
	balance          float64
	Currency         Currency
	CreatedTimestamp time.Time
	UpdatedTimestamp time.Time
}

func (ba *BankAccount) IsValid() bool {
	if ba.Name == "" {
		return false
	}
	if ba.balance < balanceMin || ba.balance > balanceMax {
		return false
	}
	if !ba.AccountNumber.IsValid() {
		return false
	}
	if !ba.SortCode.IsValid() {
		return false
	}
	if !ba.AccountType.IsValid() {
		return false
	}
	if !ba.Currency.IsValid() {
		return false
	}
	return true
}

func (ba *BankAccount) Balance() float64 {
	return ba.balance
}

func (ba *BankAccount) Withdraw(amt float64) error {
	newBalance := ba.balance - amt
	if newBalance < balanceMin {
		return errors.New("account overdrawn")
	}
	return nil
}

func (ba *BankAccount) Deposit(amt float64) error {
	newBalance := ba.balance + amt
	if newBalance > balanceMax {
		return errors.New("account underdrawn")
	}
	return nil
}

func NewBankAccount(acctNum AccountNumber, sortCode SortCode, name string, acctType AccountType, curr Currency) (*BankAccount, error) {
	now := time.Now()
	acct := BankAccount{
		AccountNumber:    acctNum,
		SortCode:         sortCode,
		Name:             name,
		AccountType:      acctType,
		balance:          0,
		Currency:         curr,
		CreatedTimestamp: now,
		UpdatedTimestamp: now,
	}
	if !acct.IsValid() {
		return &BankAccount{}, fmt.Errorf("invalid bank account %+v", acct)
	}
	return &acct, nil
}

type CreateAccountRequest struct {
	Name        string
	AccountType AccountType
}

func (r CreateAccountRequest) IsValid() bool {
	if r.Name == "" {
		return false
	}
	if !r.AccountType.IsValid() {
		return false
	}
	return true
}

func NewCreateAccountRequest(name string, acctType AccountType) (CreateAccountRequest, error) {
	req := CreateAccountRequest{
		Name:        name,
		AccountType: acctType,
	}
	if !req.IsValid() {
		return CreateAccountRequest{}, fmt.Errorf("invalid create account request %+v", req)
	}
	return req, nil
}
