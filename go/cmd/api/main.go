package main

import (
	"eaglebank/internal/accounts"
	adapters2 "eaglebank/internal/accounts/adapters"
	"eaglebank/internal/transactions"
	adapters3 "eaglebank/internal/transactions/adapters"
	"eaglebank/internal/users"
	"eaglebank/internal/users/adapters"
	"eaglebank/internal/web"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	usrStore := adapters.NewInMemoryUserStore()
	usrSvc := users.NewUserService(usrStore)

	acctStore := adapters2.NewInMemoryAccountStore()
	acctSvc := accounts.NewAccountService(acctStore)

	tanStore := adapters3.NewInMemoryTransactionStore()
	tanSvc := transactions.NewTransactionService(tanStore, acctStore)

	srv := web.NewServer(web.ServerArgs{
		Logger:  logger,
		UserSvc: usrSvc,
		AcctSvc: acctSvc,
		TanSvc:  tanSvc,
	})

	port := "8080"
	logger.Info("Starting Eagle Bank api, serving on :" + port)
	s := &http.Server{
		Addr:         ":" + port,
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		logger.Error(fmt.Errorf("fatal error in server: %v", err).Error())
	}
}
