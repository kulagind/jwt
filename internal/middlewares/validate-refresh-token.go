package middlewares

import (
	"context"
	"errors"
	"io/ioutil"
	"jwt/internal/constants"
	"jwt/internal/models"
	"jwt/internal/repo"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshCookie, err := r.Cookie(constants.TokenCookieName)
		if err != nil && refreshCookie.Value != "" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(
			refreshCookie.Value,
			&models.RefreshTokenCustomClaims{},
			func(t *jwt.Token) (interface{}, error) {
				pubBytes, err := ioutil.ReadFile(os.Getenv("REFRESH_TOKEN_PUBLIC_KEY_PATH"))
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
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*models.RefreshTokenCustomClaims)
		if !ok {
			http.Error(w, "could not parse refresh token claims", http.StatusInternalServerError)
			return
		}

		if claims.ExpiresAt < time.Now().Local().Unix() {
			http.Error(w, "refresh token is expired. Please login using email and password again", http.StatusForbidden)
			return
		}

		var candidate *models.User
		candidate, err = repo.GetUserRepo().PrivateFindBy(context.Background(), "id", claims.UserID)
		if err != nil || claims.TokenHash != candidate.TokenHash {
			http.Error(w, "refresh token is incorrect", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserToken{}, candidate)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
