package middlewares

import (
	"context"
	"errors"
	"io/ioutil"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/internal/services"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		headerToken, err := services.ParseAuthHeader(header)
		if headerToken == "" || err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(
			headerToken,
			&models.RefreshTokenCustomClaims{},
			func(t *jwt.Token) (interface{}, error) {
				pubBytes, err := ioutil.ReadFile(os.Getenv("ACCESS_TOKEN_PUBLIC_KEY_PATH"))
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
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*models.AccessTokenCustomClaims)
		if !ok {
			http.Error(w, "could not parse access token claims", http.StatusInternalServerError)
			return
		}

		if claims.ExpiresAt < time.Now().Local().Unix() {
			// TODO: renew access token
			http.Error(w, "access token is expired", http.StatusForbidden)
			return
		}

		var candidate *models.User
		candidate, err = repo.GetUserRepo().PrivateFindBy(context.Background(), "id", claims.UserID)
		if err != nil {
			http.Error(w, "access token is incorrect", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserToken{}, candidate)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
