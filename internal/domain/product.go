package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product in the system
type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Category    string             `json:"category" bson:"category"`
	Brand       string             `json:"brand" bson:"brand"`
	ImageURL    string             `json:"image_url" bson:"image_url"` // Legacy field for backward compatibility
	Images      []ProductImage     `json:"images" bson:"images"`       // New field for multiple images
	Stock       int                `json:"stock" bson:"stock"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// ProductImage represents an image associated with a product
type ProductImage struct {
	ID        string    `json:"id" bson:"id"`                 // Unique ID for this image
	URL       string    `json:"url" bson:"url"`               // Image URL (for URL-based images)
	Filename  string    `json:"filename" bson:"filename"`     // Original filename (for uploaded files)
	FilePath  string    `json:"file_path" bson:"file_path"`   // Server file path (for uploaded files)
	FileSize  int64     `json:"file_size" bson:"file_size"`   // File size in bytes
	MimeType  string    `json:"mime_type" bson:"mime_type"`   // MIME type (image/jpeg, image/png, etc.)
	IsURL     bool      `json:"is_url" bson:"is_url"`         // true if URL-based, false if uploaded file
	IsPrimary bool      `json:"is_primary" bson:"is_primary"` // true for the main product image
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// Category represents a product category
type Category struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// CreateCategoryRequest represents the request payload for creating a category
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateProductRequest represents the request payload for creating a product
type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Category    string   `json:"category" binding:"required"`
	Brand       string   `json:"brand"`
	ImageURL    string   `json:"image_url"`  // Legacy field for backward compatibility
	ImageURLs   []string `json:"image_urls"` // Multiple image URLs
	Stock       int      `json:"stock" binding:"required,gte=0"`
}

// UpdateProductRequest represents the request payload for updating a product
type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Category    string   `json:"category"`
	Brand       string   `json:"brand"`
	ImageURL    string   `json:"image_url"`  // Legacy field for backward compatibility
	ImageURLs   []string `json:"image_urls"` // Multiple image URLs
	Stock       int      `json:"stock"`
	IsActive    *bool    `json:"is_active"`
}

// ProductFilter represents filter options for products
type ProductFilter struct {
	Category string  `json:"category"`
	Brand    string  `json:"brand"`
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
	IsActive *bool   `json:"is_active"`
	Search   string  `json:"search"`
	Page     int     `json:"page"`
	Limit    int     `json:"limit"`
}
