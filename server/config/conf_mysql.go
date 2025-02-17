package config

import (
	"gorm.io/gorm/logger"
	"strings"
)

type Mysql struct {
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	Config       string `json:"config" yaml:"config"`
	DBName       string `json:"db_name" yaml:"db_name"`
	Username     string `json:"username" yaml:"username"`
	Password     string `json:"password" yaml:"password"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"`
	LogMode      string `json:"log_mode" yaml:"log_mode"`
}

func (m Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + string(m.Port) + ")/" + m.DBName + m.Config
}

func (m Mysql) LogLevel() logger.LogLevel {
	switch strings.ToLower(m.LogMode) {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}
