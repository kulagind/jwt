package middlewares

import (
	"context"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/internal/services"
	"jwt/pkg/helpers/utils"
	"net/http"
)

func UpdateRefreshTokenIfRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requiredRenewal := r.Context().Value(models.RequiredRenewalContextToken{}).(string)
		user := r.Context().Value(models.UserContextToken{}).(*models.User)
		if requiredRenewal != "" {
			err := repo.GetUserRepo().UpdateTokenhash(context.Background(), user)
			if err != nil {
				utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
				return
			}

			newRefreshToken, err := services.GenerateRefreshToken(user)
			if err != nil {
				utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
				return
			}
			err = repo.GetTokenRepo().UpdateRefresh(context.Background(), requiredRenewal, newRefreshToken)
			if err != nil {
				utils.WriteError(w, err.Error(), http.StatusInternalServerError, 0)
				return
			}

			c := services.GetRefreshCookie(newRefreshToken)
			http.SetCookie(w, &c)
		}

		next.ServeHTTP(w, r)
	})
}
