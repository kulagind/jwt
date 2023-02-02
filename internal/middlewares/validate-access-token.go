package middlewares

import (
	"context"
	"jwt/internal/models"
	"jwt/internal/services"
	"jwt/pkg/helpers/utils"
	"net/http"
	"time"
)

func ValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		headerToken, err := services.ParseAuthHeader(header)
		if headerToken == "" || err != nil {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		token, err := services.ParseAccessToken(headerToken)
		if err != nil {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		claims, ok := token.Claims.(*models.AccessTokenCustomClaims)
		if !ok {
			utils.WriteError(w, "could not parse access token claims", http.StatusInternalServerError, 0)
			return
		}

		if claims.ExpiresAt < time.Now().Local().Unix() {
			utils.WriteError(w, "access token is expired", http.StatusForbidden, 4)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserIdContextToken{}, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
