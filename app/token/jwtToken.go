package token

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v4"

	"go-storage/pkg/gserr"
)

type GsClaims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GetUserId(ctx context.Context) (uint64, error) {
	userId, err := ctx.Value("user_id").(json.Number).Int64()
	if err != nil {
		return 0, fmt.Errorf("%w: %w", gserr.ErrServerCommon, err)
	}
	return uint64(userId), nil
}
