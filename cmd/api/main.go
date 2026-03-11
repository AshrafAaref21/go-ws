package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AshrafAaref21/go-ws/internal/config"
	"github.com/AshrafAaref21/go-ws/internal/db"
	"github.com/AshrafAaref21/go-ws/internal/realtime"
	"github.com/AshrafAaref21/go-ws/internal/routes"
	"github.com/AshrafAaref21/go-ws/internal/utils"
)

func main() {
	cfg := config.LoadConfig()
	utils.InitJWT(cfg.JWTKey)

	// mux.HandleFunc("GET /api/health-check-http", handlers.HandleHealthCheckHTTP)

	err := db.InitDB(cfg.DBPath, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		err := db.CloseDB()
		if err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	hub := realtime.NewHub()
	defer hub.Shutdown()

	mux := routes.RegisterRoutes(hub)

	server := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: mux,
	}

	log.Printf("Server is running on %s.", cfg.HTTPServer.Address)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-shutdownSignal.Done()
	log.Println("Shutdown signal received.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server failed to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped.")
}
