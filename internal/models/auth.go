package models

import "github.com/golang-jwt/jwt/v4"

type AccessTokenCustomClaims struct {
	UserID  string
	KeyType string
	jwt.StandardClaims
}

type RefreshTokenCustomClaims struct {
	UserID    string
	KeyType   string
	TokenHash string
	jwt.StandardClaims
}

type AccessToken struct {
	Access_token string `json:"accessToken"`
}

type RefreshToken struct {
	Refresh_token string `json:"refreshToken"`
}

type TokensResponse struct {
	RefreshToken
	AccessToken
}
