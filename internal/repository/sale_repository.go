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

// saleRepository implements domain.SaleRepository
type saleRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

// NewSaleRepository creates a new sale repository
func NewSaleRepository(db *database.MongoDB) domain.SaleRepository {
	return &saleRepository{
		db:         db,
		collection: db.GetCollection("sales"),
	}
}

// Create creates a new sale
func (r *saleRepository) Create(ctx context.Context, sale *domain.Sale) error {
	sale.ID = primitive.NewObjectID()
	sale.CreatedAt = time.Now()
	sale.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, sale)
	return err
}

// GetByID retrieves a sale by ID
func (r *saleRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Sale, error) {
	var sale domain.Sale
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sale)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &sale, nil
}

// List retrieves a list of sales with filtering and pagination
func (r *saleRepository) List(ctx context.Context, filter domain.SaleFilter) ([]*domain.Sale, error) {
	// Build MongoDB filter
	mongoFilter := bson.M{}

	if !filter.ProductID.IsZero() {
		mongoFilter["product_id"] = filter.ProductID
	}

	if !filter.FromDate.IsZero() && !filter.ToDate.IsZero() {
		mongoFilter["date_sold"] = bson.M{
			"$gte": filter.FromDate,
			"$lte": filter.ToDate,
		}
	} else if !filter.FromDate.IsZero() {
		mongoFilter["date_sold"] = bson.M{"$gte": filter.FromDate}
	} else if !filter.ToDate.IsZero() {
		mongoFilter["date_sold"] = bson.M{"$lte": filter.ToDate}
	}

	// Set up pagination
	page := filter.Page
	limit := filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	skip := (page - 1) * limit

	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{"date_sold", -1}})

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []*domain.Sale
	for cursor.Next(ctx) {
		var sale domain.Sale
		if err := cursor.Decode(&sale); err != nil {
			return nil, err
		}
		sales = append(sales, &sale)
	}

	return sales, cursor.Err()
}

// Count counts sales with filtering
func (r *saleRepository) Count(ctx context.Context, filter domain.SaleFilter) (int64, error) {
	mongoFilter := bson.M{}

	if !filter.ProductID.IsZero() {
		mongoFilter["product_id"] = filter.ProductID
	}

	if !filter.FromDate.IsZero() && !filter.ToDate.IsZero() {
		mongoFilter["date_sold"] = bson.M{
			"$gte": filter.FromDate,
			"$lte": filter.ToDate,
		}
	} else if !filter.FromDate.IsZero() {
		mongoFilter["date_sold"] = bson.M{"$gte": filter.FromDate}
	} else if !filter.ToDate.IsZero() {
		mongoFilter["date_sold"] = bson.M{"$lte": filter.ToDate}
	}

	return r.collection.CountDocuments(ctx, mongoFilter)
}

// GetSalesSummary retrieves sales summary for a date range
func (r *saleRepository) GetSalesSummary(ctx context.Context, fromDate, toDate time.Time) (*domain.SalesSummary, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"date_sold": bson.M{
					"$gte": fromDate,
					"$lte": toDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":           nil,
				"total_sales":   bson.M{"$sum": "$total"},
				"total_revenue": bson.M{"$sum": "$total"},
				"total_items":   bson.M{"$sum": "$quantity"},
				"count":         bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		TotalSales   float64 `bson:"total_sales"`
		TotalRevenue float64 `bson:"total_revenue"`
		TotalItems   int     `bson:"total_items"`
		Count        int     `bson:"count"`
	}

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	period := fromDate.Format("2006-01-02") + " to " + toDate.Format("2006-01-02")

	return &domain.SalesSummary{
		TotalSales:   result.TotalSales,
		TotalRevenue: result.TotalRevenue,
		TotalItems:   result.TotalItems,
		Period:       period,
	}, nil
}

// GetSalesByProduct retrieves sales grouped by product
func (r *saleRepository) GetSalesByProduct(ctx context.Context, fromDate, toDate time.Time) ([]*domain.ProductSales, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"date_sold": bson.M{
					"$gte": fromDate,
					"$lte": toDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":           "$product_id",
				"total_sold":    bson.M{"$sum": "$quantity"},
				"total_revenue": bson.M{"$sum": "$total"},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "products",
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "product",
			},
		},
		{
			"$unwind": "$product",
		},
		{
			"$project": bson.M{
				"product_id":    "$_id",
				"product_name":  "$product.name",
				"total_sold":    1,
				"total_revenue": 1,
			},
		},
		{
			"$sort": bson.M{"total_revenue": -1},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var productSales []*domain.ProductSales
	for cursor.Next(ctx) {
		var ps domain.ProductSales
		if err := cursor.Decode(&ps); err != nil {
			return nil, err
		}
		productSales = append(productSales, &ps)
	}

	return productSales, cursor.Err()
}

// GetSalesByDateRange retrieves sales within a date range
func (r *saleRepository) GetSalesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]*domain.Sale, error) {
	filter := bson.M{
		"date_sold": bson.M{
			"$gte": fromDate,
			"$lte": toDate,
		},
	}

	opts := options.Find()
	opts.SetSort(bson.D{{"date_sold", -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sales []*domain.Sale
	for cursor.Next(ctx) {
		var sale domain.Sale
		if err := cursor.Decode(&sale); err != nil {
			return nil, err
		}
		sales = append(sales, &sale)
	}

	return sales, cursor.Err()
}
