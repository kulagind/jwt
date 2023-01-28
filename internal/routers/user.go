package routers

import (
	"context"
	"encoding/json"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/pkg/helpers/pg"
	"jwt/pkg/helpers/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func getCurrentUser(w http.ResponseWriter, r *http.Request) {
	privateUser := r.Context().Value(models.UserContextToken{}).(*models.User)
	user := &models.UserResponse{
		Id:    privateUser.Id,
		Email: privateUser.Email,
		Name:  privateUser.Name,
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
		return
	}
	w.Write(userJson)
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user, err := repo.GetUserRepo().FindBy(context.Background(), "id", vars["id"])
	if err != nil {
		if pg.CheckSqlError(err, "no rows in result set") {
			utils.WriteError(w, "User with ID doesn't exist", http.StatusBadRequest, -1)
			return
		}
		utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
		return
	}

	w.Write(userJson)
}
