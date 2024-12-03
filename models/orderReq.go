package models

type OrderReq struct {
	Model      `bson:",inline"`
	CouponCode string `bson:"couponCode" json:"couponCode"`
	Items      []Item `bson:"items" json:"items"`
}
