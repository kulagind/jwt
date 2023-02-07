package routers

import (
	"context"
	"errors"
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

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
)

type signUpTestCase struct {
	expectedCode int
	expectedBody []byte
	email        string
	password     string
}

func init() {
	services.LoadEnv()
	os.Setenv("APP_MODE", "testing")
}

func TestSignUp(t *testing.T) {
	mux := mux.NewRouter()
	mux.HandleFunc("/signup", signUp).Methods("POST")

	conn, err := pg.NewMockConnection(t)
	if err != nil {
		log.Fatal(err.Error())
	}
	d := repo.Init(conn)
	defer conn.Close()
	defer d()

	userExistsResponse := utils.UnsafetyToJson(models.ResponseError{
		Message:      "User with this email already exists",
		Status:       http.StatusBadRequest,
		InternalCode: 1,
		Stack:        "",
	})

	email := "test_testovich_444"
	password := "test_testovich_555"

	testCases := []signUpTestCase{
		{
			email:        email,
			password:     password,
			expectedCode: http.StatusCreated,
			expectedBody: []byte("{}"),
		},
		{
			email:        email,
			password:     password,
			expectedCode: http.StatusBadRequest,
			expectedBody: userExistsResponse,
		},
	}

	for _, test := range testCases {
		var sqlErr error
		sqlErr = nil

		user := &models.User{
			Email:    test.email,
			Password: test.password,
		}

		if test.expectedCode == http.StatusBadRequest {
			sqlErr = errors.New(pgerrcode.UniqueViolation)
		} else {
			sqlErr = nil
		}
		conn.Mock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, sqlErr).Times(1)

		req, err := http.NewRequest("POST", "/signup", nil)
		assert.Nil(t, err, "Testing sign up")

		req = req.WithContext(context.WithValue(req.Context(), models.UserContextToken{}, user))

		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)

		assert.Equal(t, test.expectedCode, res.Code, "Testing sign up")
		if test.expectedCode == http.StatusCreated {
			assert.Contains(t, res.Body.String(), test.email, "Testing sign up")
		} else {
			assert.Equal(t, string(test.expectedBody), res.Body.String(), "Testing sign up")
		}
	}
}
