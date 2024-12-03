package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	Model    `bson:",inline"`
	Items    []Item    `bson:"itmes" json:"items"`
	Products []Product `bson:"products" json:"products"`
}

type Item struct {
	ProductId primitive.ObjectID `bson:"productId" json:"productId"`
	Quantity  int8               `bson:"quantity" json:"quantity"`
}
