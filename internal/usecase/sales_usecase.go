package usecase

import (
	"agricultural-equipment-store/internal/domain"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// SalesUseCase handles sales related business logic
type SalesUseCase struct {
	saleRepo    domain.SaleRepository
	productRepo domain.ProductRepository
}

// NewSalesUseCase creates a new sales use case
func NewSalesUseCase(saleRepo domain.SaleRepository, productRepo domain.ProductRepository) *SalesUseCase {
	return &SalesUseCase{
		saleRepo:    saleRepo,
		productRepo: productRepo,
	}
}

// CreateSale creates a new sale
func (u *SalesUseCase) CreateSale(ctx context.Context, req domain.CreateSaleRequest) (*domain.Sale, error) {
	// Get product to validate and calculate total
	product, err := u.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Check stock availability
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Calculate total
	total := product.Price * float64(req.Quantity)

	// Create sale
	sale := &domain.Sale{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Price:     product.Price,
		Total:     total,
		DateSold:  time.Now(),
	}

	// Save sale
	err = u.saleRepo.Create(ctx, sale)
	if err != nil {
		return nil, err
	}

	// Update product stock
	product.Stock -= req.Quantity
	err = u.productRepo.Update(ctx, product)
	if err != nil {
		// TODO: Consider implementing transaction rollback
		return nil, err
	}

	return sale, nil
}

// GetSales retrieves sales with filtering
func (u *SalesUseCase) GetSales(ctx context.Context, filter domain.SaleFilter) ([]*domain.Sale, error) {
	return u.saleRepo.List(ctx, filter)
}

// GetSalesSummary retrieves sales summary for a period
func (u *SalesUseCase) GetSalesSummary(ctx context.Context, fromDate, toDate time.Time) (*domain.SalesSummary, error) {
	return u.saleRepo.GetSalesSummary(ctx, fromDate, toDate)
}

// GetSalesByProduct retrieves sales grouped by product
func (u *SalesUseCase) GetSalesByProduct(ctx context.Context, fromDate, toDate time.Time) ([]*domain.ProductSales, error) {
	return u.saleRepo.GetSalesByProduct(ctx, fromDate, toDate)
}

// ExportSales exports sales data in specified format
func (u *SalesUseCase) ExportSales(ctx context.Context, fromDate, toDate time.Time, format string) (string, error) {
	if format != "csv" {
		return "", errors.New("unsupported format")
	}

	// Get sales data
	sales, err := u.saleRepo.GetSalesByDateRange(ctx, fromDate, toDate)
	if err != nil {
		return "", err
	}

	// Generate CSV
	var csvBuilder strings.Builder

	// Header
	csvBuilder.WriteString("ID,Product ID,Product Name,Quantity,Price,Total,Date Sold\n")

	// Data rows
	for _, sale := range sales {
		// Get product name
		product, _ := u.productRepo.GetByID(ctx, sale.ProductID)
		productName := "Unknown"
		if product != nil {
			productName = product.Name
		}

		csvBuilder.WriteString(fmt.Sprintf("%s,%s,%s,%d,%.2f,%.2f,%s\n",
			sale.ID.Hex(),
			sale.ProductID.Hex(),
			productName,
			sale.Quantity,
			sale.Price,
			sale.Total,
			sale.DateSold.Format("2006-01-02 15:04:05"),
		))
	}

	return csvBuilder.String(), nil
}
