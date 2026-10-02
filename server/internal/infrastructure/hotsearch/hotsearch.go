// Package hotsearch 提供多平台热搜爬虫防腐层（百度/快手/头条/知乎/B站）。
// 爬虫逻辑为 utils/hotSearch 纯搬运，产物类型为 domain/website 数据形状。
package hotsearch

import (
	"context"
	"encoding/json"
	"time"

	"server/internal/domain/website"

	"github.com/go-redis/redis"
)

// Source 热搜数据源。
type Source interface {
	GetHotSearchData(maxNum int) (website.HotSearchData, error)
}

// NewSource 按数据源名构造爬虫。
func NewSource(sourceStr string) Source {
	switch sourceStr {
	case "baidu":
		return &Baidu{}
	case "kuaishou":
		return &Kuaishou{}
	case "toutiao":
		return &Toutiao{}
	case "zhihu":
		return &Zhihu{}
	case "bilibili":
		return &Bilibili{}
	default:
		return nil
	}
}

// Provider 热搜提供者（Redis 缓存 1 小时，miss 时爬取）。
type Provider struct {
	cache *redis.Client
}

// NewProvider 构造热搜提供者。
func NewProvider(cache *redis.Client) *Provider {
	return &Provider{cache: cache}
}

var _ website.HotSearchProvider = (*Provider)(nil)

// warmSources 定时预热的数据源。
var warmSources = []string{"baidu", "zhihu", "kuaishou", "toutiao"}

// GetHotSearchData 获取热搜数据（缓存优先）。
func (p *Provider) GetHotSearchData(_ context.Context, sourceStr string) (website.HotSearchData, error) {
	result, err := p.cache.Get(sourceStr).Result()
	if err != nil {
		source := NewSource(sourceStr)
		data, err := source.GetHotSearchData(30)
		if err != nil {
			return website.HotSearchData{}, err
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			return website.HotSearchData{}, err
		}
		if err := p.cache.Set(sourceStr, bytes, time.Hour).Err(); err != nil {
			return website.HotSearchData{}, err
		}
		return data, nil
	}
	var data website.HotSearchData
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return website.HotSearchData{}, err
	}
	return data, nil
}

// WarmAll 预热全部平台热搜。
func (p *Provider) WarmAll(_ context.Context) error {
	for _, sourceStr := range warmSources {
		source := NewSource(sourceStr)
		data, err := source.GetHotSearchData(30)
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if err := p.cache.Set(sourceStr, bytes, time.Hour).Err(); err != nil {
			return err
		}
	}
	return nil
}
