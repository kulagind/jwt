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

func GetExpiration(unit time.Duration, count int64) int64 {
	return time.Now().Add(time.Duration(count) * unit).Unix()
}

func GenerateAccessToken(user *models.User, expiration ...int64) (string, error) {
	expiresAt := GetExpiration(time.Minute, 15)
	if len(expiration) > 0 {
		expiresAt = GetExpiration(time.Minute, expiration[0])
	}

	claims := models.AccessTokenCustomClaims{
		UserID:  user.Id,
		KeyType: "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
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

func GenerateRefreshToken(user *models.User, expiration ...int64) (string, error) {
	expiresAt := GetExpiration(time.Hour, 24*7)
	if len(expiration) > 0 {
		expiresAt = GetExpiration(time.Hour, expiration[0])
	}

	claims := models.RefreshTokenCustomClaims{
		UserID:    user.Id,
		KeyType:   "refresh",
		TokenHash: user.TokenHash,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
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

func HashPassword(pass string) (string, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encryptedPass), nil
}

func ParseAccessToken(header string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		header,
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
	if err != nil && !strings.Contains(err.Error(), "token is expired by") {
		return nil, err
	}
	return token, nil
}

func ParseRefreshToken(refreshToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		refreshToken,
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
	if err != nil && !strings.Contains(err.Error(), "token is expired by") {
		return nil, err
	}
	return token, nil
}
