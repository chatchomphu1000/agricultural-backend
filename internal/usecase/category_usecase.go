package usecase

import (
	"agricultural-equipment-store/internal/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryUseCase handles category business logic
type CategoryUseCase struct {
	categoryRepo domain.CategoryRepository
}

// NewCategoryUseCase creates a new category use case
func NewCategoryUseCase(categoryRepo domain.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory creates a new category
func (u *CategoryUseCase) CreateCategory(ctx context.Context, req domain.CreateCategoryRequest) (*domain.Category, error) {
	// Check if category already exists
	existing, err := u.categoryRepo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category already exists")
	}

	category := &domain.Category{
		Name: req.Name,
	}

	err = u.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategories retrieves all categories
func (u *CategoryUseCase) GetCategories(ctx context.Context) ([]*domain.Category, error) {
	return u.categoryRepo.List(ctx)
}

// GetCategoryByID retrieves a category by ID
func (u *CategoryUseCase) GetCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	objID, err := parseObjectID(id)
	if err != nil {
		return nil, err
	}

	category, err := u.categoryRepo.GetByID(ctx, objID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	return category, nil
}

// DeleteCategory deletes a category
func (u *CategoryUseCase) DeleteCategory(ctx context.Context, id string) error {
	objID, err := parseObjectID(id)
	if err != nil {
		return err
	}

	// Check if category exists
	category, err := u.categoryRepo.GetByID(ctx, objID)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	return u.categoryRepo.Delete(ctx, objID)
}

// parseObjectID parses string to ObjectID
func parseObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
