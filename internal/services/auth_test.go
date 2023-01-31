package services

import (
	"jwt/internal/models"
	"jwt/pkg/helpers/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	type testCase struct {
		name           string
		input          string
		expectedResult bool
		hasError       bool
	}

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
	assert.Nil(t, err, "Generating access token. Check keys existing")

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

func TestRefreshTokenValidity(t *testing.T) {
	token, err := GenerateRefreshToken(user)
	assert.Nil(t, err, "Generating refresh token. Check keys existing")

	jwtToken, err := ParseRefreshToken(token)
	assert.Nil(t, err, "Parsing refresh token. Check that PEM keys are valid")

	claims, ok := jwtToken.Claims.(*models.RefreshTokenCustomClaims)
	assert.Equal(t, true, ok, "Parsing refresh token claims. Check that PEM keys are valid")

	assert.Equal(t, user.Id, claims.UserID, "Check refresh token claims: UserID")
	assert.Equal(t, "refresh", claims.KeyType, "Check refresh token claims: KeyType")
	assert.Equal(t, user.TokenHash, claims.TokenHash, "Check refresh token claims: TokenHash")
	difference := time.Now().Add(time.Minute*15).Unix() - claims.ExpiresAt
	lessThan10sec := difference < 10
	assert.Equal(t, true, lessThan10sec, claims.ExpiresAt, "Check refresh token claims: ExpiresAt")
}

func TestAuthHeader(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult string
		hasError       bool
	}

	testTable := []testCase{
		{
			name:           "Correct header with upperletter",
			input:          "Bearer abcde",
			expectedResult: "abcde",
			hasError:       false,
		},
		{
			name:           "Correct header with lowerletter",
			input:          "bearer abcde",
			expectedResult: "abcde",
			hasError:       false,
		},
		{
			name:           "Incorrect header (basic)",
			input:          "basic abcde",
			expectedResult: "",
			hasError:       true,
		},
		{
			name:           "Incorrect header (without the keyword bearer)",
			input:          "abcde",
			expectedResult: "",
			hasError:       true,
		},
	}

	for _, test := range testTable {
		header, err := ParseAuthHeader(test.input)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}

		assert.Equal(t, test.expectedResult, header, test.name)
	}
}
