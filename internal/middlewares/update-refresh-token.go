package middlewares

import (
	"context"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/internal/services"
	"net/http"
)

func UpdateRefreshTokenIfRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requiredRenewal := r.Context().Value(models.RequiredRenewalContextToken{}).(string)
		user := r.Context().Value(models.UserContextToken{}).(*models.User)
		if requiredRenewal != "" {
			err := repo.GetUserRepo().UpdateTokenhash(context.Background(), user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			newRefreshToken, err := services.GenerateRefreshToken(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			repo.GetTokenRepo().UpdateRefresh(context.Background(), requiredRenewal, newRefreshToken)
		}

		next.ServeHTTP(w, r)
	})
}
