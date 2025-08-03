package users

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type UserID string

var userIDRegex = regexp.MustCompile(`^usr-[A-Za-z0-9]+$`)

func NewUserID(s string) (UserID, error) {
	if s == "" {
		return "", fmt.Errorf("user ID cannot be empty")
	}
	if !userIDRegex.MatchString(s) {
		return "", fmt.Errorf("invalid user ID format: %q, must match pattern %q", s, userIDRegex.String())
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
	if s == "" {
		return "", errors.New("email cannot be empty")
	}

	_, err := mail.ParseAddress(s)
	if err != nil {
		return "", fmt.Errorf("invalid email format: %q", s)
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

var phoneNumberRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

func NewPhoneNumber(s string) (PhoneNumber, error) {
	if s == "" {
		return "", fmt.Errorf("phone number cannot be empty")
	}
	if !phoneNumberRegex.MatchString(s) {
		return "", fmt.Errorf("invalid phone number format: %q, must match pattern %q", s, phoneNumberRegex.String())
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
	Line1    string
	Town     string
	County   string
	Postcode string

	// optional
	Line2 string
	Line3 string
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
	if line1 == "" {
		return Address{}, errors.New("address line 1 is required")
	}
	if town == "" {
		return Address{}, errors.New("town is required")
	}
	if county == "" {
		return Address{}, errors.New("county is required")
	}
	if postcode == "" {
		return Address{}, errors.New("postcode is required")
	}

	addr := Address{
		Line1:    line1,
		Town:     town,
		County:   county,
		Postcode: postcode,
	}

	for _, opt := range opts {
		opt(&addr)
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
	ID          UserID
	Name        string
	Address     Address
	PhoneNumber PhoneNumber
	Email       Email
	Created     time.Time
	Updated     time.Time
}

func NewUser(id UserID, name string, address Address, phone PhoneNumber, email Email) (User, error) {
	if name == "" {
		return User{}, errors.New("name is required")
	}

	return User{
		ID:          id,
		Name:        name,
		Address:     address,
		PhoneNumber: phone,
		Email:       email,
		Created:     time.Now(),
		Updated:     time.Now(),
	}, nil
}

func MustNewUser(id UserID, name string, address Address, phone PhoneNumber, email Email) User {
	usr, err := NewUser(id, name, address, phone, email)
	if err != nil {
		panic(fmt.Sprintf("MustNewUser: %v", err))
	}
	return usr
}

type CreateUserRequest struct {
	Name        string
	Address     Address
	PhoneNumber PhoneNumber
	Email       Email
}

func NewCreateUserRequest(name string, address Address, number PhoneNumber, email Email) (CreateUserRequest, error) {
	if name == "" {
		return CreateUserRequest{}, errors.New("name is required")
	}

	return CreateUserRequest{
		Name:        name,
		Address:     address,
		PhoneNumber: number,
		Email:       email,
	}, nil
}

func MustNewCreateUserRequest(name string, address Address, number PhoneNumber, email Email) CreateUserRequest {
	req, err := NewCreateUserRequest(name, address, number, email)
	if err != nil {
		panic(fmt.Sprintf("MustNewUser: %v", err))
	}

	return req
}
