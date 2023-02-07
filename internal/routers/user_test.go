package routers

import (
	"context"
	"errors"
	"fmt"
	"jwt/internal/constants"
	"jwt/internal/models"
	"jwt/internal/repo"
	"jwt/pkg/helpers/pg"
	"jwt/pkg/helpers/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type userByIdTestCase struct {
	expectedCode int
	expectedBody []byte
	id           string
}

func TestGettingUserById(t *testing.T) {
	mux := mux.NewRouter()
	mux.HandleFunc("/user/{id}", getUserById).Methods("GET")

	conn, err := pg.NewMockConnection(t)
	if err != nil {
		log.Fatal(err.Error())
	}
	d := repo.Init(conn)
	defer conn.Close()
	defer d()

	userDoesntExistResponse := utils.UnsafetyToJson(models.ResponseError{
		Message:      "User with ID doesn't exist",
		Status:       http.StatusBadRequest,
		InternalCode: -1,
		Stack:        "",
	})

	id := "test_testovich_111"
	email := "test_testovich_444"
	password := "test_testovich_555"

	user := &models.UserResponse{
		Email: email,
		Id:    id,
	}
	userJson := utils.UnsafetyToJson(user)

	testCases := []userByIdTestCase{
		{
			id:           id,
			expectedCode: http.StatusOK,
			expectedBody: userJson,
		},
		{
			id:           password,
			expectedCode: http.StatusBadRequest,
			expectedBody: userDoesntExistResponse,
		},
		{
			id:           email,
			expectedCode: http.StatusBadRequest,
			expectedBody: userDoesntExistResponse,
		},
	}

	for _, test := range testCases {
		sqlErr := errors.New(constants.NO_ROWS)
		rows := pgxpoolmock.NewRows([]string{}).ToPgxRows()

		if test.expectedCode == http.StatusOK {
			sqlErr = nil
			rows = pgxpoolmock.NewRows([]string{"id", "email", "name"}).AddRow(user.Id, user.Email, user.Name).ToPgxRows()
		}
		conn.Mock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(rows, sqlErr).Times(1)

		req, err := http.NewRequest("GET", fmt.Sprintf("/user/%s", test.id), nil)
		assert.Nil(t, err, "Testing getting user by id")

		req = req.WithContext(context.WithValue(req.Context(), models.UserContextToken{}, user))

		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)

		assert.Equal(t, test.expectedCode, res.Code, "Testing getting user by id")
		assert.Equal(t, string(test.expectedBody), res.Body.String(), "Testing getting user by id")
	}
}
