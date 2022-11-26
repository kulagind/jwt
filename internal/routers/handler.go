package routers

import "github.com/gorilla/mux"

func HandleRequest(mux *mux.Router) {
	handleUsers(mux.PathPrefix("/users").Subrouter())
}
