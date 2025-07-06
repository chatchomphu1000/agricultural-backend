package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents MongoDB connection
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB successfully!")

	return &MongoDB{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// GetDatabase returns the database instance
func (m *MongoDB) GetDatabase() *mongo.Database {
	return m.database
}

// GetCollection returns a collection instance
func (m *MongoDB) GetCollection(collectionName string) *mongo.Collection {
	return m.database.Collection(collectionName)
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.client.Disconnect(ctx)
}

// CreateIndexes creates necessary indexes for the application
func (m *MongoDB) CreateIndexes() error {
	// Create unique index for user email
	userCollection := m.GetCollection("users")
	userIndexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCollection.Indexes().CreateOne(ctx, userIndexModel)
	if err != nil {
		return err
	}

	// Create indexes for products
	productCollection := m.GetCollection("products")
	productIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"name", "text"}, {"description", "text"}},
		},
		{
			Keys: bson.D{{"category", 1}},
		},
		{
			Keys: bson.D{{"brand", 1}},
		},
		{
			Keys: bson.D{{"price", 1}},
		},
	}

	_, err = productCollection.Indexes().CreateMany(ctx, productIndexes)
	if err != nil {
		return err
	}

	log.Println("Database indexes created successfully!")
	return nil
}
