package repository

import (
	"agricultural-equipment-store/internal/domain"
	"agricultural-equipment-store/internal/infrastructure/database"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// categoryRepository implements domain.CategoryRepository
type categoryRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *database.MongoDB) domain.CategoryRepository {
	return &categoryRepository{
		db:         db,
		collection: db.GetCollection("categories"),
	}
}

// Create creates a new category
func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	// Check if category name already exists
	existing, err := r.GetByName(ctx, category.Name)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("category already exists")
	}

	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, category)
	return err
}

// GetByID retrieves a category by ID
func (r *categoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Category, error) {
	var category domain.Category
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// GetByName retrieves a category by name
func (r *categoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	var category domain.Category
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// List retrieves all categories
func (r *categoryRepository) List(ctx context.Context) ([]*domain.Category, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*domain.Category
	for cursor.Next(ctx) {
		var category domain.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, cursor.Err()
}

// Update updates a category
func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	category.UpdatedAt = time.Now()

	filter := bson.M{"_id": category.ID}
	update := bson.M{"$set": category}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a category
func (r *categoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
