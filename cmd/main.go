package main

import (
	"agricultural-equipment-store/internal/config"
	"agricultural-equipment-store/internal/delivery/http"
	"agricultural-equipment-store/internal/infrastructure/database"
	"agricultural-equipment-store/internal/infrastructure/logger"
	"agricultural-equipment-store/internal/repository"
	"agricultural-equipment-store/internal/usecase"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "agricultural-equipment-store/docs" // Import docs for Swagger
)

// @title Agricultural Equipment Store API
// @version 1.0
// @description API for Agricultural Equipment Store with Clean Architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@agricultural-equipment-store.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8082
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.NewLogger()

	// Initialize database
	db, err := database.NewMongoDB(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	saleRepo := repository.NewSaleRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWT.Secret)
	productUseCase := usecase.NewProductUseCase(productRepo)
	inventoryUseCase := usecase.NewInventoryUseCase(productRepo)
	saleUseCase := usecase.NewSaleUseCase(saleRepo, productRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)

	// Initialize HTTP server
	server := http.NewServer(cfg, logger, authUseCase, productUseCase, inventoryUseCase, saleUseCase, categoryUseCase)

	// Start server
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Shutdown server
	server.Shutdown()
}
