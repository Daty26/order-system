package middleware

import (
	"context"
	"errors"
	"github.com/Daty26/order-system/auth/response"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.ErrorResponse(w, http.StatusUnauthorized, "invalid token")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			response.ErrorResponse(w, http.StatusUnauthorized, "couldn't find the token")
			return
		}
		tokenString := parts[1]
		claims, err := verifyToken(tokenString)
		if err != nil {
			response.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		ctx = context.WithValue(ctx, "username", claims["username"])
		ctx = context.WithValue(ctx, "role", claims["role"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func verifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}
