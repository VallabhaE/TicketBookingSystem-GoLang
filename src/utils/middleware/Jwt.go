package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


// CustomClaims defines the structure of the JWT payload
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func Jwt_sign(name string) string {
	AddingTime := time.Now().Add(time.Hour)
	claims := CustomClaims{
		UserID: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-auth",
			ExpiresAt: jwt.NewNumericDate(AddingTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "auth_token",
		},
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SIGNING_KEY)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func GetUserId(r *http.Request) int {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	token, _ := jwt.Parse(parts[1], func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SIGNING_KEY, nil
	})

	id, err := strconv.Atoi(token.Claims.(jwt.MapClaims)["user_id"].(string))
	if err != nil {
		panic(err)
	}
	return id
}
