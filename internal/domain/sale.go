package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sale represents a sale transaction
type Sale struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Product   *Product           `json:"product,omitempty" bson:"product,omitempty"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Price     float64            `json:"price" bson:"price"`
	Total     float64            `json:"total" bson:"total"`
	DateSold  time.Time          `json:"date_sold" bson:"date_sold"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// CreateSaleRequest represents the request payload for creating a sale
type CreateSaleRequest struct {
	ProductID primitive.ObjectID `json:"product_id" binding:"required"`
	Quantity  int                `json:"quantity" binding:"required,gt=0"`
	Price     float64            `json:"price" binding:"required,gt=0"`
}

// SaleFilter represents filter options for sales
type SaleFilter struct {
	ProductID primitive.ObjectID `json:"product_id"`
	FromDate  time.Time          `json:"from_date"`
	ToDate    time.Time          `json:"to_date"`
	Page      int                `json:"page"`
	Limit     int                `json:"limit"`
}

// SalesSummary represents sales summary data
type SalesSummary struct {
	TotalSales   float64 `json:"total_sales"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalItems   int     `json:"total_items"`
	Period       string  `json:"period"`
}

// ProductSales represents sales data for a specific product
type ProductSales struct {
	ProductID    primitive.ObjectID `json:"product_id"`
	ProductName  string             `json:"product_name"`
	TotalSold    int                `json:"total_sold"`
	TotalRevenue float64            `json:"total_revenue"`
}

// StockUpdateRequest represents the request payload for updating stock
type StockUpdateRequest struct {
	Stock int `json:"stock" binding:"required,gte=0"`
}

// StockSummary represents stock summary data
type StockSummary struct {
	TotalProducts    int             `json:"total_products"`
	TotalStockValue  float64         `json:"total_stock_value"`
	LowStockProducts int             `json:"low_stock_products"`
	Categories       []CategoryStock `json:"categories"`
}

// CategoryStock represents stock data for a category
type CategoryStock struct {
	Category     string  `json:"category"`
	TotalStock   int     `json:"total_stock"`
	TotalValue   float64 `json:"total_value"`
	ProductCount int     `json:"product_count"`
}

// LowStockProduct represents a product with low stock
type LowStockProduct struct {
	ID       primitive.ObjectID `json:"id"`
	Name     string             `json:"name"`
	Stock    int                `json:"stock"`
	Category string             `json:"category"`
	Price    float64            `json:"price"`
}
