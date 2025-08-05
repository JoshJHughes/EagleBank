package web

import "eaglebank/internal/users"

type UserService interface {
	CreateUser(req users.CreateUserRequest) (users.User, error)
	GetUser(userID users.UserID) (users.User, error)
}
