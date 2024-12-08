package server

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func searchInFile(filePath, target string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), target) {
			return true, nil
		}
	}

	return false, scanner.Err()
}
