// 外部测试包：集成测试经 router 组装完整路由（handler → router 的 import cycle 通过外部包规避）
package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cxh/todolist/server/internal/handler"
	"github.com/cxh/todolist/server/internal/model"
	"github.com/cxh/todolist/server/internal/repository"
	"github.com/cxh/todolist/server/internal/router"
	"github.com/cxh/todolist/server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newTestServer 起真实 Gin 路由（含中间件）+ 内存 SQLite（PRD §10.2 集成测试）。
func newTestServer(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
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
	repo := repository.NewTodoRepository(db)
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)
	r := gin.New()
	router.Setup(r, h, []string{"http://localhost:5173"})
	return r
}

type envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Meta    json.RawMessage `json:"meta"`
}

func doRequest(t *testing.T, r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func decodeEnvelope(t *testing.T, w *httptest.ResponseRecorder) envelope {
	t.Helper()
	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v; body=%s", err, w.Body.String())
	}
	return env
}

func createTodo(t *testing.T, r *gin.Engine, title string) uint {
	t.Helper()
	w := doRequest(t, r, http.MethodPost, "/api/v1/todos", fmt.Sprintf(`{"title":%q}`, title))
	if w.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body=%s", w.Code, w.Body.String())
	}
	env := decodeEnvelope(t, w)
	var todo model.Todo
	if err := json.Unmarshal(env.Data, &todo); err != nil {
		t.Fatalf("decode todo: %v", err)
	}
	return todo.ID
}

// 对应 TC-01
func TestHealthz(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodGet, "/healthz", "")
	if w.Code != http.StatusOK {
		t.Fatalf("healthz status = %d", w.Code)
	}
}

// 对应 TC-02
func TestCreateTodo(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodPost, "/api/v1/todos", `{"title":"学习 Go 语言基础"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); loc != "/api/v1/todos/1" {
		t.Fatalf("location = %q", loc)
	}
	env := decodeEnvelope(t, w)
	if env.Code != 0 {
		t.Fatalf("code = %d", env.Code)
	}
	var todo model.Todo
	if err := json.Unmarshal(env.Data, &todo); err != nil {
		t.Fatalf("decode todo: %v", err)
	}
	if todo.ID == 0 || todo.Title != "学习 Go 语言基础" || todo.Completed {
		t.Fatalf("unexpected todo: %+v", todo)
	}
}

// 对应 TC-03
func TestCreateTodo_EmptyTitle(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodPost, "/api/v1/todos", `{"title":"   "}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40001 {
		t.Fatalf("code = %d, want 40001", env.Code)
	}
}

// 对应 TC-04
func TestCreateTodo_TitleTooLong(t *testing.T) {
	r := newTestServer(t)
	title := strings.Repeat("长", 201)
	w := doRequest(t, r, http.MethodPost, "/api/v1/todos", fmt.Sprintf(`{"title":%q}`, title))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40001 {
		t.Fatalf("code = %d, want 40001", env.Code)
	}
}

func TestCreateTodo_InvalidJSON(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodPost, "/api/v1/todos", `{invalid`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40000 {
		t.Fatalf("code = %d, want 40000", env.Code)
	}
}

