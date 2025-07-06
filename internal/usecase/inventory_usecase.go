package usecase

import (
	"agricultural-equipment-store/internal/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InventoryUseCase handles inventory related business logic
type InventoryUseCase struct {
	productRepo domain.ProductRepository
}

// NewInventoryUseCase creates a new inventory use case
func NewInventoryUseCase(productRepo domain.ProductRepository) *InventoryUseCase {
	return &InventoryUseCase{
		productRepo: productRepo,
	}
}

// UpdateStock updates the stock for a product
func (u *InventoryUseCase) UpdateStock(ctx context.Context, id primitive.ObjectID, req domain.StockUpdateRequest) error {
	// Check if product exists
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	// Update stock
	return u.productRepo.UpdateStock(ctx, id, req.Stock)
}

// GetLowStockProducts retrieves products with low stock
func (u *InventoryUseCase) GetLowStockProducts(ctx context.Context, threshold int) ([]*domain.LowStockProduct, error) {
	if threshold <= 0 {
		threshold = 10 // Default threshold
	}

	return u.productRepo.GetLowStockProducts(ctx, threshold)
}

// GetStockSummary retrieves stock summary
func (u *InventoryUseCase) GetStockSummary(ctx context.Context) (*domain.StockSummary, error) {
	return u.productRepo.GetStockSummary(ctx)
}

// SaleUseCase handles sales related business logic
type SaleUseCase struct {
	saleRepo    domain.SaleRepository
	productRepo domain.ProductRepository
}

// NewSaleUseCase creates a new sale use case
func NewSaleUseCase(saleRepo domain.SaleRepository, productRepo domain.ProductRepository) *SaleUseCase {
	return &SaleUseCase{
		saleRepo:    saleRepo,
		productRepo: productRepo,
	}
}

// CreateSale creates a new sale and updates product stock
func (u *SaleUseCase) CreateSale(ctx context.Context, req domain.CreateSaleRequest) (*domain.Sale, error) {
	// Get product to verify it exists and has enough stock
	product, err := u.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Check if there's enough stock
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Calculate total
	total := req.Price * float64(req.Quantity)

	// Create sale
	sale := &domain.Sale{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Total:     total,
		DateSold:  time.Now(),
	}

	err = u.saleRepo.Create(ctx, sale)
	if err != nil {
		return nil, err
	}

	// Update product stock
	newStock := product.Stock - req.Quantity
	err = u.productRepo.UpdateStock(ctx, req.ProductID, newStock)
	if err != nil {
		return nil, err
	}

	return sale, nil
}

// GetSalesByFilter retrieves sales with filtering
func (u *SaleUseCase) GetSalesByFilter(ctx context.Context, filter domain.SaleFilter) ([]*domain.Sale, error) {
	return u.saleRepo.List(ctx, filter)
}

// GetSalesSummary retrieves sales summary for a period
func (u *SaleUseCase) GetSalesSummary(ctx context.Context, fromDate, toDate time.Time) (*domain.SalesSummary, error) {
	// If no dates provided, use current month
	if fromDate.IsZero() || toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		toDate = fromDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	return u.saleRepo.GetSalesSummary(ctx, fromDate, toDate)
}

// GetSalesByProduct retrieves sales grouped by product
func (u *SaleUseCase) GetSalesByProduct(ctx context.Context, fromDate, toDate time.Time) ([]*domain.ProductSales, error) {
	// If no dates provided, use current month
	if fromDate.IsZero() || toDate.IsZero() {
		now := time.Now()
		fromDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		toDate = fromDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	return u.saleRepo.GetSalesByProduct(ctx, fromDate, toDate)
}

// GetSalesByDateRange retrieves sales for a specific date range
func (u *SaleUseCase) GetSalesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]*domain.Sale, error) {
	return u.saleRepo.GetSalesByDateRange(ctx, fromDate, toDate)
}
