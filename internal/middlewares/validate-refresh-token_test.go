package middlewares

import (
	"jwt/internal/models"
	"jwt/internal/services"
	"jwt/pkg/helpers/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type refreshTokenTestCase struct {
	ExpectedCode int
	Cookie       string
	ExpectedBody []byte
}

func init() {
	services.LoadEnv()
	os.Setenv("APP_MODE", "testing")
}

func TestRefreshTokenValidation(t *testing.T) {
	user := &models.User{
		Id:       "test_testovich_1",
		Email:    "test_testovich_2",
		Password: "test_testovich_3",
	}
	validToken, err := services.GenerateRefreshToken(user)
	assert.Nil(t, err, "Generate refresh token")

	invalidTokenResponse := utils.UnsafetyToJson(models.ResponseError{
		Message:      "Unauthorized",
		Status:       http.StatusUnauthorized,
		InternalCode: 3,
		Stack:        "",
	})

	validTokenResponse := []byte(user.Id)

	testCases := []refreshTokenTestCase{
		{
			Cookie:       validToken,
			ExpectedCode: http.StatusOK,
			ExpectedBody: validTokenResponse,
		},
		{
			Cookie:       "invalid_token-azazazaza",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: invalidTokenResponse,
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest("POST", "/test", nil)
		res := httptest.NewRecorder()

		c := services.GetRefreshCookie(test.Cookie)
		req.AddCookie(&c)

		mw := ValidateRefreshToken(handleRefreshTokenValidation(t))
		mw.ServeHTTP(res, req)

		assert.Equal(t, test.ExpectedCode, res.Code, "Testing refresh token moddleware")
		assert.Equal(t, string(test.ExpectedBody), res.Body.String(), "Testing refresh token moddleware")
	}
}

func handleRefreshTokenValidation(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(models.UserIdContextToken{}).(string)
		claims := r.Context().Value(models.ClaimsContextToken{}).(*models.RefreshTokenCustomClaims)
		oldRefreshToken := r.Context().Value(models.RequiredRenewalContextToken{}).(string)

		assert.NotNil(t, claims, "Testing refresh token moddleware")
		assert.NotEqual(t, "", oldRefreshToken, "Testing refresh token moddleware")
		assert.NotEqual(t, "", userId, "Testing refresh token moddleware")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userId))
	})
}
