package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Type Type `json:"type"`

	jwt.RegisteredClaims
}
