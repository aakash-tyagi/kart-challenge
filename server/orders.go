package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aakash-tyagi/kart-challenge/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	orderReq := &models.OrderReq{}

	// Unmarshal the request body into the req variable
	if err := json.NewDecoder(r.Body).Decode(orderReq); err != nil {
		writeJSONResponse(w, http.StatusUnprocessableEntity, "Error decoding request body: "+err.Error())
		return
	}

	order := &models.Order{}

	// verify coupon code
	if orderReq.CouponCode != "" {
		err := s.validateCoupon(ctx, orderReq.CouponCode, 32)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, "Invalid Coupon Code: "+err.Error())
			return
		}
		order.IsCouponCodeValid = true
	}

	// Fetch products by their IDs and add to the order
	products := make([]models.Product, 0) // Assuming you have a Product model
	for _, item := range orderReq.Items {
		product, err := s.DBClient.GetProductById(ctx, item.ProductId) // Fetch product by ID
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, "Product not found")
			return
		}
		products = append(products, *product) // Dereference product to append
	}

	order = &models.Order{
		Model: models.Model{
			Id: primitive.NewObjectID(),
		},
		Items:             orderReq.Items,
		Products:          products,
		IsCouponCodeValid: order.IsCouponCodeValid,
	}

	// save the order with products
	if err := s.DBClient.SaveOrder(ctx, order); err != nil { // Save the order to the database
		s.Logger.Debug(err)
		writeJSONResponse(w, http.StatusInternalServerError, "Failed to save order")
		return
	}

	// return Response
	writeJSONResponse(w, http.StatusOK, order)
}
