package middlewares

import (
	"context"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/pkg/helpers/utils"
	"net/http"
)

func GetUserById(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(models.UserIdContextToken{}).(string)
		claims := r.Context().Value(models.ClaimsContextToken{})

		var candidate *models.User
		candidate, err := repo.GetUserRepo().PrivateFindBy(context.Background(), "id", userId)
		if err != nil || (claims != nil && claims.(*models.RefreshTokenCustomClaims).TokenHash != "" && claims.(*models.RefreshTokenCustomClaims).TokenHash != candidate.TokenHash) {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserContextToken{}, candidate)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
