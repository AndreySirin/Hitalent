package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"hitalent/internal/logger"
	"hitalent/internal/server"
	"hitalent/internal/storage"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	dbName   = "postgres"
	user     = "postgres"
	password = "secret"
	address  = "postgres:5432"
	port     = ":8080"
)

func main() {
	lg := logger.New()
	lg.Info("Starting server...")

	validate := validator.New()

	lg.Info("Connecting to database", "user", user, "db", dbName, "address", address)
	db, err := storage.New(user, password, address, dbName)
	if err != nil {
		lg.Error("Failed to connect to database", "error", err)
		return
	}
	lg.Info("Database connection established")

	lg.Info("Applying migrations")
	err = storage.MigrateUP(db)
	if err != nil {
		lg.Error("Failed to apply migrations", "error", err)
		return
	}
	lg.Info("Migrations applied successfully")

	srv := server.New(lg, validate, port, db)
	lg.Info("Server initialized", "port", port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		lg.Info("Running server")
		err = srv.Run()
		if err != nil {
			lg.Error("Failed to start server", "error", err)
			return
		}
	}()

	go func() {
		<-stop
		lg.Info("Shutdown signal received, stopping server...")
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err = srv.Shutdown(ctx)
		if err != nil {
			lg.Error("Error during server shutdown", "error", err)
		} else {
			lg.Info("Server stopped gracefully")
		}
		wg.Done()
	}()

	wg.Wait()
	lg.Info("Main function exit")
}
