package models

type Product struct {
	Model    `bson:",inline"`
	Name     string `bson:"name" json:"name"`
	Price    int32  `bson:"price" json:"proce"`
	Category string `bson:"category" json:"category"`
}
