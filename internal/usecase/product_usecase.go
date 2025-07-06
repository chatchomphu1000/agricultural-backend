package usecase

import (
	"agricultural-equipment-store/internal/domain"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductUseCase handles product related business logic
type ProductUseCase struct {
	productRepo domain.ProductRepository
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo domain.ProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (u *ProductUseCase) CreateProduct(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Brand:       req.Brand,
		ImageURL:    req.ImageURL,
		Stock:       req.Stock,
		IsActive:    true,
	}

	// Handle multiple image URLs if provided
	if len(req.ImageURLs) > 0 {
		for i, url := range req.ImageURLs {
			if url != "" {
				product.Images = append(product.Images, domain.ProductImage{
					ID:        uuid.New().String(),
					URL:       url,
					IsURL:     true,
					IsPrimary: i == 0, // First image is primary
					CreatedAt: time.Now(),
				})
			}
		}
	}

	// Handle legacy single image URL for backward compatibility
	if req.ImageURL != "" && len(product.Images) == 0 {
		product.Images = append(product.Images, domain.ProductImage{
			ID:        uuid.New().String(),
			URL:       req.ImageURL,
			IsURL:     true,
			IsPrimary: true,
			CreatedAt: time.Now(),
		})
	}

	err := u.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// CreateProductWithImages creates a new product with both uploaded images and image URLs
func (u *ProductUseCase) CreateProductWithImages(ctx context.Context, req domain.CreateProductRequest, uploadedImages []domain.ProductImage) (*domain.Product, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Brand:       req.Brand,
		ImageURL:    req.ImageURL, // Keep for backward compatibility
		Stock:       req.Stock,
		IsActive:    true,
	}

	// Add uploaded images first
	product.Images = append(product.Images, uploadedImages...)

	// Add image URLs
	if len(req.ImageURLs) > 0 {
		for _, url := range req.ImageURLs {
			if url != "" {
				product.Images = append(product.Images, domain.ProductImage{
					ID:        uuid.New().String(),
					URL:       url,
					IsURL:     true,
					IsPrimary: len(product.Images) == 0, // First image is primary
					CreatedAt: time.Now(),
				})
			}
		}
	}

	// Handle legacy single image URL for backward compatibility
	if req.ImageURL != "" && len(product.Images) == 0 {
		product.Images = append(product.Images, domain.ProductImage{
			ID:        uuid.New().String(),
			URL:       req.ImageURL,
			IsURL:     true,
			IsPrimary: true,
			CreatedAt: time.Now(),
		})
	}

	// Ensure at least one image is marked as primary
	if len(product.Images) > 0 {
		hasPrimary := false
		for _, img := range product.Images {
			if img.IsPrimary {
				hasPrimary = true
				break
			}
		}
		if !hasPrimary {
			product.Images[0].IsPrimary = true
		}
	}

	err := u.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (u *ProductUseCase) GetProductByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// UpdateProduct updates a product
func (u *ProductUseCase) UpdateProduct(ctx context.Context, id primitive.ObjectID, req domain.UpdateProductRequest) (*domain.Product, error) {
	// Get existing product
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Brand != "" {
		product.Brand = req.Brand
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
		// Also update images array for backward compatibility
		hasLegacyImage := false
		for i, img := range product.Images {
			if img.IsURL && img.URL == product.ImageURL {
				hasLegacyImage = true
				break
			}
			if img.IsURL && img.IsPrimary {
				product.Images[i].URL = req.ImageURL
				hasLegacyImage = true
				break
			}
		}
		if !hasLegacyImage {
			product.Images = append([]domain.ProductImage{
				{
					ID:        uuid.New().String(),
					URL:       req.ImageURL,
					IsURL:     true,
					IsPrimary: true,
					CreatedAt: time.Now(),
				},
			}, product.Images...)
		}
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	// Handle multiple image URLs if provided
	if len(req.ImageURLs) > 0 {
		// Remove existing URL-based images
		var newImages []domain.ProductImage
		for _, img := range product.Images {
			if !img.IsURL {
				newImages = append(newImages, img)
			}
		}

		// Add new URL-based images
		for i, url := range req.ImageURLs {
			if url != "" {
				newImages = append(newImages, domain.ProductImage{
					ID:        uuid.New().String(),
					URL:       url,
					IsURL:     true,
					IsPrimary: i == 0 && len(newImages) == 0, // First image is primary if no uploaded images
					CreatedAt: time.Now(),
				})
			}
		}

		product.Images = newImages
	}

	err = u.productRepo.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProductWithImages updates a product with both uploaded images and image URLs
func (u *ProductUseCase) UpdateProductWithImages(ctx context.Context, id primitive.ObjectID, req domain.UpdateProductRequest, uploadedImages []domain.ProductImage) (*domain.Product, error) {
	// Get existing product
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Update basic fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Brand != "" {
		product.Brand = req.Brand
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	// Handle images - replace all existing images with new ones
	var newImages []domain.ProductImage

	// Add uploaded images first
	newImages = append(newImages, uploadedImages...)

	// Add image URLs
	if len(req.ImageURLs) > 0 {
		for _, url := range req.ImageURLs {
			if url != "" {
				newImages = append(newImages, domain.ProductImage{
					ID:        uuid.New().String(),
					URL:       url,
					IsURL:     true,
					IsPrimary: len(newImages) == 0, // First image is primary
					CreatedAt: time.Now(),
				})
			}
		}
	}

	// Handle legacy single image URL for backward compatibility
	if req.ImageURL != "" && len(newImages) == 0 {
		newImages = append(newImages, domain.ProductImage{
			ID:        uuid.New().String(),
			URL:       req.ImageURL,
			IsURL:     true,
			IsPrimary: true,
			CreatedAt: time.Now(),
		})
		product.ImageURL = req.ImageURL
	}

	// Update images array
	product.Images = newImages

	// Ensure at least one image is marked as primary
	if len(product.Images) > 0 {
		hasPrimary := false
		for _, img := range product.Images {
			if img.IsPrimary {
				hasPrimary = true
				break
			}
		}
		if !hasPrimary {
			product.Images[0].IsPrimary = true
		}
	}

	err = u.productRepo.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product
func (u *ProductUseCase) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	// Check if product exists
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	return u.productRepo.Delete(ctx, id)
}

// GetProducts retrieves products with filtering and pagination
func (u *ProductUseCase) GetProducts(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	// Set default pagination values
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	// Get products
	products, err := u.productRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := u.productRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return products, count, nil
}
