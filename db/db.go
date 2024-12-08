package db

import (
	"context"

	"github.com/aakash-tyagi/kart-challenge/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	Config      *config.Config
	MongoClient *mongo.Database
}

func New(config *config.Config) *Db {
	return &Db{
		Config: config,
	}
}

func (db *Db) Connect() error {
	clientOptions := options.Client().ApplyURI(db.Config.DBUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	mongoClient := client.Database(db.Config.DatabaseName)
	db.MongoClient = mongoClient
	return nil
}
