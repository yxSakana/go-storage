package svc

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	_cacheRegisterInfoKeyf = "gsCache:register:"

	CacheRegisterInfoExpireSec = 5 * 60 * 60
)

type RegisterCacheInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type RegisterCache struct {
	rdb *redis.Redis
}

func NewRegisterCache(rdb *redis.Redis) *RegisterCache {
	return &RegisterCache{rdb: rdb}
}

func (r *RegisterCache) Exist(email string) bool {
	exist, err := r.rdb.Exists(_cacheRegisterInfoKeyf + email)
	return err == nil && exist
}

func (r *RegisterCache) Save(ctx context.Context, info *RegisterCacheInfo) error {
	key := _cacheRegisterInfoKeyf + info.Email
	val, _ := json.Marshal(info)
	return r.rdb.SetexCtx(ctx, key, string(val), CacheRegisterInfoExpireSec)
}

func (r *RegisterCache) Load(ctx context.Context, email string) (*RegisterCacheInfo, error) {
	key := _cacheRegisterInfoKeyf + email
	val, err := r.rdb.GetCtx(ctx, key)
	if err != nil {
		return nil, err
	}
	var info RegisterCacheInfo
	if err := json.Unmarshal([]byte(val), &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *RegisterCache) Delete(ctx context.Context, email string) error {
	key := _cacheRegisterInfoKeyf + email
	_, err := r.rdb.DelCtx(ctx, key)
	return err
}
