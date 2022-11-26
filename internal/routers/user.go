package routers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"jwt/internal/models"
	"jwt/internal/repo"
	"net/http"

	"github.com/gorilla/mux"
)

func handleUsers(router *mux.Router) {
	router.HandleFunc("/{id}", _getById).Methods(http.MethodGet)
	router.HandleFunc("", _create).Methods(http.MethodPost)
}

func _getById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user, err := repo.GetUserRepo().FindById(context.Background(), vars["id"])
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
