package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model defines the default fields to handle when operation happens
// import the Model in document struct to make it working
type Model struct {
	Id       primitive.ObjectID `bson:"_id"      json:"id"`
	CreateAt time.Time          `bson:"createAt" json:"createAt"`
	UpdateAt time.Time          `bson:"updateAt" json:"updateAt"`
}

// DefaultUpdateAt changes the default updateAt field
func (df *Model) DefaultUpdateAt() {
	df.UpdateAt = time.Now().UTC()
}

// DefaultCreateAt changes the default createAt field
func (df *Model) DefaultCreateAt() {
	if df.CreateAt.IsZero() {
		df.CreateAt = time.Now().UTC()
	}
}

// DefaultId changes the default _id field
func (df *Model) DefaultId() {
	if df.Id.IsZero() {
		df.Id = primitive.NewObjectID()
	}
}
