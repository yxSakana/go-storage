package token

import "github.com/golang-jwt/jwt/v4"

type GsClaims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
}
