package routers

import (
	"context"
	"encoding/json"
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

	err = repo.GetUserRepo().UpdateTokenhash(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, err := services.GenerateAccessToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := services.GenerateRefreshToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokens := models.TokensResponse{
		AccessToken:  models.AccessToken{Access_token: accessToken},
		RefreshToken: models.RefreshToken{Refresh_token: refreshToken},
	}
	response, err := json.Marshal(tokens)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
