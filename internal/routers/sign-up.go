package routers

import (
	"context"
	"encoding/json"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/internal/services"
	"jwt/pkg/helpers/pg"
	"jwt/pkg/helpers/utils"
	"net/http"

	"github.com/jackc/pgerrcode"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	candidate := r.Context().Value(models.UserContextToken{}).(*models.User)
	hashedPass, err := services.HashPassword(candidate.Password)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	candidate.Password = hashedPass

	user, err := repo.GetUserRepo().Create(context.Background(), candidate)
	if err != nil {
		if pg.CheckSqlError(err, pgerrcode.UniqueViolation) {
			utils.WriteError(w, "User with this email already exists", http.StatusBadRequest, 1)
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

	w.WriteHeader(http.StatusCreated)
	w.Write(userJson)
}
