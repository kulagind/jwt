package middlewares

import (
	"context"
	"errors"
	"jwt/internal/constants"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/internal/services"
	"jwt/pkg/helpers/utils"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		headerToken, err := services.ParseAuthHeader(header)
		if headerToken == "" || err != nil {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		token, err := jwt.ParseWithClaims(
			headerToken,
			&models.AccessTokenCustomClaims{},
			func(t *jwt.Token) (interface{}, error) {
				pubBytes, err := os.ReadFile(
					path.Join(constants.ProjectPath(), os.Getenv("ACCESS_TOKEN_PUBLIC_KEY_PATH")),
				)
				if err != nil {
					return nil, errors.New("could not parse access token. please try again later")
				}

				pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
				if err != nil {
					return nil, errors.New("could not parse access token. please try again later")
				}
				return pubKey, nil
			},
		)
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

		var candidate *models.User
		candidate, err = repo.GetUserRepo().PrivateFindBy(context.Background(), "id", claims.UserID)
		if err != nil {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserContextToken{}, candidate)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
