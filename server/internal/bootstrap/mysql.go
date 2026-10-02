package bootstrap

import (
	"strings"

	"server/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// mysqlDriver 根据配置构造 MySQL GORM 驱动。
func mysqlDriver(cfg config.Mysql) gorm.Dialector {
	return mysql.Open(cfg.Dsn())
}

// mysqlLogger 根据配置返回对应日志级别的 GORM logger。
func mysqlLogger(cfg config.Mysql) logger.Interface {
	return logger.Default.LogMode(cfg.LogLevel())
}

// LogLevel 根据配置字符串返回 GORM 日志级别（与 config.Mysql 内部逻辑一致，便于独立调用）。
func LogLevel(mode string) logger.LogLevel {
	switch strings.ToLower(mode) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}
