package routers

import (
	"context"
	"encoding/json"
	"jwt/internal/models"
	"jwt/internal/repo"
	"net/http"

	"github.com/gorilla/mux"
)

func getCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(models.UserContextToken{}).(models.UserResponse)

	userJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(userJson)
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user, err := repo.GetUserRepo().FindBy(context.Background(), "id", vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(userJson)
}
