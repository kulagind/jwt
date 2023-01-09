package routers

import (
	"jwt/internal/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRequest(mux *mux.Router) {
	// handleUsers(mux.PathPrefix("/users").Subrouter())
	postRequests := mux.Methods(http.MethodPost).Subrouter()
	postRequests.HandleFunc("/signup", signUp)
	postRequests.HandleFunc("/login", login)
	postRequests.Use(middlewares.ValidateUser)

}
