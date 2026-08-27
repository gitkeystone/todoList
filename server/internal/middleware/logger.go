package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 访问日志中间件（JSON 行格式，PRD §6.4）。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf(`{"method":%q,"path":%q,"status":%d,"latencyMs":%d,"clientIP":%q}`,
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(),
			time.Since(start).Milliseconds(), c.ClientIP())
	}
}
