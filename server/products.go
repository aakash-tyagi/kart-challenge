package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) ListProducts(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	// Get the limit and page variable from query parameters
	limit := r.URL.Query().Get("limit") // Example: "10"
	page := r.URL.Query().Get("page")   // Example: "1"
	// parse limit

	limitInt, valid := parseQueryParamToInt(limit, w)
	if !valid {
		return
	}
	pageInt, valid := parseQueryParamToInt(page, w)
	if !valid {
		return
	}

	// Fetch products from the database
	products, err := s.DBClient.ListProducts(ctx, limitInt, pageInt) // Call the database function
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, "Failed to retrieve products") // Handle error
		return
	}

	// Write the products to the response
	writeJSONResponse(w, http.StatusOK, products)
}

func (s *Server) GetProductById(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	vars := mux.Vars(r)            // Extract variables from the URL
	productId := vars["productId"] // Get the product ID from the URL

	// get the product from db by product id
	product, err := s.DBClient.GetProductById(ctx, productId) // Call the database function
	if err != nil {
		writeJSONResponse(w, http.StatusNotFound, "Product not found") // Handle error
		return
	}

	// Write the product to the response
	writeJSONResponse(w, http.StatusOK, product)
}
