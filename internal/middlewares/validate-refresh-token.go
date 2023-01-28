package middlewares

import (
	"context"
	"errors"
	"io/ioutil"
	"jwt/internal/constants"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/pkg/helpers/utils"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshCookie, err := r.Cookie(constants.TokenCookieName)
		if err != nil && refreshCookie.Value != "" {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		token, err := jwt.ParseWithClaims(
			refreshCookie.Value,
			&models.RefreshTokenCustomClaims{},
			func(t *jwt.Token) (interface{}, error) {
				pubBytes, err := ioutil.ReadFile(
					path.Join(constants.ProjectPath(), os.Getenv("REFRESH_TOKEN_PUBLIC_KEY_PATH")),
				)
				if err != nil {
					return nil, errors.New("could not parse refresh token. please try again later")
				}

				pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
				if err != nil {
					return nil, errors.New("could not parse refresh token. please try again later")
				}
				return pubKey, nil
			},
		)
		if err != nil {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		claims, ok := token.Claims.(*models.RefreshTokenCustomClaims)
		if !ok {
			utils.WriteError(w, "could not parse refresh token claims", http.StatusInternalServerError, 0)
			return
		}

		var candidate *models.User
		candidate, err = repo.GetUserRepo().PrivateFindBy(context.Background(), "id", claims.UserID)
		if err != nil || claims.TokenHash != candidate.TokenHash {
			utils.WriteError(w, "Unauthorized", http.StatusUnauthorized, 3)
			return
		}

		requiredRenewalToken := ""
		if claims.ExpiresAt < time.Now().Local().Unix() {
			err = repo.GetTokenRepo().CheckRefresh(context.Background(), refreshCookie.Value)
			if err != nil {
				utils.WriteError(w, "refresh token is blocked", http.StatusUnauthorized, 5)
				return
			}
			requiredRenewalToken = refreshCookie.Value
		}

		ctx := context.WithValue(r.Context(), models.UserContextToken{}, candidate)
		ctx = context.WithValue(ctx, models.RequiredRenewalContextToken{}, requiredRenewalToken)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
