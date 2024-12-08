package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aakash-tyagi/kart-challenge/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) AddProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// Parse the request body to get the product details
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		s.Logger.Error(err)
		writeJSONResponse(w, http.StatusBadRequest, "Invalid request payload") // Handle JSON decode error
		return
	}

	if err := product.Validate(); err != nil {
		s.Logger.Error(err)
		writeJSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call the database function to add the product
	if err := s.DBClient.AddProduct(ctx, product); err != nil {
		s.Logger.Error(err)
		writeJSONResponse(w, http.StatusUnprocessableEntity, "Failed to add product: "+err.Error()) // Handle error
		return
	}

	writeJSONResponse(w, http.StatusCreated, "Product added successfully")
}

func (s *Server) ListProducts(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	// Get the limit and page variable from query parameters
	limit := r.URL.Query().Get("limit") // Example: "10"
	page := r.URL.Query().Get("page")   // Example: "1"
	// parse limit

	if limit == "" {
		limit = "10" // Default limit
	}
	if page == "" {
		page = "1" // Default page
	}

	limitInt, valid := parseQueryParamToInt(limit, w)
	if !valid {
		return
	}
	pageInt, valid := parseQueryParamToInt(page, w)
	if !valid {
		return
	}

	// Fetch products from the database
	products, totalProducts, err := s.DBClient.ListProducts(ctx, pageInt, limitInt) // Call the database function
	if err != nil {
		s.Logger.Error(err)
		writeJSONResponse(w, http.StatusInternalServerError, "Failed to retrieve products") // Handle error
		return
	}

	listProductRes := &models.ListProduct{
		Products:      products,
		TotalProducts: int(totalProducts),
	}

	// Write the products to the response
	writeJSONResponse(w, http.StatusOK, listProductRes)
}

func (s *Server) GetProductById(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	vars := mux.Vars(r)            // Extract variables from the URL
	productId := vars["productId"] // Get the product ID from the URL

	// Validate productId as primitive.ObjectID
	objectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		s.Logger.Error(err)
		writeJSONResponse(w, http.StatusBadRequest, "Invalid product ID") // Handle invalid ID error
		return
	}

	// get the product from db by product id
	product, err := s.DBClient.GetProductById(ctx, objectId) // Call the database function
	if err != nil {
		s.Logger.Error(err)
		writeJSONResponse(w, http.StatusNotFound, "Product not found") // Handle error
		return
	}

	// Write the product to the response
	writeJSONResponse(w, http.StatusOK, product)
}
