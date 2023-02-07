package routers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"jwt/internal/constants"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/internal/services"
	"jwt/pkg/helpers/pg"
	"jwt/pkg/helpers/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type loginTestCase struct {
	expectedCode int
	body         []byte
	expectedBody []byte
	email        string
	password     string
}

func init() {
	services.LoadEnv()
	os.Setenv("APP_MODE", "testing")
}

func TestLogin(t *testing.T) {
	mux := mux.NewRouter()
	mux.HandleFunc("/login", login).Methods("POST")

	conn, err := pg.NewMockConnection(t)
	if err != nil {
		log.Fatal(err.Error())
	}
	d := repo.Init(conn)
	defer conn.Close()
	defer d()

	userNotFoundResponse := utils.UnsafetyToJson(models.ResponseError{
		Message:      "User with this email and password doesn't exist",
		Status:       http.StatusUnauthorized,
		InternalCode: 2,
		Stack:        "",
	})

	email := "test_testovich"
	password := "test_testovich"

	testCases := []loginTestCase{
		{
			body:         []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)),
			email:        email,
			password:     password,
			expectedCode: http.StatusOK,
			expectedBody: []byte(""),
		},
		{
			body:         []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, "")),
			email:        email,
			password:     "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: userNotFoundResponse,
		},
		{
			body:         []byte(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, "", password)),
			email:        "",
			password:     password,
			expectedCode: http.StatusUnauthorized,
			expectedBody: userNotFoundResponse,
		},
		{
			body:         []byte(fmt.Sprintf(`{"email": "%s"}`, email)),
			email:        email,
			expectedCode: http.StatusUnauthorized,
			expectedBody: userNotFoundResponse,
		},
		{
			body:         []byte(fmt.Sprintf(`{"password": "%s"}`, password)),
			email:        "",
			password:     password,
			expectedCode: http.StatusUnauthorized,
			expectedBody: userNotFoundResponse,
		},
	}

	for _, test := range testCases {
		rows := pgxpoolmock.NewRows([]string{}).ToPgxRows()
		sqlErr := errors.New(constants.NO_ROWS)

		hashedPass, err := services.HashPassword(test.password)
		assert.Nil(t, err, "Testing login")

		user := &models.User{
			Id:       test.email,
			Email:    test.email,
			Password: hashedPass,
		}

		candidate := &models.User{
			Email:    test.email,
			Password: test.password,
		}

		if test.expectedCode == http.StatusOK {
			sqlErr = nil
			rows = pgxpoolmock.NewRows([]string{"id", "email", "name", "password", "tokenhash", "created_at", "updated_at"}).AddRow(user.Id, user.Email, user.Name, user.Password, user.TokenHash, user.CreatedAt, user.UpdatedAt).ToPgxRows()

			token, err := services.GenerateAccessToken(user)
			assert.Nil(t, err, "Testing login")

			test.expectedBody = utils.UnsafetyToJson(models.TokensResponse{
				AccessToken: models.AccessToken{
					Access_token: token,
				},
			})
		}
		conn.Mock.EXPECT().Query(gomock.Any(), gomock.Any(), user.Email).Return(rows, sqlErr).Times(1)
		conn.Mock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).MaxTimes(1)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(test.body))
		assert.Nil(t, err, "Testing login")

		req = req.WithContext(context.WithValue(req.Context(), models.UserContextToken{}, candidate))

		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)

		assert.Equal(t, test.expectedCode, res.Code, "Testing login")
		assert.Equal(t, string(test.expectedBody), res.Body.String(), "Testing login")
	}
}
