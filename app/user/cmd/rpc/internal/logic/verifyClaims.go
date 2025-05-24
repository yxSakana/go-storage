package logic

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"go-storage/app/user/cmd/rpc/internal/svc"
)

type verifyClaims struct {
	svc.RegisterCacheInfo
	jwt.RegisteredClaims
}

func GenerateVerifyToken(info svc.RegisterCacheInfo, expireAt time.Duration, secret string) (string, error) {
	now := time.Now()
	claims := verifyClaims{
		RegisterCacheInfo: info,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expireAt)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secret))
}

func ParseVerifyToken(tokenString, secret string) (*svc.RegisterCacheInfo, error) {
	jwtResult, err := jwt.ParseWithClaims(tokenString, &verifyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := jwtResult.Claims.(*verifyClaims); ok && jwtResult.Valid {
		return &svc.RegisterCacheInfo{
			Email:    claims.Email,
			Password: claims.Password,
			Token:    claims.Token,
		}, nil
	} else {
		return nil, err
	}
}
