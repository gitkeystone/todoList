// 演示数据注入命令（make seed，PRD §7.4）。
package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gitkeystone/todolist/server/internal/config"
	"github.com/gitkeystone/todolist/server/internal/model"
	"github.com/gitkeystone/todolist/server/internal/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// 确保数据目录存在
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		log.Fatalf("mkdir data dir: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := model.Init(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// 幂等：先清空再插入
	if err := db.Where("1 = 1").Delete(&model.Todo{}).Error; err != nil {
		log.Fatalf("clear existing todos: %v", err)
	}

	repo := repository.NewTodoRepository(db)
	samples := []string{
		"学习 Go 语言基础",
		"给产品写周报",
		"整理报销单据",
		"预约牙医复诊",
		"读完《设计心理学》",
		"提交季度 KPI 自评",
		"给父母买生日礼物",
		"复习 GORM 高级用法",
		"写一篇技术博客",
		"清理邮箱未读邮件",
		"制定下月健身计划",
		"Review M1 代码评审意见",
	}
	for i, title := range samples {
		todo := &model.Todo{Title: title}
		if i%3 == 1 { // 约 1/3 演示为已完成状态
			now := time.Now().UTC()
			todo.Completed = true
			todo.CompletedAt = &now
		}
		if err := repo.Create(todo); err != nil {
			log.Fatalf("create %q: %v", title, err)
		}
	}
	log.Printf("seeded %d demo todos", len(samples))
}
