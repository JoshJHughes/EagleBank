package web

import (
	"eaglebank/internal/accounts"
	"eaglebank/internal/users"
)

type UserService interface {
	CreateUser(req users.CreateUserRequest) (users.User, error)
	GetUser(userID users.UserID) (users.User, error)
}

type AccountService interface {
	CreateAccount(req accounts.CreateAccountRequest) (*accounts.BankAccount, error)
	ListAccounts(id users.UserID) ([]accounts.BankAccount, error)
}
