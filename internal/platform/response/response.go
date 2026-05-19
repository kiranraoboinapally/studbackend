package response

import (
	"encoding/json"
	"net/http"

	"university-erp-backend/internal/platform/apperrors"
)

// JSON writes a success JSON response.
func JSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// List writes a paginated list response.
func List(w http.ResponseWriter, data interface{}, total int64, page, pageSize int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
		"meta": map[string]interface{}{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Error writes a JSON error response. Accepts *apperrors.AppError or plain error.
func Error(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.Code)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   appErr.Message,
			"detail":  appErr.Detail,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   "internal server error",
	})
}

// Created writes a 201 response.
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

// NoContent writes a 204 response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
