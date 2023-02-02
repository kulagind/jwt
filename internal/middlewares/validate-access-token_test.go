package middlewares

import (
	"fmt"
	"jwt/internal/models"
	"jwt/internal/services"
	"jwt/pkg/helpers/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type accessTokenTestCase struct {
	ExpectedCode int
	Header       string
	ExpectedBody []byte
	HeaderKey    string
}

func init() {
	services.LoadEnv()
	os.Setenv("APP_MODE", "testing")
}

func TestAccessTokenValidation(t *testing.T) {
	user := &models.User{
		Id:       "test_testovich",
		Email:    "test_testovich",
		Password: "test_testovich",
	}
	validToken, err := services.GenerateAccessToken(user)
	assert.Nil(t, err, "Generate access token")

	expiredToken, err := services.GenerateAccessToken(user, -10)
	assert.Nil(t, err, "Generate expired access token")

	invalidTokenResponse := utils.UnsafetyToJson(models.ResponseError{
		Message:      "Unauthorized",
		Status:       http.StatusUnauthorized,
		InternalCode: 3,
		Stack:        "",
	})

	expiredTokenResponse := utils.UnsafetyToJson(models.ResponseError{
		Message:      "access token is expired",
		Status:       http.StatusForbidden,
		InternalCode: 4,
		Stack:        "",
	})

	validTokenResponse := []byte(user.Id)

	testCases := []accessTokenTestCase{
		{
			Header:       fmt.Sprintf(`Bearer %s`, validToken),
			HeaderKey:    "Authorization",
			ExpectedCode: http.StatusOK,
			ExpectedBody: validTokenResponse,
		},
		{
			Header:       fmt.Sprintf(`bearer %s`, validToken),
			HeaderKey:    "authorization",
			ExpectedCode: http.StatusOK,
			ExpectedBody: validTokenResponse,
		},
		{
			Header:       fmt.Sprintf(`Bearer %s`, expiredToken),
			HeaderKey:    "Authorization",
			ExpectedCode: http.StatusForbidden,
			ExpectedBody: expiredTokenResponse,
		},
		{
			Header:       fmt.Sprintf(`Bearer %s`, "invalidTokenazazazazaz"),
			HeaderKey:    "Authorization",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: invalidTokenResponse,
		},
		{
			Header:       fmt.Sprintf(`bearer %s`, expiredToken),
			HeaderKey:    "authorization",
			ExpectedCode: http.StatusForbidden,
			ExpectedBody: expiredTokenResponse,
		},
		{
			Header:       fmt.Sprintf(`bearer %s`, ""),
			HeaderKey:    "authorization",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: invalidTokenResponse,
		},
		{
			Header:       fmt.Sprintf(`Basic %s`, validToken),
			HeaderKey:    "Authorization",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: invalidTokenResponse,
		},
		{
			Header:       fmt.Sprintf(`Bearer %s`, validToken),
			HeaderKey:    "X-Token",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: invalidTokenResponse,
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest("POST", "/test", nil)
		req.Header.Add(test.HeaderKey, test.Header)
		res := httptest.NewRecorder()
		mw := ValidateAccessToken(handleAccessTokenValidation(t))
		mw.ServeHTTP(res, req)

		assert.Equal(t, test.ExpectedCode, res.Code, "Testing access token moddleware")
		assert.Equal(t, string(test.ExpectedBody), res.Body.String(), "Testing access token moddleware")
	}
}

func handleAccessTokenValidation(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(models.UserIdContextToken{}).(string)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userId))
	})
}
