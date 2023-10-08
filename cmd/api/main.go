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
	"github.com/PongponZ/test-container-golang/pkg/infra/database/mongodb"
	"github.com/PongponZ/test-container-golang/pkg/repository"
	"github.com/PongponZ/test-container-golang/pkg/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.New()

	// Create a channel to receive signals for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Create a context for the server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient := mongodb.NewConnect(ctx, cfg.MongoURI)
	defer mongoClient.Disconnect(ctx)

	db := mongoClient.GetDatabase(cfg.Database)

	todoRepo := repository.NewToDo(db)
	todoUsecase := usecase.NewToDo(todoRepo)
	todoHandler := handler.NewTodo(todoUsecase)

	server := echo.New()
	server.Use(middleware.Logger())
	server.GET("/", handler.Index)
	server.GET("/todo", todoHandler.Gets)
	server.POST("/todo", todoHandler.Create)
	server.PUT("/todo", todoHandler.Update)
	server.DELETE("/todo/:id", todoHandler.Delete)

	// start server
	go func() {
		if err := server.Start(cfg.Port); err != nil {
			log.Fatalf("shutting down the server: %v\n", err)
		}
	}()

	<-sig // wait for SIGINT

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
	log.Println("Server shutdown gracefully")
}
