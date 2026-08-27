// Todo List 后端服务入口（M0 脚手架阶段：最小可运行 + /healthz 探针）。
package main

import (
	"log"
	"net/http"

	"github.com/cxh/todolist/server/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.GinMode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 健康检查探针（M0 联调冒烟用）
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Printf("Todo List API listening on :%s (mode=%s)", cfg.Port, cfg.GinMode)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
