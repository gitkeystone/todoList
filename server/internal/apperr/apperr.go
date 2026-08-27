// Package apperr 定义统一的业务错误类型与错误码（PRD §6.2）。
package apperr

import "net/http"

// 业务错误码（PRD §6.2 错误码表）
const (
	CodeOK                = 0     // 成功
	CodeBadRequest        = 40000 // 请求参数错误（通用）
	CodeInvalidTitle      = 40001 // 标题为空或超长（1~200）
	CodeInvalidPagination = 40002 // 分页参数非法
	CodeInvalidSort       = 40003 // 排序参数非法
	CodeInvalidID         = 40004 // id 不是合法正整数
	CodeNotFound          = 40400 // 待办不存在
	CodeInternal          = 50000 // 服务器内部错误
	CodeDB                = 50001 // 数据库操作失败
)

// AppError 业务错误：携带业务码、HTTP 状态码与可直接展示的中文信息。
type AppError struct {
	Code       int
	HTTPStatus int
	Message    string
}

func (e *AppError) Error() string { return e.Message }

// New 构造业务错误。
func New(code, httpStatus int, message string) *AppError {
	return &AppError{Code: code, HTTPStatus: httpStatus, Message: message}
}

// 预置错误实例（服务层与处理器层共用，errors.Is/As 可匹配）
var (
	ErrBadRequest        = New(CodeBadRequest, http.StatusBadRequest, "请求参数错误")
	ErrInvalidTitle      = New(CodeInvalidTitle, http.StatusBadRequest, "标题不能为空或超长（1~200 字符）")
	ErrInvalidPagination = New(CodeInvalidPagination, http.StatusBadRequest, "分页参数非法（page ≥ 1，pageSize 1~100）")
	ErrInvalidSort       = New(CodeInvalidSort, http.StatusBadRequest, "排序参数非法（field 仅支持 createdAt/updatedAt/completedAt）")
	ErrInvalidID         = New(CodeInvalidID, http.StatusBadRequest, "id 必须是正整数")
	ErrNotFound          = New(CodeNotFound, http.StatusNotFound, "待办不存在")
	ErrInternal          = New(CodeInternal, http.StatusInternalServerError, "服务器内部错误")
	ErrDB                = New(CodeDB, http.StatusInternalServerError, "数据库操作失败")
)
