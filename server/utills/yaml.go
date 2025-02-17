package utills

import (
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"server/global"
)

const configFile = "config.yaml"

func LoadYAML() ([]byte, error) {
	return os.ReadFile(configFile)
}
func SaveYAML() error {
	byteData, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, byteData, fs.ModePerm)
}
