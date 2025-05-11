package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("your_jwt_secret"), nil
			})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Преобразуем Subject в uint и сохраняем в контексте
		userID, err := strconv.ParseUint(claims.Subject, 10, 32)
		if err != nil {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}