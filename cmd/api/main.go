package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PongponZ/test-container-golang/cmd/api/handler"
	"github.com/PongponZ/test-container-golang/config"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.New()

	// Create a channel to receive signals for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	indexHandler := handler.NewIndex()
	todoHandler := handler.NewTodo()

	server := echo.New()
	server.GET("/", indexHandler.Index)
	server.POST("/todo", todoHandler.Create)

	// start server
	go func() {
		if err := server.Start(cfg.Port); err != nil {
			log.Fatalf("shutting down the server: %v\n", err)
		}
	}()

	// wait for SIGINT
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
	log.Println("Server shutdown gracefully")
}
