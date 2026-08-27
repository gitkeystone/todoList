package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cxh/todolist/server/internal/apperr"
	"github.com/cxh/todolist/server/internal/model"
	"github.com/cxh/todolist/server/internal/repository"
)

// stubRepo 内存假实现，用于纯业务逻辑单测（PRD §10.2）。
type stubRepo struct {
	todos       map[uint]*model.Todo
	nextID      uint
	updateCalls int
	lastQuery   repository.ListQuery
}

func newStubRepo() *stubRepo {
	return &stubRepo{todos: map[uint]*model.Todo{}, nextID: 1}
}

func (s *stubRepo) Create(todo *model.Todo) error {
	todo.ID = s.nextID
	s.nextID++
	s.todos[todo.ID] = todo
	return nil
}

func (s *stubRepo) GetByID(id uint) (*model.Todo, error) {
	if todo, ok := s.todos[id]; ok {
		return todo, nil
	}
	return nil, repository.ErrNotFound
}

func (s *stubRepo) Update(todo *model.Todo) error {
	s.updateCalls++
	s.todos[todo.ID] = todo
	return nil
}

func (s *stubRepo) Delete(id uint) error {
	if _, ok := s.todos[id]; !ok {
		return repository.ErrNotFound
	}
	delete(s.todos, id)
	return nil
}

func (s *stubRepo) DeleteCompleted() (int64, error) {
	var n int64
	for id, todo := range s.todos {
		if todo.Completed {
			delete(s.todos, id)
			n++
		}
	}
	return n, nil
}

func (s *stubRepo) List(q repository.ListQuery) ([]model.Todo, error) {
	s.lastQuery = q
	var out []model.Todo
	for _, todo := range s.todos {
		out = append(out, *todo)
	}
	return out, nil
}

func (s *stubRepo) Count(q repository.ListQuery) (int64, error) {
	return int64(len(s.todos)), nil
}

func newService() (TodoService, *stubRepo) {
	repo := newStubRepo()
	return NewTodoService(repo), repo
}

func mustCreate(t *testing.T, svc TodoService, title string) *model.Todo {
	t.Helper()
	todo, err := svc.Create(title)
	if err != nil {
		t.Fatalf("create %q: %v", title, err)
	}
	return todo
}

func TestCreateValidTrimsTitle(t *testing.T) {
	svc, _ := newService()
	todo := mustCreate(t, svc, "  学习 Go 语言  ")
	if todo.Title != "学习 Go 语言" {
		t.Fatalf("title should be trimmed, got %q", todo.Title)
	}
	if todo.Completed {
		t.Fatal("new todo should be active")
	}
}

