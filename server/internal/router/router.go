// Package router 路由注册（PRD §6.1）。
package router

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gitkeystone/todolist/server/internal/handler"
	"github.com/gitkeystone/todolist/server/internal/middleware"
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

// ServeFrontend 托管前端构建产物（PRD §5.1 单进程部署）。
// dist 缺失时仅记录日志，不影响 API 服务。
func ServeFrontend(r *gin.Engine, webDist string) {
	dist := filepath.Clean(webDist)
	if _, err := os.Stat(filepath.Join(dist, "index.html")); err != nil {
		log.Printf("web dist not found at %s, skip static serving", dist)
		return
	}
	r.Static("/assets", filepath.Join(dist, "assets"))
	r.StaticFile("/favicon.svg", filepath.Join(dist, "favicon.svg"))
	r.StaticFile("/icons.svg", filepath.Join(dist, "icons.svg"))
	// SPA 回退：非 API 的 GET 请求返回 index.html（前端为单页应用，无路由）
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet && !strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.File(filepath.Join(dist, "index.html"))
			return
		}
		c.Status(http.StatusNotFound)
	})
}
