// Package model 定义数据模型与数据库迁移（PRD §7）。
package model

import (
	"time"

	"gorm.io/gorm"
)

// Todo 待办事项（PRD §7.2）。
type Todo struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Completed   bool       `gorm:"not null;default:false;index" json:"completed"`
	CreatedAt   time.Time  `gorm:"index" json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

// Init 执行数据库迁移（幂等，PRD §7.4）。
func Init(db *gorm.DB) error {
	return db.AutoMigrate(&Todo{})
}