func TestCreateEmptyTitle(t *testing.T) {
	svc, _ := newService()
	_, err := svc.Create("   ")
	if !errors.Is(err, apperr.ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestCreateTitleTooLong(t *testing.T) {
	svc, _ := newService()
	_, err := svc.Create(strings.Repeat("字", 201))
	if !errors.Is(err, apperr.ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
	// 恰好 200 字符应通过
	if _, err := svc.Create(strings.Repeat("字", 200)); err != nil {
		t.Fatalf("200 chars should pass: %v", err)
	}
}

func TestGetNotFound(t *testing.T) {
	svc, _ := newService()
	_, err := svc.Get(999)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPatchComplete(t *testing.T) {
	svc, repo := newService()
	todo := mustCreate(t, svc, "完成任务")
	completed := true

	got, err := svc.Patch(todo.ID, nil, &completed)
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if !got.Completed || got.CompletedAt == nil {
		t.Fatalf("expected completed with timestamp: %+v", got)
	}
	if repo.updateCalls != 1 {
		t.Fatalf("expected 1 update, got %d", repo.updateCalls)
	}
}

func TestPatchUncompleteClearsCompletedAt(t *testing.T) {
	svc, _ := newService()
	todo := mustCreate(t, svc, "先完成再取消")
	now := time.Now().UTC()
	todo.Completed = true
	todo.CompletedAt = &now

	completed := false
	got, err := svc.Patch(todo.ID, nil, &completed)
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if got.Completed || got.CompletedAt != nil {
		t.Fatalf("expected completed=false and nil completed_at: %+v", got)
	}
}

func TestPatchUpdateTitle(t *testing.T) {
	svc, _ := newService()
	todo := mustCreate(t, svc, "旧标题")
	newTitle := "新标题"

	got, err := svc.Patch(todo.ID, &newTitle, nil)
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if got.Title != "新标题" {
		t.Fatalf("title not updated: %+v", got)
	}
}

func TestPatchNoChangeSkipsWrite(t *testing.T) {
	svc, repo := newService()
	todo := mustCreate(t, svc, "标题不变")
	same := "标题不变"

	got, err := svc.Patch(todo.ID, &same, nil)
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if got.Title != "标题不变" {
		t.Fatalf("unexpected: %+v", got)
	}
	if repo.updateCalls != 0 {
		t.Fatalf("no-change patch should not write, updateCalls=%d", repo.updateCalls)
	}
}

func TestPatchInvalidTitle(t *testing.T) {
	svc, _ := newService()
	todo := mustCreate(t, svc, "合法标题")
	bad := ""

	_, err := svc.Patch(todo.ID, &bad, nil)
	if !errors.Is(err, apperr.ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestPatchNotFound(t *testing.T) {
	svc, _ := newService()
	title := "x"
	_, err := svc.Patch(999, &title, nil)
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	svc, _ := newService()
	todo := mustCreate(t, svc, "删除我")
	if err := svc.Delete(todo.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := svc.Delete(todo.ID); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteCompleted(t *testing.T) {
	svc, _ := newService()
	a := mustCreate(t, svc, "A")
	b := mustCreate(t, svc, "B")
	completed := true
	if _, err := svc.Patch(a.ID, nil, &completed); err != nil {
		t.Fatalf("complete A: %v", err)
	}
	if err := svc.DeleteCompleted(); err != nil {
		t.Fatalf("delete completed: %v", err)
	}
	// A 已被删除
	if _, err := svc.Get(a.ID); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("A should be gone: %v", err)
	}
	if _, err := svc.Get(b.ID); err != nil {
		t.Fatalf("B should remain: %v", err)
	}
}

func TestListValidation(t *testing.T) {
	svc, _ := newService()
	mustCreate(t, svc, "A")

	cases := []struct {
		name     string
		status   string
		page     int
		pageSize int
		sort     string
		wantErr  error
	}{
		{"invalid status", "foo", 1, 20, "", apperr.ErrBadRequest},
		{"page zero", "all", 0, 20, "", apperr.ErrInvalidPagination},
		{"pageSize zero", "all", 1, 0, "", apperr.ErrInvalidPagination},
		{"pageSize too large", "all", 1, 200, "", apperr.ErrInvalidPagination},
		{"invalid sort field", "all", 1, 20, "foo:desc", apperr.ErrInvalidSort},
		{"invalid sort direction", "all", 1, 20, "createdAt:up", apperr.ErrInvalidSort},
		{"missing colon", "all", 1, 20, "createdAt", apperr.ErrInvalidSort},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := svc.List(tc.status, "", tc.sort, tc.page, tc.pageSize)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestListDefaultsAndSort(t *testing.T) {
	svc, repo := newService()
	mustCreate(t, svc, "A")

	// 缺省排序应为 created_at DESC
	todos, meta, err := svc.List("all", "", "", 1, 20)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}
	if meta.Total != 1 || meta.Page != 1 || meta.PageSize != 20 || meta.TotalPages != 1 {
		t.Fatalf("unexpected meta: %+v", meta)
	}
	if repo.lastQuery.Sort != "created_at DESC" {
		t.Fatalf("default sort = %q", repo.lastQuery.Sort)
	}

	// 多字段排序
	if _, _, err := svc.List("all", "", "completed:asc,createdAt:desc", 1, 20); err != nil {
		t.Fatalf("multi sort: %v", err)
	}
	if repo.lastQuery.Sort != "completed ASC, created_at DESC" {
		t.Fatalf("multi sort clause = %q", repo.lastQuery.Sort)
	}
}
