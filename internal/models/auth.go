package models

type AccessTokenCustomClaims struct {
	UserID  string
	KeyType string
	jwt.StandardClaims
}
