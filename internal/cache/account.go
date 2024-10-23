package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	CacheTokenPrefix    = "cacheTokenPrefix"
	CacheUserInfoPrefix = "cacheUserInfoPrefix"
)

type AccountCache interface {
	SetTokenCache(userId int64, token string, expired int) error
	GetTokenCache(userId int64) string

	SetUserInfoCache(userId int64, info map[string]interface{}, expired int) error
	GetUserInfoCache(userId int64) (map[string]interface{}, error)
	GetUserPermissions(userId int64) ([]string, error)
}

type accountCache struct {
	rdb *redis.Client
	ctx context.Context
}

func NewAccountCache(rdb *redis.Client) AccountCache {
	return &accountCache{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (a *accountCache) SetTokenCache(userId int64, token string, expired int) error {
	key := fmt.Sprintf("%v:%v", CacheTokenPrefix, userId)
	return a.rdb.Set(a.ctx, key, token, time.Duration(expired)*time.Second).Err()
}

func (a *accountCache) GetTokenCache(userId int64) string {
	key := fmt.Sprintf("%v:%v", CacheTokenPrefix, userId)
	return a.rdb.Get(a.ctx, key).String()
}

func (a *accountCache) SetUserInfoCache(userId int64, info map[string]interface{}, expired int) error {
	key := fmt.Sprintf("%v:%v", CacheUserInfoPrefix, userId)
	marshal, _ := json.Marshal(info)
	return a.rdb.Set(a.ctx, key, string(marshal), time.Duration(expired)*time.Second).Err()
}

func (a *accountCache) GetUserInfoCache(userId int64) (map[string]interface{}, error) {

	key := fmt.Sprintf("%v:%v", CacheUserInfoPrefix, userId)
	result, err := a.rdb.Get(a.ctx, key).Result()
	if err != nil || result == "" {
		return nil, errors.New("userinfo not found")
	}
	user := make(map[string]interface{})
	if err := json.Unmarshal([]byte(result), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (a *accountCache) GetUserPermissions(userId int64) ([]string, error) {
	key := fmt.Sprintf("%v:%v", CacheUserInfoPrefix, userId)
	result, err := a.rdb.Get(a.ctx, key).Result()
	if err != nil || result == "" {
		return nil, errors.New("userinfo not found")
	}
	user := make(map[string]interface{})
	if err := json.Unmarshal([]byte(result), &user); err != nil {
		return nil, err
	}
	if user["permissions"] != nil {
		var perm []string
		for _, v := range user["permissions"].([]interface{}) {
			perm = append(perm, v.(string))
		}
		return perm, nil
	}
	return nil, nil
}
