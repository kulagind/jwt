package services

import (
	"errors"
	"io/ioutil"
	"jwt/internal/constants"
	"jwt/internal/models"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(candidate *models.User, user *models.User) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(candidate.Password)); err != nil {
		return false
	}
	return true
}

func GenerateAccessToken(user *models.User) (string, error) {
	claims := models.AccessTokenCustomClaims{
		UserID:  user.Id,
		KeyType: "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	signBytes, err := ioutil.ReadFile(
		path.Join(constants.ProjectPath(), os.Getenv("ACCESS_TOKEN_PRIVATE_KEY_PATH")),
	)
	if err != nil {
		return "", errors.New("could not generate access token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", errors.New("could not generate access token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

func GenerateRefreshToken(user *models.User) (string, error) {
	claims := models.RefreshTokenCustomClaims{
		UserID:    user.Id,
		KeyType:   "refresh",
		TokenHash: user.TokenHash,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	signBytes, err := ioutil.ReadFile(
		path.Join(constants.ProjectPath(), os.Getenv("REFRESH_TOKEN_PRIVATE_KEY_PATH")),
	)
	if err != nil {
		return "", errors.New("could not generate refresh token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", errors.New("could not generate refresh token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

func ParseAuthHeader(header string) (string, error) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("Unauthorized")
	}
	return parts[1], nil
}

func GetRefreshCookie(token string) http.Cookie {
	return http.Cookie{
		HttpOnly: true,
		Name:     constants.TokenCookieName,
		Value:    token,
	}
}
