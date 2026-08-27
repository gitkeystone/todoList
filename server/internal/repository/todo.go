// Package repository 数据访问层（GORM 实现，PRD §7）。
package repository

import (
	"errors"
	"strings"

	"github.com/cxh/todolist/server/internal/model"
	"gorm.io/gorm"
)

// ErrNotFound 记录不存在（repository 层哨兵错误，service 层负责映射为业务错误）。
var ErrNotFound = errors.New("record not found")

// ListQuery 列表查询条件（字段已由 service 层校验）。
type ListQuery struct {
	Status   string // all / active / completed
	Q        string // 标题模糊搜索（已转义）
	Page     int    // >= 1
	PageSize int    // 1~100
	Sort     string // 已白名单校验的排序子句，如 "created_at DESC"
}

// TodoRepository 待办数据访问接口。
type TodoRepository interface {
	Create(todo *model.Todo) error
	GetByID(id uint) (*model.Todo, error)
	Update(todo *model.Todo) error
	Delete(id uint) error
	DeleteCompleted() (int64, error)
	List(q ListQuery) ([]model.Todo, error)
	Count(q ListQuery) (int64, error)
}

type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository 构造 GORM 实现。
func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *model.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) GetByID(id uint) (*model.Todo, error) {
	var todo model.Todo
	if err := r.db.First(&todo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Update(todo *model.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	res := r.db.Delete(&model.Todo{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *todoRepository) DeleteCompleted() (int64, error) {
	res := r.db.Where("completed = ?", true).Delete(&model.Todo{})
	return res.RowsAffected, res.Error
}

func (r *todoRepository) List(q ListQuery) ([]model.Todo, error) {
	var todos []model.Todo
	err := r.db.
		Scopes(r.scopeFilter(q)).
		Order(q.Sort).
		Limit(q.PageSize).
		Offset((q.Page - 1) * q.PageSize).
		Find(&todos).Error
	return todos, err
}

func (r *todoRepository) Count(q ListQuery) (int64, error) {
	var total int64
	err := r.db.Model(&model.Todo{}).Scopes(r.scopeFilter(q)).Count(&total).Error
	return total, err
}

// scopeFilter 组装状态筛选与关键词搜索条件。
func (r *todoRepository) scopeFilter(q ListQuery) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch q.Status {
		case "active":
			db = db.Where("completed = ?", false)
		case "completed":
			db = db.Where("completed = ?", true)
		}
		if q.Q != "" {
			// LIKE 通配符已转义，需配合 ESCAPE 子句使其生效（防止搜索结果失真）
			db = db.Where("title LIKE ? ESCAPE '\\'", "%"+escapeLike(q.Q)+"%")
		}
		return db
	}
}

// escapeLike 转义 LIKE 通配符，防止搜索结果失真。
func escapeLike(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(s)
}
