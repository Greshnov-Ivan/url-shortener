// @title URL Shortener API
// @version 1.0
// @description REST API для сервиса сокращения ссылок.
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/internal/app"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config: cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// init logger: log/slog
	log := setupLogger(cfg.Env)

	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init app
	application, err := app.NewApp(cfg, log)
	if err != nil {
		log.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}

	// run server
	errChan := make(chan error, 1)
	go func() {
		errChan <- application.Run()
	}()

	log.Info("server started", slog.String("address", cfg.Address))

	// processing completion signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-done:
		log.Info("stopping server...")
	case err := <-errChan:
		log.Error("server crashed", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("stopping server...")

	// context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.App.GracefulShutdownTimeout)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.Any("error", err))
		return
	}

	log.Info("server stopped")
}

// configuring the logger
func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return setupPrettySlog()
	case envDev:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

// pretty logger for the local environment
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
