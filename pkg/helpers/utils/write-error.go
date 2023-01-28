package utils

import (
	"encoding/json"
	"jwt/internal/models"
	"net/http"
	"os"
	"runtime/debug"
)

func WriteError(w http.ResponseWriter, message string, httpStatus int, internalCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpStatus)

	var stack string
	if os.Getenv("APP_MODE") == "dev" {
		stack = string(debug.Stack())
	}

	customError := &models.ResponseError{
		Message:      message,
		Status:       httpStatus,
		InternalCode: internalCode,
		Stack:        stack,
	}
	errorJson, err := json.Marshal(customError)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(errorJson)
}
