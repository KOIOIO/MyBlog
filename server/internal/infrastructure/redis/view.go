package redis

import (
	"context"
	"strconv"

	"server/internal/domain/article"

	"github.com/go-redis/redis"
)

// ViewCounter 文章浏览量计数器（Redis hash：article_views）。
type ViewCounter struct {
	client *redis.Client
	index  string
}

// NewViewCounter 构造浏览量计数器。
func NewViewCounter(client *redis.Client) *ViewCounter {
	return &ViewCounter{client: client, index: "article_views"}
}

// Set 浏览量 +1。
func (c *ViewCounter) Set(ctx context.Context, id string) error {
	num, _ := c.client.HGet(c.index, id).Int()
	num++
	return c.client.HSet(c.index, id, num).Err()
}

// GetInfo 取出全部计数。
func (c *ViewCounter) GetInfo(ctx context.Context) map[string]int {
	info := map[string]int{}
	maps := c.client.HGetAll(c.index).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val)
		info[id] = num
	}
	return info
}

// Clear 清空计数。
func (c *ViewCounter) Clear(ctx context.Context) {
	c.client.Del(c.index)
}

// 编译期断言。
var _ article.ViewCounter = (*ViewCounter)(nil)
