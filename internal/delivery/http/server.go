package http

import (
	"agricultural-equipment-store/internal/config"
	"agricultural-equipment-store/internal/delivery/http/middleware"
	"agricultural-equipment-store/internal/infrastructure/logger"
	"agricultural-equipment-store/internal/usecase"
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the HTTP server
type Server struct {
	config           *config.Config
	logger           logger.Logger
	authUseCase      *usecase.AuthUseCase
	productUseCase   *usecase.ProductUseCase
	inventoryUseCase *usecase.InventoryUseCase
	saleUseCase      *usecase.SaleUseCase
	categoryUseCase  *usecase.CategoryUseCase
	server           *http.Server
}

// NewServer creates a new HTTP server
func NewServer(
	config *config.Config,
	logger logger.Logger,
	authUseCase *usecase.AuthUseCase,
	productUseCase *usecase.ProductUseCase,
	inventoryUseCase *usecase.InventoryUseCase,
	saleUseCase *usecase.SaleUseCase,
	categoryUseCase *usecase.CategoryUseCase,
) *Server {
	return &Server{
		config:           config,
		logger:           logger,
		authUseCase:      authUseCase,
		productUseCase:   productUseCase,
		inventoryUseCase: inventoryUseCase,
		saleUseCase:      saleUseCase,
		categoryUseCase:  categoryUseCase,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{s.config.Frontend.URL}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Initialize handlers
	authHandler := NewAuthHandler(s.authUseCase)
	productHandler := NewProductHandler(s.productUseCase)
	inventoryHandler := NewInventoryHandler(s.inventoryUseCase)
	saleHandler := NewSaleHandler(s.saleUseCase)
	categoryHandler := NewCategoryHandler(s.categoryUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(s.authUseCase)

	// Setup routes
	s.setupRoutes(router, authHandler, productHandler, inventoryHandler, saleHandler, categoryHandler, authMiddleware)

	// Create HTTP server
	s.server = &http.Server{
		Addr:         ":" + s.config.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	s.logger.Info("Starting server on port %s", s.config.Server.Port)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// setupRoutes sets up all API routes
func (s *Server) setupRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	productHandler *ProductHandler,
	inventoryHandler *InventoryHandler,
	saleHandler *SaleHandler,
	categoryHandler *CategoryHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Agricultural Equipment Store API is running",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Static file server for uploaded images
	router.Static("/uploads", "./uploads")

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", authMiddleware.RequireAuth(), authHandler.GetProfile)
		}

		// Product routes
		products := api.Group("/products")
		{
			// Public routes
			products.GET("", productHandler.GetProducts)    // Get all products (public)
			products.GET("/:id", productHandler.GetProduct) // Get single product (public)

			// Admin routes
			products.POST("", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), productHandler.CreateProduct)
			products.PUT("/:id", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), productHandler.UpdateProduct)
			products.DELETE("/:id", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), productHandler.DeleteProduct)
		}

		// Inventory routes
		inventories := api.Group("/inventories")
		{
			inventories.PUT("/:id/stock", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), inventoryHandler.UpdateStock)
			inventories.GET("/low-stock", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), inventoryHandler.GetLowStockProducts)
			inventories.GET("/summary", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), inventoryHandler.GetStockSummary)
		}

		// Sales routes
		sales := api.Group("/sales")
		{
			sales.POST("", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), saleHandler.CreateSale)
			sales.GET("", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), saleHandler.GetSales)
			sales.GET("/summary", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), saleHandler.GetSalesSummary)
			sales.GET("/by-product", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), saleHandler.GetSalesByProduct)
			sales.GET("/export", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), saleHandler.ExportSales)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			// Public routes
			categories.GET("", categoryHandler.GetCategories)   // Get all categories (public)
			categories.GET("/:id", categoryHandler.GetCategory) // Get single category (public)

			// Admin routes
			categories.POST("", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), categoryHandler.CreateCategory)
			categories.DELETE("/:id", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin(), categoryHandler.DeleteCategory)
		}
	}

	// Print routes for debugging
	s.logger.Info("API Routes:")
	s.logger.Info("POST   /api/auth/register")
	s.logger.Info("POST   /api/auth/login")
	s.logger.Info("GET    /api/auth/profile")
	s.logger.Info("GET    /api/products")
	s.logger.Info("GET    /api/products/:id")
	s.logger.Info("POST   /api/products (admin)")
	s.logger.Info("PUT    /api/products/:id (admin)")
	s.logger.Info("DELETE /api/products/:id (admin)")
	s.logger.Info("PUT    /api/inventories/:id/stock (admin)")
	s.logger.Info("GET    /api/inventories/low-stock (admin)")
	s.logger.Info("GET    /api/inventories/summary (admin)")
	s.logger.Info("POST   /api/sales (admin)")
	s.logger.Info("GET    /api/sales (admin)")
	s.logger.Info("GET    /api/sales/summary (admin)")
	s.logger.Info("GET    /api/sales/by-product (admin)")
	s.logger.Info("GET    /api/sales/export (admin)")
	s.logger.Info("GET    /api/categories")
	s.logger.Info("GET    /api/categories/:id")
	s.logger.Info("POST   /api/categories (admin)")
	s.logger.Info("DELETE /api/categories/:id (admin)")
	s.logger.Info("GET    /swagger/index.html")
}
