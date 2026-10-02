// Package configfile 提供配置 YAML 文件持久化实现。
package configfile

import (
	"context"
	"io/fs"
	"os"

	"server/config"

	"gopkg.in/yaml.v3"
)

// Store 配置持久化（对注入的配置指针做 YAML 序列化写回，与旧 utils.SaveYAML 一致）。
type Store struct {
	cfg  *config.Config
	path string
}

// NewStore 构造配置持久化，path 为 config.yaml 路径。
func NewStore(cfg *config.Config, path string) *Store {
	return &Store{cfg: cfg, path: path}
}

// Save 保存配置为 YAML。
func (s *Store) Save(_ context.Context) error {
	byteData, err := yaml.Marshal(s.cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, byteData, fs.ModePerm)
}
