package repository

import (
	"agricultural-equipment-store/internal/domain"
	"agricultural-equipment-store/internal/infrastructure/database"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// productRepository implements domain.ProductRepository
type productRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *database.MongoDB) domain.ProductRepository {
	return &productRepository{
		db:         db,
		collection: db.GetCollection("products"),
	}
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, product)
	return err
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	var product domain.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// Update updates a product
func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	product.UpdatedAt = time.Now()

	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": product}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a product
func (r *productRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List retrieves a list of products with filtering and pagination
func (r *productRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, error) {
	// Build MongoDB filter
	mongoFilter := bson.M{}

	if filter.Category != "" {
		mongoFilter["category"] = filter.Category
	}
	if filter.Brand != "" {
		mongoFilter["brand"] = filter.Brand
	}
	if filter.MinPrice > 0 || filter.MaxPrice > 0 {
		priceFilter := bson.M{}
		if filter.MinPrice > 0 {
			priceFilter["$gte"] = filter.MinPrice
		}
		if filter.MaxPrice > 0 {
			priceFilter["$lte"] = filter.MaxPrice
		}
		mongoFilter["price"] = priceFilter
	}
	if filter.IsActive != nil {
		mongoFilter["is_active"] = *filter.IsActive
	}
	if filter.Search != "" {
		// Search primarily in name only for more precise results
		mongoFilter["name"] = bson.M{"$regex": filter.Search, "$options": "i"}
	}

	// Build options
	opts := options.Find()

	// Pagination
	if filter.Page > 0 && filter.Limit > 0 {
		skip := (filter.Page - 1) * filter.Limit
		opts.SetSkip(int64(skip))
		opts.SetLimit(int64(filter.Limit))
	}

	// Sort by creation date (newest first)
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	for cursor.Next(ctx) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, cursor.Err()
}

// Count returns the total count of products matching the filter
func (r *productRepository) Count(ctx context.Context, filter domain.ProductFilter) (int64, error) {
	// Build MongoDB filter (same as List method)
	mongoFilter := bson.M{}

	if filter.Category != "" {
		mongoFilter["category"] = filter.Category
	}
	if filter.Brand != "" {
		mongoFilter["brand"] = filter.Brand
	}
	if filter.MinPrice > 0 || filter.MaxPrice > 0 {
		priceFilter := bson.M{}
		if filter.MinPrice > 0 {
			priceFilter["$gte"] = filter.MinPrice
		}
		if filter.MaxPrice > 0 {
			priceFilter["$lte"] = filter.MaxPrice
		}
		mongoFilter["price"] = priceFilter
	}
	if filter.IsActive != nil {
		mongoFilter["is_active"] = *filter.IsActive
	}
	if filter.Search != "" {
		// Search primarily in name only for more precise results
		mongoFilter["name"] = bson.M{"$regex": filter.Search, "$options": "i"}
	}

	return r.collection.CountDocuments(ctx, mongoFilter)
}

// UpdateStock updates the stock quantity for a product
func (r *productRepository) UpdateStock(ctx context.Context, id primitive.ObjectID, stock int) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"stock":      stock,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// GetLowStockProducts retrieves products with stock below the threshold
func (r *productRepository) GetLowStockProducts(ctx context.Context, threshold int) ([]*domain.LowStockProduct, error) {
	filter := bson.M{
		"stock":     bson.M{"$lt": threshold},
		"is_active": true,
	}

	opts := options.Find()
	opts.SetSort(bson.D{{"stock", 1}}) // Sort by stock ascending (lowest first)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var lowStockProducts []*domain.LowStockProduct
	for cursor.Next(ctx) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}

		lowStockProduct := &domain.LowStockProduct{
			ID:       product.ID,
			Name:     product.Name,
			Stock:    product.Stock,
			Category: product.Category,
			Price:    product.Price,
		}
		lowStockProducts = append(lowStockProducts, lowStockProduct)
	}

	return lowStockProducts, cursor.Err()
}

// GetStockSummary retrieves stock summary data
func (r *productRepository) GetStockSummary(ctx context.Context) (*domain.StockSummary, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"is_active": true},
		},
		{
			"$group": bson.M{
				"_id":           "$category",
				"total_stock":   bson.M{"$sum": "$stock"},
				"total_value":   bson.M{"$sum": bson.M{"$multiply": []interface{}{"$stock", "$price"}}},
				"product_count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []domain.CategoryStock
	var totalProducts int
	var totalStockValue float64

	for cursor.Next(ctx) {
		var categoryStock domain.CategoryStock
		if err := cursor.Decode(&categoryStock); err != nil {
			return nil, err
		}
		categories = append(categories, categoryStock)
		totalProducts += categoryStock.ProductCount
		totalStockValue += categoryStock.TotalValue
	}

	// Get low stock products count
	lowStockCount, err := r.collection.CountDocuments(ctx, bson.M{
		"stock":     bson.M{"$lt": 10},
		"is_active": true,
	})
	if err != nil {
		return nil, err
	}

	return &domain.StockSummary{
		TotalProducts:    totalProducts,
		TotalStockValue:  totalStockValue,
		LowStockProducts: int(lowStockCount),
		Categories:       categories,
	}, nil
}
