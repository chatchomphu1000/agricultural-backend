package main

import (
	"agricultural-equipment-store/internal/config"
	"agricultural-equipment-store/internal/domain"
	"agricultural-equipment-store/internal/infrastructure/database"
	"agricultural-equipment-store/internal/repository"
	"agricultural-equipment-store/internal/usecase"
	"context"
	"log"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewMongoDB(cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create indexes
	if err := db.CreateIndexes(); err != nil {
		log.Fatal("Failed to create indexes:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWT.Secret)
	productUseCase := usecase.NewProductUseCase(productRepo)

	ctx := context.Background()

	// Create admin user
	adminReq := domain.CreateUserRequest{
		Email:    cfg.Admin.Email,
		Password: cfg.Admin.Password,
		Name:     "Administrator",
		Role:     "admin",
	}

	existingAdmin, err := userRepo.GetByEmail(ctx, adminReq.Email)
	if err != nil {
		log.Fatal("Failed to check existing admin:", err)
	}

	if existingAdmin == nil {
		_, err = authUseCase.Register(ctx, adminReq)
		if err != nil {
			log.Fatal("Failed to create admin user:", err)
		}
		log.Println("Admin user created successfully")
	} else {
		log.Println("Admin user already exists")
	}

	// Create sample products
	sampleProducts := []domain.CreateProductRequest{
		{
			Name:        "John Deere X350 Lawn Tractor",
			Description: "42-inch cutting deck, 17.5 HP engine, comfortable seat",
			Price:       2499.99,
			Category:    "Lawn Mowers",
			Brand:       "John Deere",
			ImageURL:    "https://example.com/images/john-deere-x350.jpg",
			Stock:       15,
		},
		{
			Name:        "Husqvarna 450 Chainsaw",
			Description: "18-inch bar, 50.2cc engine, professional grade",
			Price:       329.99,
			Category:    "Chainsaws",
			Brand:       "Husqvarna",
			ImageURL:    "https://example.com/images/husqvarna-450.jpg",
			Stock:       25,
		},
		{
			Name:        "Kubota BX23S Compact Tractor",
			Description: "23 HP diesel engine, 4WD, backhoe attachment",
			Price:       28999.99,
			Category:    "Tractors",
			Brand:       "Kubota",
			ImageURL:    "https://example.com/images/kubota-bx23s.jpg",
			Stock:       8,
		},
		{
			Name:        "STIHL MS 170 Chainsaw",
			Description: "16-inch bar, 30.1cc engine, lightweight design",
			Price:       179.99,
			Category:    "Chainsaws",
			Brand:       "STIHL",
			ImageURL:    "https://example.com/images/stihl-ms170.jpg",
			Stock:       30,
		},
		{
			Name:        "Troy-Bilt Pony 42 Riding Mower",
			Description: "42-inch cutting deck, 17.5 HP engine, automatic transmission",
			Price:       1299.99,
			Category:    "Lawn Mowers",
			Brand:       "Troy-Bilt",
			ImageURL:    "https://example.com/images/troy-bilt-pony42.jpg",
			Stock:       12,
		},
	}

	for _, productReq := range sampleProducts {
		_, err = productUseCase.CreateProduct(ctx, productReq)
		if err != nil {
			log.Printf("Failed to create product %s: %v", productReq.Name, err)
		} else {
			log.Printf("Product created: %s", productReq.Name)
		}
	}

	log.Println("Data seeding completed successfully!")
}
