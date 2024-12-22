package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/gommon/log"
	"quote-management-tech-task/cmd"
	"quote-management-tech-task/config"
	"quote-management-tech-task/db/sqlc"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var appConfig config.Config
	if err := env.Parse(&appConfig); err != nil {
		log.Fatal("failed to parse common environment variables", "err", err)
	}

	db, err := pgxpool.New(context.Background(), appConfig.DbURL)
	if err != nil {
		log.Fatal("failed to create new database connection pool", "err", err)
	}

	if err = db.Ping(context.Background()); err != nil {
		log.Fatal("failed to ping database", "err", err)
	}

	dbQuerier := sqlc.New(db)
	defer db.Close()

	app := cmd.NewServer(appConfig, dbQuerier)

	// Start app
	go func() {
		if err = app.Run(); err != nil {
			log.Fatal("failed to run app", "err", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	const defaultWaitTimeForGracefulShutdown = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), defaultWaitTimeForGracefulShutdown)
	defer cancel()
	if err = app.Shutdown(ctx); err != nil {
		log.Info("failed to shutdown app", "err", err)
	}
}
