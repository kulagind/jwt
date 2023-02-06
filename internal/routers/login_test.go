package routers

import (
	"bytes"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/pkg/helpers/pg"
	"jwt/pkg/helpers/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type loginTestCase struct {
	expectedCode int
	body         []byte
	expectedBody []byte
	email        string
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

	validResponse := utils.UnsafetyToJson(models.TokensResponse{
		AccessToken: models.AccessToken{
			Access_token: "",
		},
	})

	testCases := []loginTestCase{
		{
			body:         []byte(`{"email": "test_testovich", "password": "test_testovich"}`),
			email:        "test_testovich",
			expectedCode: http.StatusOK,
			expectedBody: validResponse,
		},
		{
			body:         []byte(`{"email": "test_testovich", "password": ""}`),
			email:        "test_testovich",
			expectedCode: http.StatusUnauthorized,
			expectedBody: userNotFoundResponse,
		},
		{
			body:         []byte(`{"email": "", "password": "test_testovich"}`),
			email:        "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: userNotFoundResponse,
		},
		{
			body:         []byte(`{"email": "test_testovich"}`),
			email:        "test_testovich",
			expectedCode: http.StatusUnauthorized,
			expectedBody: userNotFoundResponse,
		},
		{
			body:         []byte(`{"password": "test_testovich"}`),
			email:        "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: userNotFoundResponse,
		},
	}

	// rows := pgxpoolmock.NewRows([]string{}).ToPgxRows()
	// sqlErr := errors.New(constants.NO_ROWS)

	for _, test := range testCases {
		// if test.expectedCode == http.StatusOK {
		// 	sqlErr = nil
		// 	rows = pgxpoolmock.NewRows([]string{"id", "email", "name", "password", "tokenhash", "created_at", "updated_at"}).AddRow("test_testovich", "test_testovich", "name", "test_testovich", "test_testovich", time.Now(), time.Now()).ToPgxRows()
		// }
		// conn.Mock.EXPECT().Query(gomock.Any(), gomock.Any(), test.email).Return(rows, sqlErr).Times(1)
		// conn.Mock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(test.body))
		assert.Nil(t, err, "Testing login")

		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)

		assert.Equal(t, test.expectedCode, res.Code, "Testing login")
		assert.Equal(t, string(test.expectedBody), res.Body.String(), "Testing login")
	}
}
