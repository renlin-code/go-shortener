package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/renlin-code/go-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/renlin-code/go-shortener/internal/http-server/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/renlin-code/go-shortener/internal/config"
	"github.com/renlin-code/go-shortener/internal/lib/logger/handlers/slogpretty"
	sl "github.com/renlin-code/go-shortener/internal/lib/logger/slog"
	"github.com/renlin-code/go-shortener/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages enabled")

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/api/v1/url", save.NewHandler(log, storage))

	log.Info("starting server...", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Timeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(1)
	}

	log.Error("server stopped")
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
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
