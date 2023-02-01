package middlewares

import (
	"bytes"
	"encoding/json"
	"jwt/internal/models"
	"jwt/pkg/helpers/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type userBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	OddField string `json:"oddField"`
}

type testCase struct {
	Body                 []byte
	ExpectedCode         int
	ExpectedBody         []byte
	ExpectedInternalCode int
}

var testCases []testCase

func init() {
	invalidBodyJson := utils.UnsafetyToJson(models.ResponseError{
		Message:      "Invalid body fields",
		Status:       http.StatusBadRequest,
		InternalCode: 0,
		Stack:        "",
	})

	validBodyJson := utils.UnsafetyToJson(models.User{
		Email:    "test_testovich",
		Password: "test_testovich",
	})

	testCases = []testCase{
		{
			Body:         []byte(`{"email": "test_testovich", "password": "test_testovich"}`),
			ExpectedCode: http.StatusOK,
			ExpectedBody: validBodyJson,
		},
		{
			Body:         []byte(`{"email": "test_testovich", "password": ""}`),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: invalidBodyJson,
		},
		{
			Body:         []byte(`{"email": "", "password": "test_testovich"}`),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: invalidBodyJson,
		},
		{
			Body:         []byte(`{"email": "test_testovich"}`),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: invalidBodyJson,
		},
		{
			Body:         []byte(`{"password": "test_testovich"}`),
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: invalidBodyJson,
		},
	}
}

func TestValidateUserMiddleware(t *testing.T) {
	for _, test := range testCases {
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(test.Body))
		res := httptest.NewRecorder()
		mw := ValidateUser(next(t))
		mw.ServeHTTP(res, req)

		assert.Equal(t, test.ExpectedCode, res.Code, "Testing login and password moddleware")
		assert.Equal(t, string(test.ExpectedBody), res.Body.String(), "Testing login and password moddleware")
	}
}

func next(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		candidate := r.Context().Value(models.UserContextToken{}).(*models.User)

		userJson, err := json.Marshal(candidate)
		assert.Nil(t, err, "Marshal candidate to json")

		w.WriteHeader(http.StatusOK)
		w.Write(userJson)
	})
}
