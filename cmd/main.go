package main

import (
	"SongsLibrary/internal/app"
	"SongsLibrary/internal/config"
	"SongsLibrary/pkg/logger/sl"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// @title Test songs library
// @version 1.0

// @host localhost:8082
// @BasePath /
func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Debug("failed to load config due to error: %w", sl.Err(err))
	}

	app := app.New(context.Background(), log, cfg)
	go func() {
		app.Start()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Stop()
	log.Info("Gracefully stopped")
}
