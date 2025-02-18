package global

import (
	"gorm.io/gorm"
	"time"
)

type MODEL struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeleteAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
