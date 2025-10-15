package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Success bool  `json:"success"`
	Error   Error `json:"error"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func JsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	successResponse := &SuccessResponse{
		Success: true,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(successResponse); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func ErrorJsonResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := &ErrorResponse{
		Success: false,
		Error: Error{
			Message: message,
		},
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}
