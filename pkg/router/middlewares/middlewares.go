package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/youssef-182/vue-skeleton-server/pkg/db"
	"github.com/youssef-182/vue-skeleton-server/pkg/models"
	"net/http"
	"os"
	"strings"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "You are unauthorized to perform this action.", http.StatusUnauthorized)
			return
		}

		token, found := strings.CutPrefix(authorizationHeader, "Bearer ")
		if !found {
			http.Error(w, "You are either unauthorized to perform this action or your authorization header setting is wrong.", http.StatusUnauthorized)
			return
		}

		tokenData, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if claims, ok := tokenData.Claims.(jwt.MapClaims); ok && tokenData.Valid {
			userID := claims["id"].(float64)
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	})
}

func BookCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bookID := chi.URLParam(r, "bookID")

		var book models.BookResp
		query := `SELECT 
    		books.id,
    		books.title,
    		books.alternative_title,
    		books.description,
    		books.author,
    		books.release_date,
    		user.id as "user.id",
    		user.first_name as "user.first_name",
    		user.last_name as "user.last_name",
    		user.username as "user.username",
    		user.email_address as "user.email_address",
    		user.birthdate as "user.birthdate",
    		user.gender as "user.gender"
    	FROM books JOIN users as user on user.id = books.user_id WHERE books.id=?`
		if err := db.DB.Get(&book, query, bookID); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"message": fmt.Sprintf("An error has occurred while trying to fetch the book. %v", err),
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		ctx := context.WithValue(r.Context(), "book", book)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
