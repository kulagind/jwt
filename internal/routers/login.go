package routers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/internal/services"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	candidate := r.Context().Value(models.UserToken{}).(models.User)

	user, err := repo.GetUserRepo().PrivateFindBy(context.Background(), "email", candidate.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isAuth := services.Authenticate(&candidate, user); !isAuth {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken, err := 

	userJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(userJson)
}

func _create(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var candidate models.User
	if err := json.Unmarshal(data, &candidate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !candidate.Valid() {
		http.Error(w, "Invalid body fields", http.StatusBadRequest)
		return
	}

	user, err := repo.GetUserRepo().Create(context.Background(), &candidate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(userJson)
}
