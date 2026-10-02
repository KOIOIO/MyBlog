// Package database 提供数据库实体与 DTO（持久化映射层）。
package database

import (
	"time"

	"gorm.io/gorm"
)

// MODEL 通用模型字段（原 global.MODEL，随 global 包删除迁入）。
type MODEL struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
