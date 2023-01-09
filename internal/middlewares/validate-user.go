package middlewares

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"jwt/internal/models"
	"net/http"
)

func ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		r.Body.Close()

		var candidate models.User
		if err := json.Unmarshal(data, &candidate); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !candidate.Valid() {
			http.Error(w, "Invalid body fields", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserToken{}, candidate)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
