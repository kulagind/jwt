package services

import (
	"jwt/internal/models"
	"jwt/pkg/helpers/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name           string
	input          string
	expectedResult bool
	hasError       bool
}

var user *models.User

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

	user = &models.User{
		Id:        "test-id",
		Name:      "Test Testovich",
		TokenHash: utils.GenerateRandomString(15),
		Password:  hashedPass,
		Email:     "test-id@test.ru",
	}

	for _, test := range testTable {
		candidate := &models.User{
			Password: test.input,
			Email:    "test-id@test.ru",
		}

		isAuth := Authenticate(candidate, user)
		assert.Equal(t, test.expectedResult, isAuth, test.name)
	}
}

func TestGenerateAccessToken(t *testing.T) {
	token, err := GenerateAccessToken(user)
	assert.Nil(t, err, "Couldn't generate token. Check keys existing")
}
