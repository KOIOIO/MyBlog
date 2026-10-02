// Package config 提供配置加载与持久化能力（构造注入，不依赖 global）。
package config

import (
	"os"

	"server/config"

	"gopkg.in/yaml.v3"
)

// DefaultConfigFile 默认配置文件路径（相对进程工作目录）。
const DefaultConfigFile = "config.yaml"

// Load 从 path 指定的 YAML 文件加载配置并返回实例；path 为空时使用默认路径。
func Load(path string) (*config.Config, error) {
	if path == "" {
		path = DefaultConfigFile
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &config.Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save 将配置实例序列化并写回 path 指定的 YAML 文件；path 为空时使用默认路径。
func Save(path string, cfg *config.Config) error {
	if path == "" {
		path = DefaultConfigFile
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
