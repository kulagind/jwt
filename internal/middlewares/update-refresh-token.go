package middlewares

import (
	"context"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/internal/services"
	"jwt/pkg/helpers/utils"
	"net/http"
	"time"
)

func UpdateRefreshTokenIfRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oldRefreshToken := r.Context().Value(models.RequiredRenewalContextToken{}).(string)
		claims := r.Context().Value(models.ClaimsContextToken{}).(*models.RefreshTokenCustomClaims)

		requiredRenewal := ""
		if claims.ExpiresAt < time.Now().Local().Unix() {
			err := repo.GetTokenRepo().CheckRefresh(context.Background(), oldRefreshToken)
			if err != nil {
				utils.WriteError(w, "refresh token is blocked", http.StatusUnauthorized, 5)
				return
			}
			requiredRenewal = oldRefreshToken
		}

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
