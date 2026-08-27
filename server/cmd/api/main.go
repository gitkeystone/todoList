// Todo List 后端服务入口（M1：完整分层组装 + 优雅退出）。
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gitkeystone/todolist/server/internal/config"
	"github.com/gitkeystone/todolist/server/internal/handler"
	"github.com/gitkeystone/todolist/server/internal/model"
	"github.com/gitkeystone/todolist/server/internal/repository"
	"github.com/gitkeystone/todolist/server/internal/router"
	"github.com/gitkeystone/todolist/server/internal/service"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	if cfg.GinMode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库（SQLite + 迁移，PRD §7）
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("init database: %v", err)
	}

	// 分层组装：repository → service → handler → router
	repo := repository.NewTodoRepository(db)
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	r := gin.New()
	router.Setup(r, h, cfg.AllowedOrigins)
	// 生产模式：单进程托管前端构建产物（PRD §5.1）
	if cfg.GinMode == gin.ReleaseMode {
		router.ServeFrontend(r, cfg.WebDist)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Todo List API listening on :%s (mode=%s)", cfg.Port, cfg.GinMode)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server exited: %v", err)
		}
	}()

	// 优雅退出：等待 SIGINT/SIGTERM，5 秒内完成收尾
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
	log.Println("server stopped")
}

// initDB 打开 SQLite（含 PRAGMA 与连接池限制）并执行迁移。
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// 确保数据目录存在（首次运行时 data/ 尚未创建）
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// SQLite 单写者：限制连接数，避免 "database is locked"（PRD §7.1）
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	if err := model.Init(db); err != nil {
		return nil, err
	}
	return db, nil
}
