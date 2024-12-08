package db

import (
	"context"
	"fmt"

	"github.com/aakash-tyagi/kart-challenge/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *Db) AddProduct(ctx context.Context, product models.Product) error {

	product.DefaultCreateAt()
	product.DefaultUpdateAt()

	// Insert the product into the database
	product.Id = primitive.NewObjectID()
	_, err := db.MongoClient.Collection("products").InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to add product: %v", err)
	}
	return nil
}

func (db *Db) ListProducts(ctx context.Context, page int, limit int) ([]models.Product, int64, error) {
	// Calculate the number of documents to skip
	skip := (page - 1) * limit

	// Create a slice to hold the products
	var products []models.Product

	// Query the database with pagination
	cursor, err := db.MongoClient.Collection("products").Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each product
	for cursor.Next(ctx) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, 0, err
		}

		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	// Count total number of entries
	total, err := db.MongoClient.Collection("products").CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (db *Db) GetProductById(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	// Create a variable to hold the product
	var product models.Product

	// Query the database for the product by ID
	err := db.MongoClient.Collection("products").FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		// If no document is found, return nil and the error
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no product found with id: %s", id)
		}
		return nil, err
	}

	return &product, nil
}
