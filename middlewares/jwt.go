package middlewares

import (
	"github.com/dgrijalva/jwt-go"
)

type MyCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
