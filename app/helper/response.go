package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// NewErrorResponse for error response
func NewErrorResponse(w http.ResponseWriter, statusCode int, errors string) {
	response, _ := json.Marshal(map[string]interface{}{
		"message": errors,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
	fmt.Printf("error: %s", errors)
}

// NewSuccessResponse for success response
func NewSuccessResponse(w http.ResponseWriter, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}
