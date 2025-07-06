package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, page, limit int) ([]*User, error)
}

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, filter ProductFilter) ([]*Product, error)
	Count(ctx context.Context, filter ProductFilter) (int64, error)

	// Stock management methods
	UpdateStock(ctx context.Context, id primitive.ObjectID, stock int) error
	GetLowStockProducts(ctx context.Context, threshold int) ([]*LowStockProduct, error)
	GetStockSummary(ctx context.Context) (*StockSummary, error)
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Category, error)
	GetByName(ctx context.Context, name string) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// SaleRepository defines the interface for sale data operations
type SaleRepository interface {
	Create(ctx context.Context, sale *Sale) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Sale, error)
	List(ctx context.Context, filter SaleFilter) ([]*Sale, error)
	Count(ctx context.Context, filter SaleFilter) (int64, error)

	// Sales analytics methods
	GetSalesSummary(ctx context.Context, fromDate, toDate time.Time) (*SalesSummary, error)
	GetSalesByProduct(ctx context.Context, fromDate, toDate time.Time) ([]*ProductSales, error)
	GetSalesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]*Sale, error)
}
