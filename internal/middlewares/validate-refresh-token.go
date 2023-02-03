package middlewares

import (
	"context"
	"jwt/internal/constants"
	"jwt/internal/models"
	"jwt/internal/services"
	"jwt/pkg/helpers/utils"
	"net/http"
)

func ValidateRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshCookie, err := r.Cookie(constants.TokenCookieName)
		if err != nil || (refreshCookie == nil && refreshCookie.Value != "") {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		token, err := services.ParseRefreshToken(refreshCookie.Value)
		if err != nil {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		claims, ok := token.Claims.(*models.RefreshTokenCustomClaims)
		if !ok {
			utils.WriteError(w, "could not parse refresh token claims", http.StatusInternalServerError, 0)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserIdContextToken{}, claims.UserID)
		ctx = context.WithValue(ctx, models.ClaimsContextToken{}, claims)
		ctx = context.WithValue(ctx, models.RequiredRenewalContextToken{}, refreshCookie.Value)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
