package core

import (
	"gopkg.in/yaml.v3"
	"log"
	"server/config"
	"server/utills"
)

// InitConf 从 YAML 文件加载配置
func InitConf() *config.Config {
	c := &config.Config{}
	yamlConf, err := utills.LoadYAML()
	if err != nil {
		log.Fatalln("加载配置文件失败")
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalln("解析配置文件失败")
	}
	return c
}
