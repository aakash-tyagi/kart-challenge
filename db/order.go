package db

import (
	"context"

	"github.com/aakash-tyagi/kart-challenge/models"
)

func (db *Db) SaveOrder(ctx context.Context, order *models.Order) error {

	order.DefaultCreateAt()
	order.DefaultUpdateAt()

	_, err := db.MongoClient.Collection("orders").InsertOne(ctx, order)
	if err != nil {

		return err
	}
	return nil
}
