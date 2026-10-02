// Package config 提供系统配置领域模型（持久化端口 + 差异规则）。
package config

import "context"

// ConfigStore 配置持久化端口（保存 YAML 到文件）。
type ConfigStore interface {
	Save(ctx context.Context) error
}

// DiffArrays 计算两个数组的差异：added 为新增元素，removed 为移除元素
// （行为与原 utils.DiffArrays 一致）。
func DiffArrays(oldArray, newArray []string) (added []string, removed []string) {
	oldSet := make(map[string]struct{}, len(oldArray))
	for _, v := range oldArray {
		oldSet[v] = struct{}{}
	}
	newSet := make(map[string]struct{}, len(newArray))
	for _, v := range newArray {
		newSet[v] = struct{}{}
	}
	for _, v := range newArray {
		if _, exists := oldSet[v]; !exists {
			added = append(added, v)
		}
	}
	for _, v := range oldArray {
		if _, exists := newSet[v]; !exists {
			removed = append(removed, v)
		}
	}
	return added, removed
}
