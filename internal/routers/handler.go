package routers

import (
	"jwt/internal/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRequest(mux *mux.Router) {
	// public API
	publicRouter := mux.Methods(http.MethodPost).Subrouter()
	publicRouter.HandleFunc("/signup", signUp)
	publicRouter.HandleFunc("/login", login)
	publicRouter.Use(middlewares.ValidateUser)

	publicAccessRouter := mux.Methods(http.MethodPost).Subrouter()
	publicAccessRouter.HandleFunc("/update_access", updateAccessToken)
	publicAccessRouter.Use(middlewares.ValidateRefreshToken, middlewares.UpdateRefreshTokenIfRequired)

	// private API
	privateRouter := mux.PathPrefix("/private").Subrouter()
	privateRouter.Use(middlewares.ValidateAccessToken)
	usersRouter := privateRouter.PathPrefix("/user").Methods(http.MethodGet).Subrouter()
	usersRouter.HandleFunc("", getCurrentUser)
	usersRouter.HandleFunc("/{id}", getUserById)
}
