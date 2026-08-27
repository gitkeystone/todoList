package repository

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cxh/todolist/server/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newTestDB 创建隔离的内存 SQLite（每个测试独立数据库，PRD §10.2）。
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open memory db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	if err := db.AutoMigrate(&model.Todo{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func seed(t *testing.T, repo TodoRepository, titles ...string) []model.Todo {
	t.Helper()
	var todos []model.Todo
	for _, title := range titles {
		todo := &model.Todo{Title: title}
		if err := repo.Create(todo); err != nil {
			t.Fatalf("seed %q: %v", title, err)
		}
		todos = append(todos, *todo)
	}
	return todos
}

func TestCreate(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepository(db)

	todo := &model.Todo{Title: "学习 Go 语言"}
	if err := repo.Create(todo); err != nil {
		t.Fatalf("create: %v", err)
	}
	if todo.ID == 0 {
		t.Fatal("expected auto-increment id")
	}
	if todo.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be set by gorm")
	}

	got, err := repo.GetByID(todo.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "学习 Go 语言" || got.Completed {
		t.Fatalf("unexpected todo: %+v", got)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepository(db)

	_, err := repo.GetByID(999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepository(db)
	created := seed(t, repo, "旧标题")[0]

	now := time.Now().UTC()
	created.Title = "新标题"
	created.Completed = true
	created.CompletedAt = &now
	if err := repo.Update(&created); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != "新标题" || !got.Completed || got.CompletedAt == nil {
		t.Fatalf("unexpected after update: %+v", got)
	}

	// 取消完成：completedAt 应被清空
	got.Completed = false
	got.CompletedAt = nil
	if err := repo.Update(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	got2, _ := repo.GetByID(created.ID)
	if got2.Completed || got2.CompletedAt != nil {
		t.Fatalf("expected completed=false and nil completed_at: %+v", got2)
	}
}

func TestDelete(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepository(db)
	created := seed(t, repo, "待删除")[0]

	if err := repo.Delete(created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := repo.GetByID(created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
	if err := repo.Delete(created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound on second delete, got %v", err)
	}
}

func TestDeleteCompleted(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepository(db)
	todos := seed(t, repo, "A", "B", "C", "D")

	now := time.Now().UTC()
	todos[0].Completed, todos[0].CompletedAt = true, &now
	todos[2].Completed, todos[2].CompletedAt = true, &now
	if err := repo.Update(&todos[0]); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := repo.Update(&todos[2]); err != nil {
		t.Fatalf("update: %v", err)
	}

	n, err := repo.DeleteCompleted()
	if err != nil {
		t.Fatalf("delete completed: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 deleted, got %d", n)
	}

	all, _ := repo.List(ListQuery{Status: "all", Page: 1, PageSize: 20, Sort: "id ASC"})
	if len(all) != 2 {
		t.Fatalf("expected 2 remaining, got %d", len(all))
	}
	for _, todo := range all {
		if todo.Completed {
			t.Fatalf("completed todo should be deleted: %+v", todo)
		}
	}
}

func TestListStatusFilter(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepository(db)
	todos := seed(t, repo, "A", "B", "C")

	now := time.Now().UTC()
	todos[1].Completed, todos[1].CompletedAt = true, &now
	if err := repo.Update(&todos[1]); err != nil {
		t.Fatalf("update: %v", err)
	}

	active, _ := repo.List(ListQuery{Status: "active", Page: 1, PageSize: 20, Sort: "id ASC"})
	if len(active) != 2 {
		t.Fatalf("active: expected 2, got %d", len(active))
	}
	completed, _ := repo.List(ListQuery{Status: "completed", Page: 1, PageSize: 20, Sort: "id ASC"})
	if len(completed) != 1 {
		t.Fatalf("completed: expected 1, got %d", len(completed))
	}
	all, _ := repo.List(ListQuery{Status: "all", Page: 1, PageSize: 20, Sort: "id ASC"})
	if len(all) != 3 {
		t.Fatalf("all: expected 3, got %d", len(all))
	}
}

func TestListKeyword(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepository(db)
	seed(t, repo, "学习 Go 语言", "给产品写周报", "进度 100% 完成")

	got, err := repo.List(ListQuery{Status: "all", Q: "Go", Page: 1, PageSize: 20, Sort: "id ASC"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 1 || got[0].Title != "学习 Go 语言" {
		t.Fatalf("keyword Go: %+v", got)
	}

	// LIKE 通配符转义：搜索字面量 "0%" 应精确匹配，不把 % 当通配符
	got2, _ := repo.List(ListQuery{Status: "all", Q: "0%", Page: 1, PageSize: 20, Sort: "id ASC"})
	if len(got2) != 1 || got2[0].Title != "进度 100% 完成" {
		t.Fatalf("keyword with percent: %+v", got2)
	}
}

func TestListPagination(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepository(db)
	seed(t, repo, "1", "2", "3", "4", "5")

	page1, _ := repo.List(ListQuery{Status: "all", Page: 1, PageSize: 2, Sort: "id ASC"})
	if len(page1) != 2 || page1[0].Title != "1" {
		t.Fatalf("page1: %+v", page1)
	}
	page3, _ := repo.List(ListQuery{Status: "all", Page: 3, PageSize: 2, Sort: "id ASC"})
	if len(page3) != 1 || page3[0].Title != "5" {
		t.Fatalf("page3: %+v", page3)
	}
}

func TestListSort(t *testing.T) {
	db := newTestDB(t)
	repo := NewTodoRepository(db)
	seed(t, repo, "A", "B", "C")

	// 默认 created_at DESC：最后创建的排在最前
	got, _ := repo.List(ListQuery{Status: "all", Page: 1, PageSize: 20, Sort: "created_at DESC"})
	if got[0].Title != "C" {
		t.Fatalf("desc sort: %+v", got)
	}
	// 显式 ASC
	gotAsc, _ := repo.List(ListQuery{Status: "all", Page: 1, PageSize: 20, Sort: "created_at ASC"})
	if gotAsc[0].Title != "A" {
		t.Fatalf("asc sort: %+v", gotAsc)
	}
}
