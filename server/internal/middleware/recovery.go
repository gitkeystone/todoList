package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/gitkeystone/todolist/server/internal/apperr"
)

// Recovery panic 恢复中间件：记录堆栈并返回统一错误信封（PRD §6.4）。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v\n%s", rec, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    apperr.CodeInternal,
					"message": "服务器内部错误",
					"data":    nil,
					"meta":    nil,
				})
			}
		}()
		c.Next()
	}
}
