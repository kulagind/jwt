package routers

import (
	"encoding/json"
	"jwt/internal/models"
	"jwt/internal/services"
	"jwt/pkg/helpers/utils"
	"net/http"
)

func updateAccessToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(models.UserContextToken{}).(*models.User)

	accessToken, err := services.GenerateAccessToken(user)
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

	w.Write(response)
}
