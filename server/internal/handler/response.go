// Package handler HTTP 处理器层：参数绑定与响应封装（PRD §6）。
package handler

import (
	"errors"
	"net/http"

	"github.com/cxh/todolist/server/internal/apperr"
	"github.com/gin-gonic/gin"
)

// Envelope 统一响应信封（PRD §6.2）。
type Envelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta"`
}

// OK 返回 200 成功响应。
func OK(c *gin.Context, data, meta any) {
	c.JSON(http.StatusOK, Envelope{Code: apperr.CodeOK, Message: "ok", Data: data, Meta: meta})
}

// Created 返回 201 创建成功响应（带 Location 头）。
func Created(c *gin.Context, data any, location string) {
	if location != "" {
		c.Header("Location", location)
	}
	c.JSON(http.StatusCreated, Envelope{Code: apperr.CodeOK, Message: "ok", Data: data, Meta: nil})
}

// NoContent 返回 204 无内容响应。
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Fail 返回错误响应；未知错误统一映射为 50000 服务器内部错误。
func Fail(c *gin.Context, err error) {
	var ae *apperr.AppError
	if errors.As(err, &ae) {
		c.JSON(ae.HTTPStatus, Envelope{Code: ae.Code, Message: ae.Message, Data: nil, Meta: nil})
		return
	}
	c.JSON(http.StatusInternalServerError,
		Envelope{Code: apperr.CodeInternal, Message: apperr.ErrInternal.Message, Data: nil, Meta: nil})
}
