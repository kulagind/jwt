package middlewares

import (
	"context"
	"errors"
	"jwt/internal/constants"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/pkg/helpers/pg"
	"jwt/pkg/helpers/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type getUserTestCase struct {
	ExpectedCode int
	UserId       string
	Claims       *models.RefreshTokenCustomClaims
	ExpectedBody []byte
}

func TestGettingUserMiddleware(t *testing.T) {
	conn, err := pg.NewMockConnection(t)
	if err != nil {
		log.Fatal(err.Error())
	}
	d := repo.Init(conn)
	defer conn.Close()
	defer d()

	userId := "test_testovich_id"
	userTokenHash := utils.GenerateRandomString(15)
	user := &models.User{
		Id:        userId,
		Email:     "test_testovich",
		Password:  "test_testovich",
		TokenHash: userTokenHash,
	}

	invalidResponse := utils.UnsafetyToJson(models.ResponseError{
		Message:      "Unauthorized",
		Status:       http.StatusUnauthorized,
		InternalCode: 3,
		Stack:        "",
	})

	validResponse := []byte(user.Id)

	testCases := []getUserTestCase{
		{
			UserId:       user.Id,
			ExpectedCode: http.StatusOK,
			ExpectedBody: validResponse,
		},
		{
			UserId:       "invalididazazazaza",
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: invalidResponse,
		},
		{
			UserId: user.Id,
			Claims: &models.RefreshTokenCustomClaims{
				UserID:    user.Id,
				KeyType:   "refresh",
				TokenHash: userTokenHash,
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: validResponse,
		},
		{
			UserId: user.Id,
			Claims: &models.RefreshTokenCustomClaims{
				UserID:    user.Id,
				KeyType:   "refresh",
				TokenHash: utils.GenerateRandomString(15),
			},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: invalidResponse,
		},
		{
			UserId: "anotheruser",
			Claims: &models.RefreshTokenCustomClaims{
				UserID:    user.Id,
				KeyType:   "refresh",
				TokenHash: userTokenHash,
			},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: invalidResponse,
		},
	}

	for _, test := range testCases {
		req := httptest.NewRequest("POST", "/test", nil)
		res := httptest.NewRecorder()

		ctx := context.WithValue(req.Context(), models.UserIdContextToken{}, test.UserId)
		if test.Claims != nil {
			ctx = context.WithValue(ctx, models.ClaimsContextToken{}, test.Claims)
		}
		req = req.WithContext(ctx)

		rows := pgxpoolmock.NewRows([]string{}).ToPgxRows()
		err := errors.New(constants.NO_ROWS)
		if test.UserId == user.Id {
			rows = pgxpoolmock.NewRows([]string{"id", "email", "name", "password", "tokenhash", "created_at", "updated_at"}).AddRow(userId, user.Email, "name", user.Password, user.TokenHash, time.Now(), time.Now()).ToPgxRows()
			err = nil
		}
		conn.Mock.EXPECT().Query(gomock.Any(), gomock.Any(), test.UserId).Return(rows, err).Times(1)

		mw := GetUserById(handleGettingUser(t))
		mw.ServeHTTP(res, req)

		assert.Equal(t, test.ExpectedCode, res.Code, "Testing gettting user by ID middleware")
		assert.Equal(t, string(test.ExpectedBody), res.Body.String(), "Testing gettting user by ID middleware")
	}
}

func handleGettingUser(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(models.UserContextToken{}).(*models.User)
		assert.NotNil(t, user, "Getting user from middleware's context")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(user.Id))
	})
}
