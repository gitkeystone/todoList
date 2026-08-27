package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gitkeystone/todolist/server/internal/apperr"
	"github.com/gitkeystone/todolist/server/internal/service"
)

// TodoHandler 待办 HTTP 处理器。
type TodoHandler struct {
	svc service.TodoService
}

// NewTodoHandler 构造处理器。
func NewTodoHandler(svc service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

type createTodoRequest struct {
	Title string `json:"title"`
}

type patchTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

// List GET /api/v1/todos — 查询列表（筛选/搜索/分页/排序，PRD §6.3.1）。
func (h *TodoHandler) List(c *gin.Context) {
	status := c.DefaultQuery("status", "all")
	q := c.Query("q")
	sort := c.Query("sort")
	page, err := parsePositiveInt(c.Query("page"), 1)
	if err != nil {
		Fail(c, apperr.ErrInvalidPagination)
		return
	}
	pageSize, err := parsePositiveInt(c.DefaultQuery("pageSize", "20"), 20)
	if err != nil {
		Fail(c, apperr.ErrInvalidPagination)
		return
	}

	todos, meta, err := h.svc.List(status, q, sort, page, pageSize)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, todos, meta)
}

// Create POST /api/v1/todos — 创建待办（PRD §6.3.2）。
func (h *TodoHandler) Create(c *gin.Context) {
	var req createTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperr.New(apperr.CodeBadRequest, 400, "请求体格式错误"))
		return
	}
	todo, err := h.svc.Create(req.Title)
	if err != nil {
		Fail(c, err)
		return
	}
	Created(c, todo, "/api/v1/todos/"+strconv.FormatUint(uint64(todo.ID), 10))
}

// Get GET /api/v1/todos/:id — 查询单条（PRD §6.3.3）。
func (h *TodoHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		Fail(c, apperr.ErrInvalidID)
		return
	}
	todo, err := h.svc.Get(id)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, todo, nil)
}

// Patch PATCH /api/v1/todos/:id — 更新标题/完成状态（PRD §6.3.4）。
func (h *TodoHandler) Patch(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		Fail(c, apperr.ErrInvalidID)
		return
	}
	var req patchTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, apperr.New(apperr.CodeBadRequest, 400, "请求体格式错误"))
		return
	}
	if req.Title == nil && req.Completed == nil {
		Fail(c, apperr.New(apperr.CodeBadRequest, 400, "至少提供一个要更新的字段"))
		return
	}
	todo, err := h.svc.Patch(id, req.Title, req.Completed)
	if err != nil {
		Fail(c, err)
		return
	}
	OK(c, todo, nil)
}

// Delete DELETE /api/v1/todos/:id — 删除单条（PRD §6.3.5）。
func (h *TodoHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		Fail(c, apperr.ErrInvalidID)
		return
	}
	if err := h.svc.Delete(id); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}

// DeleteCompleted DELETE /api/v1/todos/completed — 清除已完成（PRD §6.3.6，幂等）。
func (h *TodoHandler) DeleteCompleted(c *gin.Context) {
	if err := h.svc.DeleteCompleted(); err != nil {
		Fail(c, err)
		return
	}
	NoContent(c)
}

// parsePositiveInt 解析正整数查询参数；空值回退默认值。
func parsePositiveInt(s string, def int) (int, error) {
	if s == "" {
		return def, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 0, errors.New("invalid positive int")
	}
	return n, nil
}

// parseID 解析路径参数 id（正整数，PRD §6.2）。
func parseID(c *gin.Context) (uint, bool) {
	n, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || n == 0 {
		return 0, false
	}
	return uint(n), true
}
