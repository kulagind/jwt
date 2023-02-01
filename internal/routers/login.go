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
)

func login(w http.ResponseWriter, r *http.Request) {
	candidate := r.Context().Value(models.UserContextToken{}).(*models.User)

	user, err := repo.GetUserRepo().PrivateFindBy(context.Background(), "email", candidate.Email)
	if err != nil {
		if pg.CheckSqlError(err, "no rows in result set") {
			utils.WriteError(w, "User with this email and password doesn't exist", http.StatusUnauthorized, 2)
			return
		}
		utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
		return
	}

	if isAuth := services.Authenticate(candidate, user); !isAuth {
		utils.WriteError(w, "User with this email and password doesn't exist", http.StatusUnauthorized, 2)
		return
	}

	err = repo.GetUserRepo().UpdateTokenhash(context.Background(), user)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
		return
	}

	accessToken, err := services.GenerateAccessToken(user)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
		return
	}

	refreshToken, err := services.GenerateRefreshToken(user)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
		return
	}

	tokens := models.TokensResponse{
		AccessToken: models.AccessToken{Access_token: accessToken},
	}
	response, err := json.Marshal(tokens)
	if err != nil {
		utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
		return
	}

	c := services.GetRefreshCookie(refreshToken)
	http.SetCookie(w, &c)

	w.Write(response)
}
