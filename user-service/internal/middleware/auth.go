package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	transport_http_response "github.com/Daty26/order-system/user-service/internal/transport/http/response"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			transport_http_response.ErrorJSON(w, http.StatusUnauthorized, "invalid token")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			transport_http_response.ErrorJSON(w, http.StatusUnauthorized, "couldn't find the token")
			return
		}
		tokenString := parts[1]
		claims, err := verifyToken(tokenString)
		if err != nil {
			transport_http_response.ErrorJSON(w, http.StatusUnauthorized, err.Error())
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
