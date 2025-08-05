package users

import (
	"eaglebank/internal/validation"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type UserID string

func NewUserID(s string) (UserID, error) {
	err := validation.Get().Var(s, "required,userID")
	if err != nil {
		return "", err
	}
	return UserID(s), nil
}

func MustNewUserID(s string) UserID {
	id, err := NewUserID(s)
	if err != nil {
		panic(fmt.Sprintf("MustNewUserID: %v", err))
	}
	return id
}

func NewRandUserID() (UserID, error) {
	id := uuid.New()
	clean := strings.ReplaceAll(id.String(), "-", "")
	return NewUserID("usr-" + clean)
}

func MustNewRandUserID() UserID {
	id, err := NewRandUserID()
	if err != nil {
		panic(fmt.Sprintf("MustNewRandUserID: %v", err))
	}
	return id
}

func (u UserID) String() string {
	return string(u)
}

type Email string

func NewEmail(s string) (Email, error) {
	err := validation.Get().Var(s, "required,email")
	if err != nil {
		return "", err
	}
	return Email(s), nil
}

func MustNewEmail(s string) Email {
	email, err := NewEmail(s)
	if err != nil {
		panic(fmt.Sprintf("MustNewEmail: %v", err))
	}
	return email
}

func (e Email) String() string {
	return string(e)
}

type PhoneNumber string

func NewPhoneNumber(s string) (PhoneNumber, error) {
	err := validation.Get().Var(s, "required,phone")
	if err != nil {
		return "", err
	}
	return PhoneNumber(s), nil
}

func MustNewPhoneNumber(s string) PhoneNumber {
	id, err := NewPhoneNumber(s)
	if err != nil {
		panic(fmt.Sprintf("MustNewPhoneNumber: %v", err))
	}
	return id
}

func (n PhoneNumber) String() string {
	return string(n)
}

type Address struct {
	Line1    string `validate:"required"`
	Line2    string
	Line3    string
	Town     string `validate:"required"`
	County   string `validate:"required"`
	Postcode string `validate:"required"`
}

type AddressOption func(*Address)

func WithLine2(line2 string) AddressOption {
	return func(a *Address) {
		a.Line2 = line2
	}
}

func WithLine3(line3 string) AddressOption {
	return func(a *Address) {
		a.Line3 = line3
	}
}

func NewAddress(line1, town, county, postcode string, opts ...AddressOption) (Address, error) {
	addr := Address{
		Line1:    line1,
		Town:     town,
		County:   county,
		Postcode: postcode,
	}
	for _, opt := range opts {
		opt(&addr)
	}
	err := validation.Get().Struct(addr)
	if err != nil {
		return Address{}, err
	}
	return addr, nil
}

func MustNewAddress(line1, town, county, postcode string, opts ...AddressOption) Address {
	addr, err := NewAddress(line1, town, county, postcode, opts...)
	if err != nil {
		panic(fmt.Sprintf("MustNewAddress: %v", err))
	}
	return addr
}

type User struct {
	ID          UserID      `validate:"required,userID"`
	Name        string      `validate:"required"`
	Address     Address     `validate:"required"`
	PhoneNumber PhoneNumber `validate:"required,phone"`
	Email       Email       `validate:"required,email"`
	Created     time.Time   `validate:"required"`
	Updated     time.Time   `validate:"required"`
}

func NewUser(id UserID, name string, address Address, phone PhoneNumber, email Email) (User, error) {
	now := time.Now()
	usr := User{
		ID:          id,
		Name:        name,
		Address:     address,
		PhoneNumber: phone,
		Email:       email,
		Created:     now,
		Updated:     now,
	}
	err := validation.Get().Struct(usr)
	if err != nil {
		return User{}, err
	}
	return usr, nil
}

func MustNewUser(id UserID, name string, address Address, phone PhoneNumber, email Email) User {
	usr, err := NewUser(id, name, address, phone, email)
	if err != nil {
		panic(fmt.Sprintf("MustNewUser: %v", err))
	}
	return usr
}

type CreateUserRequest struct {
	Name        string      `validate:"required"`
	Address     Address     `validate:"required"`
	PhoneNumber PhoneNumber `validate:"required,phone"`
	Email       Email       `validate:"required,email"`
}

func NewCreateUserRequest(name string, address Address, number PhoneNumber, email Email) (CreateUserRequest, error) {
	req := CreateUserRequest{
		Name:        name,
		Address:     address,
		PhoneNumber: number,
		Email:       email,
	}
	err := validation.Get().Struct(req)
	if err != nil {
		return CreateUserRequest{}, err
	}
	return req, nil
}

func MustNewCreateUserRequest(name string, address Address, number PhoneNumber, email Email) CreateUserRequest {
	req, err := NewCreateUserRequest(name, address, number, email)
	if err != nil {
		panic(fmt.Sprintf("MustNewUser: %v", err))
	}

	return req
}
