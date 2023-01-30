package services

import (
	"jwt/internal/models"
	"jwt/pkg/helpers/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name           string
	input          string
	expectedResult bool
	hasError       bool
}

var user *models.User

func init() {
	user = &models.User{
		Id:        "test-id",
		Name:      "Test Testovich",
		TokenHash: utils.GenerateRandomString(15),
		Password:  "",
		Email:     "test-id@test.ru",
	}

	LoadEnv()
}

func TestAuthenticate(t *testing.T) {
	valid_password := "test_password"
	testTable := []testCase{
		{
			name:           "compare correct passwords",
			input:          valid_password,
			expectedResult: true,
			hasError:       false,
		},
		{
			name:           "compare different password",
			input:          "invalid_password",
			expectedResult: false,
			hasError:       false,
		},
	}

	hashedPass, err := HashPassword(valid_password)
	assert.Nil(t, err, "encrypts password")

	user.Password = hashedPass

	for _, test := range testTable {
		candidate := &models.User{
			Password: test.input,
			Email:    "test-id@test.ru",
		}

		isAuth := Authenticate(candidate, user)
		assert.Equal(t, test.expectedResult, isAuth, test.name)
	}
}

func TestAccessTokenValidity(t *testing.T) {
	token, err := GenerateAccessToken(user)
	assert.Nil(t, err, "Generating token. Check keys existing")

	jwtToken, err := ParseAccessToken(token)
	assert.Nil(t, err, "Parsing access token. Check that PEM keys are valid")

	claims, ok := jwtToken.Claims.(*models.AccessTokenCustomClaims)
	assert.Equal(t, true, ok, "Parsing access token claims. Check that PEM keys are valid")

	assert.Equal(t, user.Id, claims.UserID, "Check access token claims: UserID")
	assert.Equal(t, "access", claims.KeyType, "Check access token claims: KeyType")
	difference := time.Now().Add(time.Minute*15).Unix() - claims.ExpiresAt
	lessThan10sec := difference < 10
	assert.Equal(t, true, lessThan10sec, claims.ExpiresAt, "Check access token claims: ExpiresAt")
}
