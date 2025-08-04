package main

import (
	"eaglebank/internal/users"
	"eaglebank/internal/users/adapters"
	"eaglebank/internal/validation"
	"eaglebank/internal/web"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	validate := validation.MustNewValidator()

	usrStore := adapters.NewInMemoryUserStore()
	usrSvc := users.NewUserService(usrStore)

	srv := web.NewServer(logger, validate, usrSvc)

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
