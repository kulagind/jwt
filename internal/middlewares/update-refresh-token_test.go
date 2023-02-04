package middlewares

import (
	"context"
	"jwt/internal/constants"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/internal/services"
	"jwt/pkg/helpers/pg"
	"jwt/pkg/helpers/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type updateRefreshTokenTestCase struct {
	oldRefreshToken string
	isBlocked       bool
	isExpired       bool
	isExchanged     bool
	Claims          *models.RefreshTokenCustomClaims
	ExpectedCode    int
	ExpectedBody    []byte
}

func TestUpdatingRefreshToken(t *testing.T) {
	conn, err := pg.NewMockConnection(t)
	if err != nil {
		log.Fatal(err.Error())
	}
	d := repo.Init(conn)
	defer conn.Close()
	defer d()

	userTokenHash := utils.GenerateRandomString(15)
	user := &models.User{
		Id:        "test_testovich_111",
		Email:     "test_testovich_111",
		Password:  "test_testovich_111",
		TokenHash: userTokenHash,
	}

	validToken, err := services.GenerateRefreshToken(user)
	assert.Nil(t, err, "Generate refresh token")

	blockedTokenResponse := utils.UnsafetyToJson(models.ResponseError{
		Message:      "refresh token is blocked",
		Status:       http.StatusUnauthorized,
		InternalCode: 5,
		Stack:        "",
	})

	validResponse := []byte(user.Id)

	testCases := []updateRefreshTokenTestCase{
		{
			oldRefreshToken: validToken,
			Claims: &models.RefreshTokenCustomClaims{
				UserID:    user.Id,
				TokenHash: userTokenHash,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: services.GetExpiration(time.Hour, 15),
				},
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: validResponse,
		},
		{
			oldRefreshToken: validToken,
			isExpired:       true,
			Claims: &models.RefreshTokenCustomClaims{
				UserID:    user.Id,
				TokenHash: userTokenHash,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: services.GetExpiration(time.Hour, -15),
				},
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: validResponse,
		},
		{
			oldRefreshToken: validToken,
			isBlocked:       true,
			Claims: &models.RefreshTokenCustomClaims{
				UserID:    user.Id,
				TokenHash: userTokenHash,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: services.GetExpiration(time.Hour, 15),
				},
			},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: blockedTokenResponse,
		},
		{
			oldRefreshToken: validToken,
			isExchanged:     true,
			Claims: &models.RefreshTokenCustomClaims{
				UserID:    user.Id,
				TokenHash: userTokenHash,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: services.GetExpiration(time.Hour, -15),
				},
			},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: blockedTokenResponse,
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest("POST", "/test", nil)
		res := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), models.RequiredRenewalContextToken{}, test.oldRefreshToken)
		ctx = context.WithValue(ctx, models.ClaimsContextToken{}, test.Claims)
		ctx = context.WithValue(ctx, models.UserContextToken{}, user)
		req = req.WithContext(ctx)

		emptyRows := pgxpoolmock.NewRows([]string{}).ToPgxRows()

		if test.isExchanged {
			newToken := "new_token"
			updatedTokenRows := pgxpoolmock.NewRows([]string{"old_token", "new_token"}).AddRow(test.oldRefreshToken, newToken).ToPgxRows()
			// Select updated tokens
			conn.Mock.EXPECT().Query(gomock.Any(), gomock.Any(), test.oldRefreshToken).Return(updatedTokenRows, nil).Times(1)
			// Add old and new tokens to black list
			conn.Mock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)
		} else if test.isBlocked {
			blockedTokenRows := pgxpoolmock.NewRows([]string{"token"}).AddRow(test.oldRefreshToken).ToPgxRows()
			// Try to select updated tokens
			conn.Mock.EXPECT().Query(gomock.Any(), gomock.Any(), test.oldRefreshToken).Return(emptyRows, nil).Times(1)
			// Get token as blocked
			conn.Mock.EXPECT().Query(gomock.Any(), gomock.Any(), test.oldRefreshToken).Return(blockedTokenRows, nil).Times(1)
		} else {
			// Try to select updated tokens and select blocked token
			conn.Mock.EXPECT().Query(gomock.Any(), gomock.Any(), test.oldRefreshToken).Return(emptyRows, nil).Times(2)

			if test.isExpired {
				// Update tokenhash and refresh token
				conn.Mock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)
			}
		}

		mw := UpdateRefreshTokenIfRequired(handleUpdatingRefreshToken(test, t))
		mw.ServeHTTP(res, req)

		assert.Equal(t, test.ExpectedCode, res.Code, "Testing updating refresh token moddleware")
		assert.Equal(t, string(test.ExpectedBody), res.Body.String(), "Testing updating refresh token moddleware")
	}
}

func handleUpdatingRefreshToken(test updateRefreshTokenTestCase, t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(models.UserContextToken{}).(*models.User)

		assert.NotNil(t, user, "Testing updating refresh token moddleware")

		refreshCookie, err := r.Cookie(constants.TokenCookieName)
		if test.isExpired {
			assert.Nil(t, err, "Extract new refresh cookie")
			assert.NotEqual(t, test.oldRefreshToken, refreshCookie.Value, "Extract new refresh cookie")
		} else {
			assert.NotNil(t, err, "Extract new refresh cookie")
			assert.Nil(t, refreshCookie, "Extract new refresh cookie")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(user.Id))
	})
}
