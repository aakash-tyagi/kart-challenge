package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) ListProducts(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	// Get the limit and page variable from query parameters
	limit := r.URL.Query().Get("limit") // Example: "10"
	page := r.URL.Query().Get("page")   // Example: "1"

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}

	// Fetch products from the database
	products, err := s.DBClient.ListProducts(ctx, limitInt, pageInt) // Call the database function
	if err != nil {
		http.Error(w, "Failed to retrieve products", http.StatusInternalServerError) // Handle error
		return
	}

	// Write the products to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode products", http.StatusInternalServerError) // Handle encoding error
	}
}

func (s *Server) GetProductById(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	vars := mux.Vars(r)            // Extract variables from the URL
	productId := vars["productId"] // Get the product ID from the URL

	// get the product from db by product id
	product, err := s.DBClient.GetProductById(ctx, productId) // Call the database function
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound) // Handle error
		return
	}

	// Write the product to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode product", http.StatusInternalServerError) // Handle encoding error
	}
}
