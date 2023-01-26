package utils

import (
	"encoding/json"
	"fmt"
	"jwt/internal/models"
	"net/http"
	"os"
	"runtime/debug"
)

func WriteError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	var stack string
	if os.Getenv("APP_MODE") == "dev" {
		stack = string(debug.Stack())
		fmt.Println(stack)
	}

	customError := &models.ResponseError{
		Message: message,
		Status:  status,
		Stack:   stack,
	}
	errorJson, err := json.Marshal(customError)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(errorJson)
}
