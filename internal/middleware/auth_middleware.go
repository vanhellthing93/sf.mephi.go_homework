package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/vanhellthing93/sf.mephi.go_homework/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	if err := godotenv.Load(); err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Warning loading .env:")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        jwtSecret := os.Getenv("JWT_SECRET")
        if jwtSecret == "" {
            panic("JWT_SECRET not set in environment")
        }
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
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