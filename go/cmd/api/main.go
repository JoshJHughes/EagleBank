package main

import (
	"eaglebank/internal/web"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := web.NewServer(logger)

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
