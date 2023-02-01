package middlewares

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"jwt/internal/models"
	"jwt/pkg/helpers/utils"
	"net/http"
)

func ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteError(w, err.Error(), http.StatusBadRequest, 0)
			return
		}
		r.Body.Close()

		var candidate *models.User
		if err := json.Unmarshal(data, &candidate); err != nil {
			utils.WriteError(w, err.Error(), http.StatusBadRequest, 0)
			return
		}

		if !candidate.Valid() {
			utils.WriteError(w, "Invalid body fields", http.StatusBadRequest, 0)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserContextToken{}, candidate)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
