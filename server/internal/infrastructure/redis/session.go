// Package redis 提供 Redis 基础设施实现（构造注入，不依赖 global）。
package redis

import (
	"context"
	"errors"
	"time"

	"server/internal/domain/auth"
	"server/internal/domain/user"

	"github.com/go-redis/redis"
)

// SessionStore 基于 Redis 的 JWT 会话存储（uuid → refresh token）。
type SessionStore struct {
	client *redis.Client
	ttl    time.Duration
}

// NewSessionStore 构造会话存储；ttl 为 refresh token 过期时长。
func NewSessionStore(client *redis.Client, ttl time.Duration) *SessionStore {
	return &SessionStore{client: client, ttl: ttl}
}

// Set 保存会话。
func (s *SessionStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return s.client.Set(key, value, ttl).Err()
}

// Get 读取会话；不存在返回 auth.ErrSessionNotFound。
func (s *SessionStore) Get(ctx context.Context, key string) (string, error) {
	val, err := s.client.Get(key).Result()
	if errors.Is(err, redis.Nil) {
		return "", auth.ErrSessionNotFound
	}
	return val, err
}

// Del 删除会话。
func (s *SessionStore) Del(ctx context.Context, key string) error {
	return s.client.Del(key).Err()
}

// Cache 基于 Redis 的通用键值缓存。
type Cache struct {
	client *redis.Client
}

// NewCache 构造通用缓存。
func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

// Get 读取缓存值；不存在返回 redis.Nil（由调用方按 miss 处理）。
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(key).Result()
}

// Set 写入缓存。
func (c *Cache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(key, value, ttl).Err()
}

// Del 删除缓存键。
func (c *Cache) Del(ctx context.Context, key string) error {
	return c.client.Del(key).Err()
}

// 接口编译期校验。
var (
	_ auth.SessionStore = (*SessionStore)(nil)
	_ user.Cache        = (*Cache)(nil)
)
