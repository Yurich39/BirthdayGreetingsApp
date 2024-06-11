package app

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/slog"

	"github.com/go-chi/chi/v5"

	"BirthdayGreetingsService/config"
	"BirthdayGreetingsService/internal/controller/api"
	"BirthdayGreetingsService/internal/usecase"
	"BirthdayGreetingsService/internal/usecase/repo"
	"BirthdayGreetingsService/pkg/httpserver"
	"BirthdayGreetingsService/pkg/logger"
	"BirthdayGreetingsService/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Env)

	// Repository
	db, err := postgres.New(cfg.StorageConfig)
	if err != nil {
		l.Debug("failed to init storage", l.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	// Creating usecase for employees
	employeesUseCase := usecase.NewEmployees(
		repo.NewEmployeesRepo(db),
		l,
	)

	// Creating usecase for one-to-many relationship between employee and employee
	subscribersUseCase := usecase.NewSubscribers(
		repo.NewSubscriberRepo(db),
		l,
	)

	// HTTP Server
	r := chi.NewRouter()
	api.NewRouter(cfg, r, l, employeesUseCase, subscribersUseCase)

	l.Info("starting server", slog.String("address", cfg.Address))

	httpServer := httpserver.New(r, cfg)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Debug("Failed to start server", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Debug("Server shutdown", err)
	}
}
