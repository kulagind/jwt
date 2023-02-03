package models

type AccessTokenContextToken struct{}

type RefreshTokenContextToken struct{}

type RequiredRenewalContextToken struct{}

type UpdatedRefreshTokenContextToken struct{}

type ClaimsContextToken struct{}

type RepoUpdatedTokens struct {
	OldToken string `validate:"required" sql:"old_token"`
	NewToken string `validate:"required" sql:"new_token"`
}

type RepoBlocked struct {
	Token string `validate:"required" sql:"token"`
}
