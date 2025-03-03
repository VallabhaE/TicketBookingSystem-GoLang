package middleware


import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SIGNING_KEY = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var Bearer = "Bearer"

func MethodCheck(expectedMethod string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectedMethod {
			http.Error(w, ("Invalid request method"),
				http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := ExtractTokenFromHeader(r)
		if err != nil {
			http.Error(w, ("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		_, err = ValidateToken(tokenString, SIGNING_KEY)

		if err != nil {
			http.Error(w, ("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateToken(tokenString string, signingKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)

			if expirationTime.Before(time.Now()) {
				return nil, fmt.Errorf("token has expired")
			}

			return token, nil
		}
	}

	return nil, fmt.Errorf("invalid token claims")
}

func ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != Bearer {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return parts[1], nil
}
