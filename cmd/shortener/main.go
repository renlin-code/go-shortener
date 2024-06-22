package main

import (
	"log/slog"
	"os"

	"github.com/renlin-code/go-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	//TODO: init logger (slog)

	log := setupLogger(cfg.Env)

	log.Info("starting shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages enabled")

	//TODO: init storage (SQL Lite)

	//TODO: init router (chi)

	//TODO: run server
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
