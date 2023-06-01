package middlewares

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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

		//token = string(regexp.MustCompile(`\s*$`).ReplaceAll([]byte(token), []byte{}))
		fmt.Println(token)

		claims := jwt.MapClaims{}
		tokenData, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if tokenData.Valid {
			userID := claims["id"].(int64)
			//username := claims["username"].(string)
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		http.Error(w, "Token is invalid", http.StatusUnauthorized)
		return
	})
}
