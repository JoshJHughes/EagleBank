package accounts

import "errors"

var ErrAccountNotFound = errors.New("account not found")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrTooManyFunds = errors.New("you have too much money")
