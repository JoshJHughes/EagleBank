# EagleBank

## Submission Note
I have written the api in Go, my recruiter indicated that this was fine on the phone, I apologise if that isn't the case

## Installation
- download and install go
- navigate to `/go` folder under repo root
- install dependencies `go mod download`
- run `go run ./cmd/api/main.go`

## Run tests
`go test ./...`

## Routes implemented
`GET /health`

`POST /login`

`POST /v1/users`


`GET /v1/users/{userId}`


`POST /v1/accounts`

`GET /v1/accounts`

`GET /v1/accounts/{accountNumber}`


`POST /v1/accounts/{accountNumber}/transactions`

`GET /v1/accounts/{accountNumber}/transactions`

`GET /v1/accounts/{accountNumber}/transactions/{transactionId}`


## Architecture overview
- 3 services: users, accounts, transactions
- 3 layers:
  - web for authentication, authorisation, validation, and parsing
  - application for business logic
  - infrastructure for stores
- Each layer tested using TDD
- Hexagonal architecture with ports/adapter model


## Technical decision & tradeoffs
- Current implementation of the account and balance stores & updates is vulnerable to partial failures & race conditions
  - Ideally the store updates would each be atomic
  - For the in-memory solution could expose the lock for both stores and update while both locks are taken to avoid race conditions but not necessarily keep the values in-sync
  - For real systems could fix potential inconsistencies by making transactions the source-of-truth and reconciling the account balance out-of-band if the account update fails


- I included logic to make the transaction store read-only in the store itself but in hindsight it may be better for the store to have update and delete methods to facilitate admin tasks and move the read-only logic into the application later
- I put the account balance update logic in the transactions service to avoid writing another handler in the account service. It might be desirable to separate responsibility for adding transactions and reconciling the account balance, but it seemed unnecessary here.


- I handled authentication by sending a hashed password with the http request, auth would probably be better done using a 3rd party service in prod.
- I also hard-coded the jwt secret key, which is clearly bad practice and I would not do so in a real system 
- I chose to use single global logger and to not abstract it behind an interface for simplicity and to declutter function signatures. In a larger project it may be worth constructing an interface and passing it down through the context. 
- I have also used a single global validator. I experimented using a validator for domain type validation in the users package but in hindsight I preferred to set up my own validation rules within the object constructors as it seems easier to follow, breaks the coupling between web and domain layers, and is more idiomatic in Go.
- I have used custom errors only where I needed to for the test
- I have tended to use exported members for dataclasses, as is idiomatic in Go. I have made type constructors that validate their arguments for use in application functions but exported members make for easier testing.
- I have propagated data by value as I used in-memory stores and passing by-reference would introduce potential bugs by bypassing the store methods.
- As a note I changed the transaction ID regex to allow more than one character after 'tan-', which I assume was a typo in the spec
