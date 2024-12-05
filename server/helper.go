package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func parseQueryParamToInt(param string, w http.ResponseWriter) (int, bool) {
	value, err := strconv.Atoi(param)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, "Invalid parameter") // Handle error
		return 0, false
	}
	return value, true
}

// Helper function to write JSON response
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError) // Handle encoding error
	}
}