// 对应 TC-06（前置 TC-05 创建多条）
func TestListTodos(t *testing.T) {
	r := newTestServer(t)
	for _, title := range []string{"A", "B", "C"} {
		createTodo(t, r, title)
	}
	w := doRequest(t, r, http.MethodGet, "/api/v1/todos", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	env := decodeEnvelope(t, w)
	var todos []model.Todo
	if err := json.Unmarshal(env.Data, &todos); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(todos) != 3 {
		t.Fatalf("expected 3 todos, got %d", len(todos))
	}
	var meta struct {
		Page       int   `json:"page"`
		PageSize   int   `json:"pageSize"`
		Total      int64 `json:"total"`
		TotalPages int   `json:"totalPages"`
	}
	if err := json.Unmarshal(env.Meta, &meta); err != nil {
		t.Fatalf("decode meta: %v", err)
	}
	if meta.Total != 3 || meta.PageSize != 20 {
		t.Fatalf("unexpected meta: %+v", meta)
	}
}

// 对应 TC-07
func TestListTodos_StatusActive(t *testing.T) {
	r := newTestServer(t)
	id := createTodo(t, r, "A")
	if w := doRequest(t, r, http.MethodPatch, fmt.Sprintf("/api/v1/todos/%d", id), `{"completed":true}`); w.Code != http.StatusOK {
		t.Fatalf("patch status = %d", w.Code)
	}
	createTodo(t, r, "B")

	w := doRequest(t, r, http.MethodGet, "/api/v1/todos?status=active", "")
	env := decodeEnvelope(t, w)
	var todos []model.Todo
	_ = json.Unmarshal(env.Data, &todos)
	if len(todos) != 1 || todos[0].Title != "B" {
		t.Fatalf("active list: %+v", todos)
	}
}

// 对应 TC-08
func TestListTodos_StatusCompleted(t *testing.T) {
	r := newTestServer(t)
	id := createTodo(t, r, "A")
	createTodo(t, r, "B")
	if w := doRequest(t, r, http.MethodPatch, fmt.Sprintf("/api/v1/todos/%d", id), `{"completed":true}`); w.Code != http.StatusOK {
		t.Fatalf("patch status = %d", w.Code)
	}

	w := doRequest(t, r, http.MethodGet, "/api/v1/todos?status=completed", "")
	env := decodeEnvelope(t, w)
	var todos []model.Todo
	_ = json.Unmarshal(env.Data, &todos)
	if len(todos) != 1 || todos[0].ID != id {
		t.Fatalf("completed list: %+v", todos)
	}
}

// 对应 TC-09
func TestListTodos_Search(t *testing.T) {
	r := newTestServer(t)
	createTodo(t, r, "学习 Go 语言")
	createTodo(t, r, "给产品写周报")

	w := doRequest(t, r, http.MethodGet, "/api/v1/todos?q=Go", "")
	env := decodeEnvelope(t, w)
	var todos []model.Todo
	_ = json.Unmarshal(env.Data, &todos)
	if len(todos) != 1 || todos[0].Title != "学习 Go 语言" {
		t.Fatalf("search result: %+v", todos)
	}
}

// 对应 TC-10
func TestListTodos_Pagination(t *testing.T) {
	r := newTestServer(t)
	for i := 0; i < 5; i++ {
		createTodo(t, r, fmt.Sprintf("任务%d", i))
	}
	w := doRequest(t, r, http.MethodGet, "/api/v1/todos?page=1&pageSize=2", "")
	env := decodeEnvelope(t, w)
	var todos []model.Todo
	_ = json.Unmarshal(env.Data, &todos)
	if len(todos) != 2 {
		t.Fatalf("expected 2 items, got %d", len(todos))
	}
	var meta struct {
		Page     int   `json:"page"`
		PageSize int   `json:"pageSize"`
		Total    int64 `json:"total"`
	}
	_ = json.Unmarshal(env.Meta, &meta)
	if meta.Page != 1 || meta.PageSize != 2 || meta.Total != 5 {
		t.Fatalf("meta: %+v", meta)
	}
}

// 对应 TC-11
func TestListTodos_InvalidPage(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodGet, "/api/v1/todos?page=0", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40002 {
		t.Fatalf("code = %d, want 40002", env.Code)
	}
}

// 对应 TC-12
func TestListTodos_InvalidSort(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodGet, "/api/v1/todos?sort=foo:desc", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40003 {
		t.Fatalf("code = %d, want 40003", env.Code)
	}
}

// 对应 TC-13
func TestGetTodo(t *testing.T) {
	r := newTestServer(t)
	id := createTodo(t, r, "单条查询")
	w := doRequest(t, r, http.MethodGet, fmt.Sprintf("/api/v1/todos/%d", id), "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	env := decodeEnvelope(t, w)
	var todo model.Todo
	if err := json.Unmarshal(env.Data, &todo); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if todo.ID != id || todo.Title != "单条查询" {
		t.Fatalf("unexpected: %+v", todo)
	}
}

// 对应 TC-14
func TestGetTodo_NotFound(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodGet, "/api/v1/todos/999999", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40400 {
		t.Fatalf("code = %d, want 40400", env.Code)
	}
}

// 对应 TC-15
func TestGetTodo_InvalidID(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodGet, "/api/v1/todos/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40004 {
		t.Fatalf("code = %d, want 40004", env.Code)
	}
}

// 对应 TC-16：完成待办（M1 提前覆盖，M4 补 test.rest 用例）
func TestPatchTodo_Complete(t *testing.T) {
	r := newTestServer(t)
	id := createTodo(t, r, "完成我")
	w := doRequest(t, r, http.MethodPatch, fmt.Sprintf("/api/v1/todos/%d", id), `{"completed":true}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
	env := decodeEnvelope(t, w)
	var todo model.Todo
	_ = json.Unmarshal(env.Data, &todo)
	if !todo.Completed || todo.CompletedAt == nil {
		t.Fatalf("expected completed with timestamp: %+v", todo)
	}
}

// 对应 TC-17：取消完成
func TestPatchTodo_Uncomplete(t *testing.T) {
	r := newTestServer(t)
	id := createTodo(t, r, "先完成")
	if w := doRequest(t, r, http.MethodPatch, fmt.Sprintf("/api/v1/todos/%d", id), `{"completed":true}`); w.Code != http.StatusOK {
		t.Fatalf("complete: %d", w.Code)
	}
	w := doRequest(t, r, http.MethodPatch, fmt.Sprintf("/api/v1/todos/%d", id), `{"completed":false}`)
	env := decodeEnvelope(t, w)
	var todo model.Todo
	_ = json.Unmarshal(env.Data, &todo)
	if todo.Completed || todo.CompletedAt != nil {
		t.Fatalf("expected active with nil completed_at: %+v", todo)
	}
}

// 对应 TC-18：修改标题
func TestPatchTodo_Title(t *testing.T) {
	r := newTestServer(t)
	id := createTodo(t, r, "旧标题")
	w := doRequest(t, r, http.MethodPatch, fmt.Sprintf("/api/v1/todos/%d", id), `{"title":"学习 Go 与 GORM"}`)
	env := decodeEnvelope(t, w)
	var todo model.Todo
	_ = json.Unmarshal(env.Data, &todo)
	if todo.Title != "学习 Go 与 GORM" {
		t.Fatalf("title = %q", todo.Title)
	}
}

// 对应 TC-19：更新不存在
func TestPatchTodo_NotFound(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodPatch, "/api/v1/todos/999999", `{"title":"x"}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40400 {
		t.Fatalf("code = %d", env.Code)
	}
}

func TestPatchTodo_EmptyBody(t *testing.T) {
	r := newTestServer(t)
	id := createTodo(t, r, "A")
	w := doRequest(t, r, http.MethodPatch, fmt.Sprintf("/api/v1/todos/%d", id), `{}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40000 {
		t.Fatalf("code = %d, want 40000", env.Code)
	}
}

// 对应 TC-20：删除单条
func TestDeleteTodo(t *testing.T) {
	r := newTestServer(t)
	id := createTodo(t, r, "删除我")
	w := doRequest(t, r, http.MethodDelete, fmt.Sprintf("/api/v1/todos/%d", id), "")
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %s", w.Body.String())
	}
}

// 对应 TC-21：删除不存在
func TestDeleteTodo_NotFound(t *testing.T) {
	r := newTestServer(t)
	w := doRequest(t, r, http.MethodDelete, "/api/v1/todos/999999", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
	if env := decodeEnvelope(t, w); env.Code != 40400 {
		t.Fatalf("code = %d", env.Code)
	}
}

// 对应 TC-22：清除已完成（幂等）
func TestDeleteCompleted(t *testing.T) {
	r := newTestServer(t)
	idA := createTodo(t, r, "A")
	createTodo(t, r, "B")
	if w := doRequest(t, r, http.MethodPatch, fmt.Sprintf("/api/v1/todos/%d", idA), `{"completed":true}`); w.Code != http.StatusOK {
		t.Fatalf("complete: %d", w.Code)
	}

	w := doRequest(t, r, http.MethodDelete, "/api/v1/todos/completed", "")
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
	listW := doRequest(t, r, http.MethodGet, "/api/v1/todos", "")
	env := decodeEnvelope(t, listW)
	var todos []model.Todo
	_ = json.Unmarshal(env.Data, &todos)
	if len(todos) != 1 || todos[0].Completed {
		t.Fatalf("after clear: %+v", todos)
	}

	// 幂等：再次清除仍返回 204
	w2 := doRequest(t, r, http.MethodDelete, "/api/v1/todos/completed", "")
	if w2.Code != http.StatusNoContent {
		t.Fatalf("idempotent status = %d", w2.Code)
	}
}

// 对应 TC-23：静态路由 /todos/completed 与参数路由 /todos/:id 不冲突
func TestRoutePrecedence(t *testing.T) {
	r := newTestServer(t)
	createTodo(t, r, "A")

	// DELETE /todos/completed 走静态路由（而非 :id）
	if w := doRequest(t, r, http.MethodDelete, "/api/v1/todos/completed", ""); w.Code != http.StatusNoContent {
		t.Fatalf("static route status = %d", w.Code)
	}
	// GET /todos/completed 落入 :id，id 非法 → 40004（证明参数路由未遮蔽静态路由的删除语义）
	if w := doRequest(t, r, http.MethodGet, "/api/v1/todos/completed", ""); w.Code != http.StatusBadRequest {
		t.Fatalf("param route status = %d", w.Code)
	}
}

func TestCORSAllowedOrigin(t *testing.T) {
	r := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("ACAO = %q", got)
	}

	// 非白名单来源不放行
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	req2.Header.Set("Origin", "http://evil.example.com")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if got := w2.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("unexpected ACAO = %q", got)
	}
}
