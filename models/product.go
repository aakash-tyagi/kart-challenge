package models

import (
	"errors"
)

type Product struct {
	Model    `bson:",inline"`
	Name     string `bson:"name" json:"name"`
	Price    int32  `bson:"price" json:"price"`
	Category string `bson:"category" json:"category"`
}

type ListProduct struct {
	Products      []Product `json:"products"`
	TotalProducts int       `json:"totalProducts"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if p.Category == "" {
		return errors.New("category is required")
	}
	return nil
}
