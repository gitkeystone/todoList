// Package router 路由注册（PRD §6.1）。
package router

import (
	"net/http"

	"github.com/cxh/todolist/server/internal/handler"
	"github.com/cxh/todolist/server/internal/middleware"
	"github.com/gin-gonic/gin"
)

// Setup 组装中间件与路由。
// 注意：静态路由 /todos/completed 必须先于参数路由 /todos/:id 注册（TC-23，PRD §6.1）。
func Setup(r *gin.Engine, h *handler.TodoHandler, allowedOrigins []string) {
	r.Use(middleware.Logger(), middleware.Recovery(), middleware.CORS(allowedOrigins))

	// 健康检查探针（API 前缀外）
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	g := r.Group("/api/v1")
	g.POST("/todos", h.Create)
	g.GET("/todos", h.List)
	// 静态路由优先于参数路由
	g.DELETE("/todos/completed", h.DeleteCompleted)
	g.GET("/todos/:id", h.Get)
	g.PATCH("/todos/:id", h.Patch)
	g.DELETE("/todos/:id", h.Delete)
}
