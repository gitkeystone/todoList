// Package service 业务逻辑层：校验与业务规则（PRD §6.5）。
package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cxh/todolist/server/internal/apperr"
	"github.com/cxh/todolist/server/internal/model"
	"github.com/cxh/todolist/server/internal/repository"
)

const (
	maxTitleLen     = 200 // 标题最大长度（PRD FR-01）
	maxKeywordLen   = 100 // 搜索关键词最大长度（PRD §6.1）
	defaultPageSize = 20  // 默认每页条数（PRD §6.2）
	maxPageSize     = 100 // 每页条数上限
)

// PageMeta 分页元信息（PRD §6.2）。
type PageMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

// TodoService 待办业务接口。
type TodoService interface {
	Create(title string) (*model.Todo, error)
	Get(id uint) (*model.Todo, error)
	List(status, q, sort string, page, pageSize int) ([]model.Todo, *PageMeta, error)
	Patch(id uint, title *string, completed *bool) (*model.Todo, error)
	Delete(id uint) error
	DeleteCompleted() error
}

type todoService struct {
	repo repository.TodoRepository
}

// NewTodoService 构造业务服务。
func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) Create(title string) (*model.Todo, error) {
	title = strings.TrimSpace(title)
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	todo := &model.Todo{Title: title}
	if err := s.repo.Create(todo); err != nil {
		return nil, dbError(err)
	}
	return todo, nil
}

func (s *todoService) Get(id uint) (*model.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, dbError(err)
	}
	return todo, nil
}

func (s *todoService) List(status, q, sort string, page, pageSize int) ([]model.Todo, *PageMeta, error) {
	if err := validateStatus(status); err != nil {
		return nil, nil, err
	}
	if page < 1 || pageSize < 1 || pageSize > maxPageSize {
		return nil, nil, apperr.ErrInvalidPagination
	}
	q = strings.TrimSpace(q)
	if utf8.RuneCountInString(q) > maxKeywordLen {
		return nil, nil, fmt.Errorf("%w: 搜索关键词不能超过 %d 字符", apperr.ErrBadRequest, maxKeywordLen)
	}
	orderBy, err := buildOrderBy(sort)
	if err != nil {
		return nil, nil, err
	}

	query := repository.ListQuery{
		Status:   status,
		Q:        q,
		Page:     page,
		PageSize: pageSize,
		Sort:     orderBy,
	}
	todos, err := s.repo.List(query)
	if err != nil {
		return nil, nil, dbError(err)
	}
	total, err := s.repo.Count(query)
	if err != nil {
		return nil, nil, dbError(err)
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}
	meta := &PageMeta{Page: page, PageSize: pageSize, Total: total, TotalPages: totalPages}
	return todos, meta, nil
}

func (s *todoService) Patch(id uint, title *string, completed *bool) (*model.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, dbError(err)
	}

	changed := false
	if title != nil {
		t := strings.TrimSpace(*title)
		if err := validateTitle(t); err != nil {
			return nil, err
		}
		if t != todo.Title {
			todo.Title = t
			changed = true
		}
	}
	if completed != nil {
		if *completed != todo.Completed {
			now := time.Now().UTC()
			if *completed {
				todo.CompletedAt = &now // 完成：记录完成时间
			} else {
				todo.CompletedAt = nil // 取消完成：清空完成时间
			}
			todo.Completed = *completed
			changed = true
		}
	}
	if !changed {
		return todo, nil // 无字段变化：不写库，直接返回现有数据
	}
	if err := s.repo.Update(todo); err != nil {
		return nil, dbError(err)
	}
	return todo, nil
}

func (s *todoService) Delete(id uint) error {
	err := s.repo.Delete(id)
	if errors.Is(err, repository.ErrNotFound) {
		return apperr.ErrNotFound
	}
	if err != nil {
		return dbError(err)
	}
	return nil
}

func (s *todoService) DeleteCompleted() error {
	if _, err := s.repo.DeleteCompleted(); err != nil {
		return dbError(err)
	}
	return nil
}

// validateTitle 标题校验：trim 后 1~200 字符（PRD FR-01）。
func validateTitle(title string) error {
	if n := utf8.RuneCountInString(title); n < 1 || n > maxTitleLen {
		return apperr.ErrInvalidTitle
	}
	return nil
}

// validateStatus 状态筛选白名单校验。
func validateStatus(status string) error {
	switch status {
	case "all", "active", "completed":
		return nil
	}
	return fmt.Errorf("%w: status 参数非法（仅支持 all/active/completed）", apperr.ErrBadRequest)
}

// sortFieldMap 排序字段白名单：前端字段名 → 数据库列名（PRD §6.2）。
var sortFieldMap = map[string]string{
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
	"completedAt": "completed_at",
	"completed":   "completed",
}

// buildOrderBy 解析并白名单校验排序参数，返回 ORDER BY 子句。
// 支持多字段：sort=completed:asc,createdAt:desc；缺省为 createdAt:desc。
func buildOrderBy(sort string) (string, error) {
	if strings.TrimSpace(sort) == "" {
		sort = "createdAt:desc"
	}
	parts := strings.Split(sort, ",")
	clauses := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		field, dir, ok := strings.Cut(p, ":")
		if !ok {
			return "", apperr.ErrInvalidSort
		}
		col, ok := sortFieldMap[strings.TrimSpace(field)]
		if !ok {
			return "", apperr.ErrInvalidSort
		}
		d := strings.ToLower(strings.TrimSpace(dir))
		if d != "asc" && d != "desc" {
			return "", apperr.ErrInvalidSort
		}
		clauses = append(clauses, col+" "+strings.ToUpper(d))
	}
	return strings.Join(clauses, ", "), nil
}

// dbError 包装数据库错误为业务错误 50001，并保留底层错误信息便于排查。
func dbError(err error) error {
	return fmt.Errorf("%w: %v", apperr.ErrDB, err)
}
