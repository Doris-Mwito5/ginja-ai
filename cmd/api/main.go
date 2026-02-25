package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Doris-Mwito5/ginja-ai/internal/configs"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/domain"
	"github.com/Doris-Mwito5/ginja-ai/internal/logger"
	"github.com/Doris-Mwito5/ginja-ai/web/routes"
)

func main() {

	configs.InitializeEnvironment()

	log.Printf("Console: %s", configs.Config.Environment)
	
	logger.InitLogger("ginja-ai")
	// Init db
	dB := db.InitDB(configs.Config.DatabaseURL)
	defer dB.Close()

	logger.Info("Starting ginja-ai")


	//domain store
	domainStore := domain.NewStore()

	appRouter := routes.BuildRouter(
		dB,
		domainStore,
	)

	server := &http.Server{
		Addr:    ":" + configs.Config.Port,
		Handler: appRouter,
	}

	done := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		logger.Info("Process terminated...shutting down")

		if err := server.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Server shut down error: %v", err)
		}

		close(done)
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Info("Server shut down")
		} else {
			logger.Fatal("Server shut down unexpectedly!")
		}
	}

	timeout := 30 * time.Second
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	code := 0
	select {
	case <-sigint:
		code = 1
		logger.Info("Process forcibly terminated")
	case <-time.After(timeout):
		code = 1
		logger.Info("Shutdown timeout. Forcibly shutting down...")
	case <-done:
		logger.Info("Shutdown completed...")
	}

	logger.Info("Server exiting...")

	os.Exit(code)

}
