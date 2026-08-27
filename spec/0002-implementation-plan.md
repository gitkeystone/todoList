# Todo List — 实施计划（Implementation Plan）

| 项目 | 内容 |
| --- | --- |
| 文档编号 | 0002 |
| 文档名称 | Todo List 待办事项应用实施计划 |
| 版本 | v0.1.0（初稿，待评审） |
| 状态 | Draft / 评审中 |
| 上游文档 | `spec/0001-prd-spec.md`（需求与设计文档，简称 PRD） |
| 适用范围 | 研发排期、任务拆分、开发执行、测试验收 |
| 预计总工时 | 约 7.5 人日（M0~M4 五个阶段） |

> **一句话说明**：本文档把 PRD §11 的五个里程碑（M0~M4）展开为**可执行的任务级实施计划**——每个任务给出编号、前置依赖、涉及文件、关键实现要点、验收标准与预估工时；同时约定开发规范、质量门禁与测试推进顺序，保证任何人接手都能按序落地。

---

## 修订记录

| 版本 | 日期 | 修订人 | 说明 |
| --- | --- | --- | --- |
| v0.1.0 | 2026-08-27 | 陈晓会 | 初稿：M0~M4 任务级 WBS、开发约定、质量门禁、测试计划 |

---

## 目录

1. [计划目标与执行原则](#1-计划目标与执行原则)
2. [环境与前置条件](#2-环境与前置条件)
3. [总体路线图](#3-总体路线图)
4. [开发约定](#4-开发约定)
5. [阶段详解（WBS）](#5-阶段详解wbs)
   - 5.1 [M0 项目脚手架（1d）](#51-m0-项目脚手架1d)
   - 5.2 [M1 后端 API（2d）](#52-m1-后端-api2d)
   - 5.3 [M2 前端基础（2d）](#53-m2-前端基础2d)
   - 5.4 [M3 打磨（1.5d）](#54-m3-打磨15d)
   - 5.5 [M4 测试与收尾（1d）](#55-m4-测试与收尾1d)
6. [质量门禁（Quality Gate）](#6-质量门禁quality-gate)
7. [测试实施计划](#7-测试实施计划)
8. [风险与缓解（实施视角）](#8-风险与缓解实施视角)
9. [交付物清单](#9-交付物清单)
10. [附录](#10-附录)

---

## 1. 计划目标与执行原则

### 1.1 计划目标

- 在 **7.5 人日**内交付一个可运行、可测试、可维护的 Todo List 全栈应用；
- 每个阶段结束都有**可演示、可验收**的增量产物（而不是最后一次性交付）；
- 全程落实 PRD 的**契约**：API 信封、错误码、字段命名、设计令牌，作为前后端与测试的共同基线。

### 1.2 执行原则

| 原则 | 说明 |
| --- | --- |
| P-1 契约先行 | 先定接口（PRD §6）与测试用例（PRD §10.3），再写实现；前后端都对着同一份契约开发 |
| P-2 测试伴随 | 后端"实现一个模块即补一个测试"；禁止攒到最后统一补测试 |
| P-3 小步提交 | 每个任务一次可编译、可运行的提交；禁止跨阶段的大杂烩提交 |
| P-4 命令即验收 | 一切以 Makefile 目标为准：能 `make dev` 跑起来、`make test` 全绿才算完成 |
| P-5 文档同步 | 任何影响契约/流程的变更，需同步更新 PRD（0001）与本计划（0002） |
| P-6 单用户假设 | 不引入用户系统、不做并发一致性过度设计，聚焦 PRD §1.4 范围内需求 |

---

## 2. 环境与前置条件

### 2.1 工具链清单（版本下限）

| 工具 | 版本要求 | 用途 | 校验命令 |
| --- | --- | --- | --- |
| Go | ≥ 1.22 | 后端编译 | `go version` |
| Node.js | ≥ 20 | 前端工具链 | `node -v` |
| pnpm | ≥ 9 | 前端包管理 | `pnpm -v` |
| GNU Make | ≥ 4.0 | 生命周期编排 | `make --version` |
| git | ≥ 2.40 | 版本管理 | `git --version` |
| VS Code | 最新稳定版 | 开发 IDE | — |
| REST Client 插件 | 最新 | `.rest` 接口测试 | 扩展商店安装 humao.rest-client |
| golangci-lint | ≥ 1.60 | Go 静态检查 | `golangci-lint version` |
| air | ≥ 1.51（可选） | Go 热重载（`make dev-api`） | `air -v` |
| httpyac | 最新（可选） | CI 中运行 `.rest` 文件 | `httpyac --version` |

### 2.2 前置检查清单

- [ ] 以上工具均安装且版本满足要求；
- [ ] 本机可访问 Go Proxy 与 npm registry（或已配置镜像）；
- [ ] 工作目录 `todolist/` 已初始化 git 仓库；
- [ ] 已通读 PRD（`spec/0001-prd-spec.md`）第 4、6、7、8、9、10 章。

---

## 3. 总体路线图

### 3.1 阶段总览

| 阶段 | 名称 | 工时 | 核心产物 | 阶段出口验收（Gate） |
| --- | --- | --- | --- | --- |
| M0 | 项目脚手架 | 1d | 目录骨架、Makefile、双端工程初始化 | `make setup && make dev` 跑通空页面 + `/healthz` 200 |
| M1 | 后端 API | 2d | 6 个端点、分层代码、单测/集成测试、`test.rest`（TC-01~15） | `make test` 全绿；`test.rest` TC-01~TC-15 通过 |
| M2 | 前端基础 | 2d | 布局组件、Pinia store、CRUD 联调 | 页面可完成 创建/完成/删除/清除，刷新持久 |
| M3 | 打磨 | 1.5d | 动效、空状态、深色模式、响应式、无障碍、toast | 视觉/交互符合 PRD §4；375px 可操作 |
| M4 | 测试与收尾 | 1d | 前端 Vitest、`test.rest` 补全（TC-16~23）、README | 全部质量门禁通过，PRD §13 DoD 达成 |

### 3.2 阶段依赖与关键路径

```mermaid
flowchart LR
    M0[M0 脚手架] --> M1[M1 后端 API]
    M1 --> M2[M2 前端基础]
    M2 --> M3[M3 打磨]
    M3 --> M4[M4 测试与收尾]
    M1 -. 契约（test.rest/信封/错误码）先行 .-> M2
```

- **关键路径**：M0 → M1 → M2 → M3 → M4，串行推进；
- **并行优化**：M1 的后端实现与 `test.rest` 编写、M3 的动效与无障碍可并行小步推进；
- **里程碑节奏**：建议每完成一个阶段跑一次完整质量门禁（见 §6），不合格不进入下一阶段。

---

## 4. 开发约定

### 4.1 分支与提交

- 分支模型：`main`（始终可运行）+ 每个里程碑一个 `feat/Mx-xxx` 分支，合并前过质量门禁；
- 提交信息：`<type>(<scope>): <subject>`，type ∈ `feat|fix|test|docs|chore|refactor`，scope 如 `server|web|makefile|spec`；
- 单次提交粒度：一个任务（如 `feat(server): 新增 Todo repository 层`）。

### 4.2 代码规范

| 端 | 规范 |
| --- | --- |
| Go | `gofmt`/`goimports` 格式化；`golangci-lint`（含 govet、errcheck、staticcheck、ineffassign）无 error；导出类型/函数写注释；错误统一包装返回，不吞错 |
| 前端 | ESLint（vue 官方推荐）+ Prettier；TS 严格模式（`strict: true`）；组件 props/emits 显式声明；命名：组件 PascalCase、文件 kebab-case、store 用 `useXxxStore` |
| 通用 | 提交前跑 `make fmt && make lint` |

### 4.3 契约基线（实现时须严格遵守，源自 PRD §6）

| 项 | 约定 |
| --- | --- |
| API 前缀 | `/api/v1`，健康检查 `/healthz` |
| 响应信封 | `{ "code": 0, "message": "ok", "data": …, "meta": … }` |
| 错误码 | `40000/40001/40002/40003/40004/40400/50000/50001`（见 PRD §6.2） |
| 字段命名 | JSON 一律 camelCase：`createdAt/updatedAt/completedAt` |
| 日期 | RFC3339（UTC 存储，前端本地化展示） |
| 分页 | `page`(≥1) / `pageSize`(1~100)，`meta` 返回 `{page,pageSize,total,totalPages}` |
| 排序 | `sort=field:direction`，field ∈ `createdAt/updatedAt/completedAt` |

### 4.4 端口与代理约定

- 后端 `:8080`，前端 `:5173`；
- 开发时前端经 Vite 代理 `/api` → `http://localhost:8080`，避免 CORS；
- 生产模式由后端托管 `web/dist`（`make build && make run`）。

---

## 5. 阶段详解（WBS）

> 任务编号规则：`T{M}-{NN}`（如 T1-03 表示 M1 阶段第 3 个任务）。每个任务包含：目的、涉及文件、关键要点、验收、前置、工时。工时单位为"人日"，0.25d ≈ 2 小时。

---

### 5.1 M0 项目脚手架（1d）

**目标**：建立可编译、可运行、可扩展的工程骨架，为 M1/M2 提供稳定的地基。

| 编号 | 任务 | 工时 | 前置 |
| --- | --- | --- | --- |
| T0-01 | 初始化仓库与目录骨架 | 0.25d | — |
| T0-02 | 根 Makefile（全目标） | 0.25d | T0-01 |
| T0-03 | Go 后端最小可运行（`/healthz`） | 0.25d | T0-01 |
| T0-04 | Vue 前端工程初始化（Vite+TS+Tailwind+shadcn-vue） | 0.25d | T0-01 |
| T0-05 | 联调冒烟与门禁检查 | 0.25d | T0-02~T0-04 |

**阶段出口**：`make setup && make dev` 一键启动；浏览器 :5173 显示空页面；`curl http://localhost:8080/healthz` 返回 200。

#### T0-01 初始化仓库与目录骨架

- **目的**：落地 PRD §5.3 的目录结构，建立版本管理与配置基线。
- **涉及文件**：`.gitignore`、`.env.example`、`README.md`、`spec/0001-prd-spec.md`（已有）。
- **关键要点**：
  - `.gitignore` 至少包含：`bin/`、`web/dist/`、`data/`、`*.db`、`*.db-wal`、`*.db-shm`、`node_modules/`、`.env`、IDE 目录；
  - `.env.example` 提供 PRD §9.3 全部变量样例：`PORT=8080`、`DB_PATH=data/todolist.db`、`GIN_MODE=debug`、`ALLOWED_ORIGINS=http://localhost:5173`、`WEB_DIST=web/dist`；
  - 按 PRD §5.3 创建 `server/` 与 `web/` 骨架目录（空目录放 `.gitkeep`）。
- **验收**：`git status` 干净；`.gitignore` 生效（`data/` 不被跟踪）。

#### T0-02 根 Makefile（全目标）

- **目的**：一次性落地 PRD §9 的全部生命周期目标，后续阶段只填实现、不增删目标名。
- **涉及文件**：`Makefile`。
- **关键要点**：
  - 实现 `help/setup/dev/dev-api/dev-web/build/run/test/test-web/test-api/lint/fmt/seed/db-reset/clean` 15 个目标（骨架可先用占位命令，M1/M2 逐步充实）；
  - `help` 用 `awk` 解析 `##` 注释（照抄 PRD §9.2 片段）；
  - `setup`、`clean`、`db-reset` 必须**幂等可重入**；
  - `dev` 用 `trap 'kill 0' INT TERM EXIT` 保证 Ctrl+C 双进程同退。
- **验收**：`make help` 输出完整目标表；`make setup && make setup` 第二次不报错；`make clean` 清理 `bin/ web/dist data/*.db`。

#### T0-03 Go 后端最小可运行

- **目的**：验证 Go 工具链与 Gin 集成，提供 `/healthz` 作为联调探针。
- **涉及文件**：`server/go.mod`、`server/cmd/api/main.go`、`server/internal/config/config.go`（最小版）。
- **关键要点**：
  - `go mod init github.com/cxh/todolist/server`（模块路径可按实际 Git 仓库调整）；
  - 依赖：`gin-gonic/gin`、`gorm.io/gorm`、`github.com/glebarez/sqlite`（纯 Go，规避 CGO，PRD §12 R-1）；
  - `main.go`：加载 `PORT`（默认 8080）→ 创建 gin 引擎 → 注册 `GET /healthz` 返回 `{"status":"ok"}` → `r.Run(":"+port)`；
  - `config` 用标准库 `os.Getenv` + 默认值，暂不引第三方配置库。
- **验收**：`cd server && go build ./...` 通过；`go run ./cmd/api` 后 `curl :8080/healthz` 返回 200。

#### T0-04 Vue 前端工程初始化

- **目的**：搭好 Vite + TS + Tailwind + shadcn-vue 的前端地基（PRD §8.1）。
- **涉及文件**：`web/`（package.json、vite.config.ts、tsconfig.json、tailwind.config.ts、index.html、src/main.ts、App.vue）。
- **关键要点**：
  - `pnpm create vite web --template vue-ts`，然后 `pnpm add -D tailwindcss@^3 postcss autoprefixer` 并 `npx tailwindcss init -p`；
  - `pnpm add vue-router?`——**不需要**：单页应用，无路由（PRD §4.3 单页布局）；
  - `pnpm add pinia axios`；`pnpm add -D vitest @vue/test-utils jsdom`；
  - shadcn-vue：`pnpm dlx shadcn-vue@latest init`（组件输出到 `src/components`，主题色映射预留）；
  - `vite.config.ts` 配好 `/api`、`/healthz` 代理与 `@` 别名（PRD §8.1 片段）；
  - `index.html` 预留防 FOUC 主题脚本占位（M3 填充）。
- **验收**：`cd web && pnpm dev` 后 :5173 打开显示空白页无报错；`pnpm build` 通过；`pnpm vitest run` 可跑通一个示例测试。

#### T0-05 联调冒烟与门禁检查

- **目的**：验证 M0 全部产物能协同工作。
- **步骤**：`make setup` → `make dev` → 浏览器 :5173 → `curl :8080/healthz` → `make build` → `make clean`。
- **验收**：对照 §3.1 阶段出口全部满足；`make lint`（如已配置）无致命告警。

---

### 5.2 M1 后端 API（2d）

**目标**：按 PRD §6、§7 交付完整后端——6 个端点、分层结构、中间件、数据库迁移、单元与集成测试、`test.rest` 前半段（TC-01~TC-15）。

**实现顺序建议（自底向上）**：config → model → repository → service → response/错误码 → handler → middleware → router → main → 测试 → test.rest → seed。

| 编号 | 任务 | 工时 | 前置 |
| --- | --- | --- | --- |
| T1-01 | config 模块完善 | 0.1d | T0-03 |
| T1-02 | Todo 模型 + 自动迁移 | 0.15d | T1-01 |
| T1-03 | repository 层（CRUD + 筛选分页排序） | 0.35d | T1-02 |
| T1-04 | service 层（校验与业务规则） | 0.2d | T1-03 |
| T1-05 | 统一响应信封与错误码 | 0.15d | T1-02 |
| T1-06 | handler 层（6 端点） | 0.3d | T1-04, T1-05 |
| T1-07 | 中间件（CORS/Logger/Recovery） | 0.15d | T1-05 |
| T1-08 | 路由注册（含静态路由优先级） | 0.1d | T1-06, T1-07 |
| T1-09 | main.go 组装与优雅退出 | 0.1d | T1-01~T1-08 |
| T1-10 | 单元 + 集成测试 | 0.3d | T1-08 |
| T1-11 | test.rest（TC-01~TC-15） | 0.2d | T1-09 |
| T1-12 | seed 演示数据命令 | 0.1d | T1-09 |

**阶段出口**：`make test` 全绿（覆盖率达标）；`make dev-api` 后 VS Code 执行 `test.rest` TC-01~TC-15 全通过；`make seed` 可注入演示数据。

#### T1-01 config 模块完善

- **涉及文件**：`server/internal/config/config.go`。
- **关键要点**：定义 `Config{Port, DBPath, GinMode, AllowedOrigins []string, WebDist}`；`Load()` 读环境变量并应用默认值（PRD §9.3）；提供 `Addrs()`/`DSN()` 便捷方法。
- **验收**：带环境变量启动时配置生效；缺失时回退默认值。

#### T1-02 Todo 模型 + 自动迁移

- **涉及文件**：`server/internal/model/todo.go`。
- **关键要点**：照 PRD §7.2 定义 `Todo`（含 GORM tag 与 json tag，camelCase）；`CompletedAt *time.Time` 可空；`Init(db)` 内 `db.AutoMigrate(&Todo{})`。
- **验收**：启动后 `data/todolist.db` 生成 `todos` 表，字段与 PRD §7.2 DDL 一致；`completed`、`created_at` 索引存在。

#### T1-03 repository 层

- **涉及文件**：`server/internal/repository/todo.go`（接口 + 实现）。
- **关键要点**：
  - 接口 `TodoRepository`：`Create/GetByID/Update/Delete/DeleteCompleted/List(ListQuery)/Count(ListQuery)`；
  - `ListQuery` 结构：`Status string / Q string / Page, PageSize int / Sort string`；
  - 实现要点：
    - 筛选：`status=active` → `completed = ?`（false）；`completed` → `completed = ?`（true）；
    - 搜索：`q != ""` 时 `WHERE title LIKE ?`（`%q%`），注意转义 `%`/`_`；
    - 排序：白名单映射（`createdAt→created_at` 等）+ 方向校验，拼接 `ORDER BY` 防注入；
    - 分页：`LIMIT ? OFFSET ?`，`Offset=(page-1)*pageSize`；
    - 连接池：`sqlDB.SetMaxOpenConns(1)`（SQLite 单写者，PRD §7.1）。
- **验收**：`go vet` 通过；各方法行为与 PRD §6.3 一致。

#### T1-04 service 层

- **涉及文件**：`server/internal/service/todo.go`。
- **关键要点**：
  - 标题校验：`strings.TrimSpace` 后 `1 ≤ len(rune) ≤ 200`，否则返回业务错误 `ErrInvalidTitle`（映射 40001）；
  - PATCH 语义：仅处理请求中出现的字段；`completed=true` 且原值为 false 时写 `CompletedAt=now`；`completed=false` 时置 `CompletedAt=nil`；
  - 将 repository 错误映射为业务错误（`ErrNotFound` 40400 / `ErrDB` 50001）。
- **验收**：边界用例（空标题、201 字符、重复 completed=true 不刷新 completedAt）单测覆盖。

#### T1-05 统一响应信封与错误码

- **涉及文件**：`server/internal/handler/response.go`（或 `pkg/response`）。
- **关键要点**：
  - 定义 `OK(c, data, meta)` / `Created(c, data, location)` / `NoContent(c)` / `Fail(c, httpStatus, code, msg)`；
  - 错误码常量 + `AppError{Code, HTTPStatus, Message}` 类型，`Error()` 实现 error 接口；
  - 与 PRD §6.2 错误码表一一对应。
- **验收**：所有响应结构符合信封约定；错误响应 `data` 为 `null`。

#### T1-06 handler 层（6 端点）

- **涉及文件**：`server/internal/handler/todo.go`。
- **关键要点**：
  - `List`：绑定 query → 参数校验（page/pageSize/sort 非法返回 40002/40003）→ service/repository → 信封 + meta；
  - `Create`：`ShouldBindJSON` → 校验 → 201 + `Location` 头；
  - `Get`：id 解析失败返回 40004，不存在返回 40400；
  - `Patch`：先查后改（或 Update 返回影响行数），保证 404 语义；
  - `Delete`：影响行数 0 → 40400；成功 204；
  - `DeleteCompleted`：幂等 204。
- **验收**：6 个端点行为与 PRD §6.3 逐条一致。

#### T1-07 中间件

- **涉及文件**：`server/internal/middleware/{cors,logger,recovery}.go`。
- **关键要点**：CORS 读 `Config.AllowedOrigins`（含预检 OPTIONS）；Logger 用 JSON 行格式输出方法/路径/状态/耗时；Recovery 捕获 panic 返回 `{code:50000}` 信封并记录堆栈。
- **验收**：跨域预检正常；panic 不泄露堆栈给客户端。

#### T1-08 路由注册

- **涉及文件**：`server/internal/router/router.go`。
- **关键要点**：**先注册 `DELETE /todos/completed`，再注册 `/todos/:id` 参数路由**（PRD §6.1 风险说明、TC-23）；组装 `Group("/api/v1")`；`/healthz` 挂在根路径。
- **验收**：`DELETE /todos/completed` 与 `GET /todos/:id` 互不冲突。

#### T1-09 main.go 组装

- **涉及文件**：`server/cmd/api/main.go`。
- **关键要点**：Load config → 打开 SQLite（DSN 含 PRAGMA：WAL/busy_timeout/foreign_keys/synchronous，PRD §7.1）→ 迁移 → 注入依赖 → 注册路由 → 启动；监听 `os.Interrupt`/`SIGTERM` 优雅退出（`http.Server.Shutdown` 带超时）。
- **验收**：Ctrl+C 可正常退出且无数据损坏；启动日志清晰。

#### T1-10 单元 + 集成测试

- **涉及文件**：`server/internal/**/*_test.go`。
- **关键要点**：
  - repository 测试：`file::memory:?cache=shared` 内存 SQLite，覆盖 CRUD/筛选/搜索/分页/排序/删除已完成；
  - service 测试：标题校验、completedAt 置空规则（PRD §10.2）；
  - handler 集成测试：`httptest` + 真实路由（含中间件），逐端点断言状态码/信封/错误码（对齐 TC-01~TC-15）；
  - 每个测试独立内存库或事务回滚，互不污染；
  - 覆盖率目标：service ≥ 85%，handler ≥ 80%。
- **验收**：`make test` 全绿；`go test ./... -cover` 覆盖率达标。

#### T1-11 test.rest（TC-01~TC-15）

- **涉及文件**：`server/test.rest`。
- **关键要点**：照 PRD §10.3 用例清单实现 TC-01~TC-15；文件头定义 `@baseUrl = http://localhost:8080/api/v1`；TC-02 用 `# @name createTodo` 命名，TC-13 链式引用 `{{createTodo.response.body.$.data.id}}`。
- **验收**：VS Code 逐条执行，状态码与 PRD 表格期望一致。

#### T1-12 seed 命令

- **涉及文件**：`server/cmd/seed/main.go`。
- **关键要点**：复用 config/model/repository，插入 12 条演示待办（含中英文、部分已完成，PRD §7.4）；重复执行需幂等（可先清空再插入）。
- **验收**：`make seed` 后 `GET /todos` 返回 12 条。

---

### 5.3 M2 前端基础（2d）

**目标**：按 PRD §8 完成前端基础能力——工程配置、设计令牌、API 客户端、Pinia store、核心布局与列表组件，与 M1 后端完成 CRUD 联调。

| 编号 | 任务 | 工时 | 前置 |
| --- | --- | --- | --- |
| T2-01 | Vite/TS 工程配置细化 | 0.15d | T0-04 |
| T2-02 | 设计令牌与全局样式 | 0.3d | T0-04 |
| T2-03 | shadcn-vue 基础组件接入 | 0.15d | T0-04 |
| T2-04 | 类型定义与 API 客户端 | 0.2d | T2-01 |
| T2-05 | Pinia todo store | 0.3d | T2-04 |
| T2-06 | 布局组件（Header/Hero/Input） | 0.35d | T2-02, T2-05 |
| T2-07 | 列表与单项组件（List/Item） | 0.35d | T2-05 |
| T2-08 | 分段控件与页脚（Segmented/Footer） | 0.2d | T2-05 |
| T2-09 | 前后端联调冒烟 | 0.25d | T2-06~T2-08 |

**阶段出口**：`make dev` 下，页面可完成 创建 → 完成/取消 → 删除 → 清除已完成，刷新后数据持久；统计数字实时正确。

#### T2-01 Vite/TS 工程配置细化

- **涉及文件**：`web/vite.config.ts`、`web/tsconfig.json`、`web/package.json`。
- **关键要点**：`@` 别名、`/api`+`/healthz` 代理、`manualChunks`（vue/pinia 分包，PRD §8.1）；tsconfig 开 `strict`；scripts 补 `lint`、`format`（eslint + prettier）。
- **验收**：`pnpm build` 无 TS 错误；产物分包生效。

#### T2-02 设计令牌与全局样式

- **涉及文件**：`web/src/assets/design-tokens.css`、`web/tailwind.config.ts`、`web/index.html`。
- **关键要点**：按 PRD §4.2.1 定义浅/深两套 CSS 变量（`--bg/--bg-elevated/--text-primary/--text-secondary/--accent/--success/--danger/--border/--frosted`）；Tailwind 通过 `hsl(var(--token))` 映射颜色令牌，`darkMode: 'class'`；扩展圆角/阴影/字体栈/关键帧（PRD §8.1 片段）；`index.html` 预置防 FOUC 内联脚本（读 `localStorage['todolist.theme']`）。
- **验收**：`<html class="dark">` 时整页令牌切换生效；无闪烁。

#### T2-03 shadcn-vue 基础组件接入

- **涉及文件**：`web/components.json`、`web/src/components/*`。
- **关键要点**：接入 `button`、`input`、`skeleton`、`sonner`（toast）、`tooltip`；主题色映射到 §T2-02 令牌（accent → `hsl(var(--accent))` 等）。
- **验收**：`pnpm dlx shadcn-vue@latest add` 各组件可正常渲染并跟随主题。

#### T2-04 类型定义与 API 客户端

- **涉及文件**：`web/src/types/todo.ts`、`web/src/api/todo.ts`。
- **关键要点**：定义 `Todo`、`TodoQuery`、`PageResult<T>`、`ApiEnvelope<T>`（与 PRD §6.2 信封一致）；axios 实例 `baseURL:'/api/v1'`、`timeout:10_000`；响应拦截器解包信封，`code!==0` 或 HTTP 非 2xx 时统一 reject 并提取 `message`；导出 `todoApi`（list/create/get/patch/remove/clearCompleted，PRD §8.3）。
- **验收**：TS 编译通过；错误 message 可被上层直接展示。

#### T2-05 Pinia todo store

- **涉及文件**：`web/src/stores/todo.ts`。
- **关键要点**：state（items/status/keyword/meta/loading/error）；getters（activeCount/completedCount/allCount）；actions（fetchList/create/toggle/updateTitle/remove/clearCompleted）；**乐观更新**：toggle/remove 先改本地，失败回滚 + 抛错；写操作后用响应体合并本地项（PRD §8.2）。
- **验收**：store 单测可覆盖乐观更新与回滚（M4 补测，此处先实现）。

#### T2-06 布局组件（Header/Hero/Input）

- **涉及文件**：`web/src/features/todo/TodoInput.vue`、`SiteHeader.vue`、`Hero.vue`（或并入 App.vue 的局部组件）。
- **关键要点**：
  - Header：`position: sticky` + `backdrop-blur(20px)` + 半透明背景（`--frosted`）；右侧主题切换按钮（M3 接逻辑）；
  - Hero：大标题（负字距）+ 文案统计（"还剩 N 项 · 已完成 M 项"）；
  - Input：大号胶囊输入框（`rounded-full`），占位"添加一项待办…"；`@keydown.enter` 提交；**监听 `compositionstart/end` 保护中文输入**（PRD §12 R-3）；空输入抖动（M3 动效）。
- **验收**：视觉符合 PRD §4.3 布局图；回车触发 `store.create`。

#### T2-07 列表与单项组件（List/Item）

- **涉及文件**：`web/src/features/todo/TodoList.vue`、`TodoItem.vue`。
- **关键要点**：`TodoItem`：自定义圆形勾选框（SVG，M3 补描边动画）、标题（completed 时 line-through + 变灰）、相对时间（`createdAt` 本地化显示，PRD §8.5 时区约定）、hover 显示删除按钮（移动端常显）；`TodoList`：`<TransitionGroup>` 包裹（M3 补 FLIP），`v-for` 渲染 store.items。
- **验收**：勾选/删除事件正确调用 store；刷新后状态一致。

#### T2-08 分段控件与页脚（Segmented/Footer）

- **涉及文件**：`web/src/features/todo/TodoSegmented.vue`、`TodoFooter.vue`。
- **关键要点**：分段控件三项（全部/进行中/已完成）带数量角标，切换更新 `store.status` 并重新 fetch；页脚"还剩 N 项未完成"（store getter）+ "清除已完成 (N)"（N=0 隐藏，PRD FR-05）。
- **验收**：切换视图数据正确；角标与后端 total 一致。

#### T2-09 前后端联调冒烟

- **步骤**：`make dev` → 浏览器完成完整用户旅程（创建→完成→取消→删除→清除已完成→刷新持久）→ `pnpm build` 后 `make run` 验证生产模式单进程托管。
- **验收**：对照 §3.1 阶段出口全部满足。

---

### 5.4 M3 打磨（1.5d）

**目标**：把"能用"提升为"好看好用"——落实 PRD §4 的 Apple 风格动效、空状态、深色模式、响应式、无障碍与反馈体系。

| 编号 | 任务 | 工时 | 前置 |
| --- | --- | --- | --- |
| T3-01 | 动效体系 | 0.4d | T2-07 |
| T3-02 | 空状态 / 骨架 / 错误态 | 0.2d | T2-06 |
| T3-03 | 深色模式（useTheme） | 0.2d | T2-02 |
| T3-04 | 响应式与无障碍 | 0.3d | T2-06~T2-08 |
| T3-05 | 键盘交互与快捷操作 | 0.2d | T3-01 |
| T3-06 | toast 反馈体系 | 0.2d | T2-03 |

**阶段出口**：视觉/交互符合 PRD §4 规范；375px 完整可操作；`prefers-reduced-motion` 下动效降级；深浅色无闪烁。

#### T3-01 动效体系

- **涉及文件**：`web/tailwind.config.ts`（关键帧）、各业务组件、`web/src/assets/design-tokens.css`。
- **关键要点**（对应 PRD §4.2.4 动效表）：
  - 统一缓动变量：`--ease-apple: cubic-bezier(0.32, 0.72, 0, 1)`；
  - 列表入场：`<TransitionGroup>` + FLIP（`v-move`），位移 8px + 淡入，相邻项 40ms 错峰；
  - 勾选动画：SVG `stroke-dashoffset` 描边绘制 + 圆框填充（250ms）；
  - 删除：先量高度 → `height:0` + `opacity:0`（250ms）→ 移除 DOM；
  - 分段滑块：绝对定位背景 + `transform` 过渡；
  - 空输入抖动：300ms 水平 shake；
  - 全局 `@media (prefers-reduced-motion: reduce)` 关闭位移类动画。
- **验收**：各动效时长/曲线符合规范表；开启"减少动态效果"后无位移动画。

#### T3-02 空状态 / 骨架 / 错误态

- **涉及文件**：`web/src/features/todo/`（EmptyState.vue、Skeleton.vue 复用 shadcn skeleton）。
- **关键要点**：三种空状态（首次使用/筛选无结果/搜索无结果）与加载骨架、错误重试按钮，对应 PRD §4.5。
- **验收**：各状态文案与图标符合 PRD 表格；加载过程无布局跳动。

#### T3-03 深色模式（useTheme）

- **涉及文件**：`web/src/composables/useTheme.ts`、`SiteHeader.vue`。
- **关键要点**：三态（light/dark/system）；`system` 用 `matchMedia('(prefers-color-scheme: dark)')` 解析并监听变化；持久化到 `localStorage['todolist.theme']`；切换时更新 `<html class>`。
- **验收**：三态切换即时生效无闪烁；system 模式跟随系统变化。

#### T3-04 响应式与无障碍

- **涉及文件**：全局样式 + 各组件。
- **关键要点**：375px 断点（输入框图标隐藏、删除按钮常显、字号 -2px）；触屏点击目标 ≥44px；`focus-visible` 焦点环；`aria-pressed`（分段控件）、`aria-live="polite"`（统计数字）、`role="alert"`（toast）；对比度达 WCAG AA（PRD §4.6）。
- **验收**：375px 宽度走查全流程；Tab 键全程可达；axe 扫描无 critical 问题。

#### T3-05 键盘交互与快捷操作

- **涉及文件**：`TodoList.vue`、`TodoInput.vue`。
- **关键要点**：列表项聚焦时 `Space` 切换完成、`Delete`/`Backspace` 删除（PRD §4.4 I-7）；`⌘K`/`Ctrl+K` 聚焦搜索（搜索为 P1，若无则聚焦输入框）；自动聚焦输入框（桌面端）。
- **验收**：纯键盘可完成创建/完成/删除全流程。

#### T3-06 toast 反馈体系

- **涉及文件**：`App.vue`（`<Sonner />` 出口）、`api/todo.ts`（拦截器触发）、store。
- **关键要点**：错误统一 toast（顶部居中，`role="alert"`）；写操作成功轻提示（删除/清除完成）；toast 时长与样式符合 Apple 风格（浅色白底/深色黑底，PRD §4.4 I-9）。
- **验收**：接口失败、删除成功均有明确反馈；不打断操作流。

---

### 5.5 M4 测试与收尾（1d）

**目标**：补齐测试与工程收尾——前端 Vitest、`test.rest` 全量（TC-16~23）、lint/format 全绿、README 完善、对照 PRD §13 DoD 逐项验收。

| 编号 | 任务 | 工时 | 前置 |
| --- | --- | --- | --- |
| T4-01 | 前端 Vitest 测试 | 0.3d | M3 |
| T4-02 | test.rest 补全（TC-16~TC-23） | 0.2d | M1, M3 |
| T4-03 | lint/format 全绿 + CI 脚本 | 0.2d | M3 |
| T4-04 | README 与文档收尾 | 0.15d | T4-03 |
| T4-05 | 手动验收走查（PRD §10.5/§13） | 0.15d | T4-01~T4-04 |

**阶段出口**：全部质量门禁（§6）通过；PRD §13 DoD 清单勾选完毕。

#### T4-01 前端 Vitest 测试

- **涉及文件**：`web/src/__tests__/`。
- **关键要点**：`TodoInput`（回车创建、空输入不触发、composition 保护）、`TodoItem`（勾选/删除事件、completed 样式类）、`useTodoStore`（fetch/create 乐观插入/toggle 乐观更新与失败回滚/remove/clearCompleted）、`useTheme`（三态切换与持久化）；网络层用 `axios-mock-adapter` 模拟 PRD §6 契约（状态码/信封）（PRD §10.4）。
- **验收**：`make test-web` 全绿。

#### T4-02 test.rest 补全（TC-16~TC-23）

- **涉及文件**：`server/test.rest`。
- **关键要点**：补齐完成/取消完成/改标题/更新不存在/删除/删除不存在/清除已完成/路由优先级用例（PRD §10.3 表格）；**全量回归 TC-01~TC-23**。
- **验收**：VS Code 全量通过；TC-23 验证静态/参数路由不冲突。

#### T4-03 lint/format 全绿 + CI 脚本

- **涉及文件**：`Makefile`（可选加 `ci` 目标）、`scripts/ci.sh`。
- **关键要点**：`golangci-lint run ./...`、`pnpm lint`、`pnpm format --check` 零告警；`ci.sh` 按序执行：fmt 检查 → lint → test → test-web → build。
- **验收**：`make lint && make fmt` 干净；`scripts/ci.sh` 一键通过。

#### T4-04 README 与文档收尾

- **涉及文件**：`README.md`、`.env.example`。
- **关键要点**：快速开始（前置条件 → `make setup` → `make dev` → 访问 :5173）、Makefile 目标表、REST Client 测试说明、目录结构、文档索引（0001/0002）。
- **验收**：新机器按 README 可独立跑通。

#### T4-05 手动验收走查

- **执行**：PRD §10.5 手动验收清单 + §13 DoD 全项勾选；发现缺陷走"修复 → 回归 → 关闭"闭环。
- **验收**：DoD 100% 达成，本次交付完成。

---

## 6. 质量门禁（Quality Gate）

每个阶段合入 `main` 前必须通过对应门禁；任一不通过则回到所属阶段修复。

| 门禁 | M0 | M1 | M2 | M3 | M4 |
| --- | :-: | :-: | :-: | :-: | :-: |
| `make help` 输出完整 | ✅ | ✅ | ✅ | ✅ | ✅ |
| `make setup` 幂等可重入 | ✅ | ✅ | ✅ | ✅ | ✅ |
| `make dev` 一键启动双端 | ✅ | ✅ | ✅ | ✅ | ✅ |
| `make test`（后端单测+集成）全绿 | — | ✅ | ✅ | ✅ | ✅ |
| handler 覆盖率 ≥ 80% / service ≥ 85% | — | ✅ | ✅ | ✅ | ✅ |
| `test.rest` 通过用例 | — | TC-01~15 | TC-01~15 | TC-01~15 | TC-01~23 |
| `make build` 成功 | — | ✅ | ✅ | ✅ | ✅ |
| `make run` 生产模式冒烟 | — | — | ✅ | ✅ | ✅ |
| `make lint` / `fmt` 零告警 | ⚠️ | ⚠️ | ⚠️ | ⚠️ | ✅ |
| `make test-web` 全绿 | — | — | — | — | ✅ |
| 手动验收（PRD §10.5） | — | — | ⚠️ | ⚠️ | ✅ |

> ⚠️ = 该阶段起开始要求（M0 起要求 lint 无致命项；M2/M3 手动验收做部分走查，M4 全量）。

---

## 7. 测试实施计划

### 7.1 测试推进顺序（与阶段绑定）

| 阶段 | 新增测试 | 对应 PRD 章节 |
| --- | --- | --- |
| M1 | repository 单测（内存 SQLite）、service 单测、handler 集成测试（httptest）、test.rest TC-01~15 | §10.2 / §10.3 |
| M2 | 联调冒烟（手动），store 行为由 M4 补自动化 | §10.4 |
| M3 | 无障碍走查（axe）、reduced-motion 检查（手动） | §10.5 |
| M4 | 前端 Vitest（组件/store/composable + axios-mock-adapter）、test.rest TC-16~23 全量回归、手动验收 | §10.3 / §10.4 / §10.5 |

### 7.2 契约一致性保障

- 前端 `axios-mock-adapter` 的 mock 数据与 `test.rest` 断言同一份契约（信封结构、错误码、字段名）；
- 任何契约变更走 PRD 评审流程，同步更新：PRD §6.2/§6.3 → `test.rest` → 前端 mock → 后端集成测试断言。

### 7.3 覆盖率目标与度量

- `make test` 输出 `-cover`；M1 出口检查 handler/service 覆盖率（§6 门禁表）；
- 前端以"关键路径用例覆盖"为准（不追求行覆盖数字），关键路径：创建 → 完成 → 删除 → 清除已完成。

---

## 8. 风险与缓解（实施视角）

> 除 PRD §12 已列风险外，补充实施过程中的操作性风险。

| # | 风险 | 影响 | 缓解 |
| --- | --- | --- | --- |
| E-1 | `glebarez/sqlite` 与 GORM 版本兼容问题 | 编译/运行异常 | 初始化时固定依赖版本并 `go mod tidy`；遇到问题先查 GORM 官方 sqlite 支持表 |
| E-2 | pnpm 安装 shadcn-vue 需交互式 init | 自动化受阻 | 采用非交互参数（`--defaults` 或直接手写 components.json），文档给出两种方式 |
| E-3 | M1 联调时后端未起导致前端大量报错 | 误判前端缺陷 | 联调前先跑 `make test-api` 确认后端就绪；前端错误统一走 toast 不阻塞开发 |
| E-4 | 动效导致性能问题（列表项多时） | 卡顿 | 动画只作用于可视列表；>200 项时先降级为无 FLIP（PRD B-04 虚拟滚动为 P2） |
| E-5 | `test.rest` 依赖数据库状态（如 total ≥ 5） | 用例不稳定 | 用例前置用 seed/清理步骤保证可重复；TC 之间依赖通过请求链传递 id，不硬编码 |
| E-6 | 单阶段工时预估偏差 | 排期失准 | 每个任务 ≤ 0.5d，超时即拆分上报；阶段出口按"功能可用"判定，不追求完美 |
| E-7 | 前后端命名不一致（camelCase vs snake_case） | 联调返工 | 以 PRD §6.2 契约为准，后端 json tag 显式声明 camelCase，前端类型定义照抄 |

---

## 9. 交付物清单

| # | 交付物 | 来源阶段 | 验收说明 |
| --- | --- | --- | --- |
| D-1 | `Makefile`（15 个目标） | M0 | `make help` 完整、目标幂等 |
| D-2 | Go 后端（config/model/repository/service/handler/middleware/router） | M1 | `make test` 全绿，覆盖率达标 |
| D-3 | `server/test.rest`（TC-01~TC-23） | M1+M4 | VS Code 全量通过 |
| D-4 | SQLite 数据库文件（`data/todolist.db`，运行期生成） | M1 | 结构符合 PRD §7.2 |
| D-5 | Vue 前端（布局/动效/主题/store/API 层） | M2+M3 | 视觉符合 PRD §4，功能完整 |
| D-6 | 前端测试（Vitest） | M4 | `make test-web` 全绿 |
| D-7 | `README.md`、`.env.example` | M0+M4 | 新环境可按文档独立跑通 |
| D-8 | `spec/0001-prd-spec.md`、`spec/0002-implementation-plan.md` | — | 契约与计划同步一致 |

---

## 10. 附录

### 10.1 常用命令速查

```bash
make help          # 查看全部生命周期目标
make setup         # 安装前后端依赖
make dev           # 一键启动开发环境（:5173 前端 / :8080 后端）
make dev-api       # 仅后端（air 热重载）
make dev-web       # 仅前端
make test          # 后端测试 + 覆盖率
make test-web      # 前端 Vitest
make test-api      # REST Client 接口测试（VS Code / HttpYac 双路径）
make lint / fmt    # 静态检查 / 格式化
make seed          # 注入演示数据
make build && make run   # 生产构建与运行（后端托管前端产物）
make db-reset / clean    # 重置数据库 / 清理全部产物
```

### 10.2 阶段工时汇总

| 阶段 | 工时（人日） | 占比 |
| --- | --- | --- |
| M0 脚手架 | 1.0 | 13% |
| M1 后端 API | 2.0 | 27% |
| M2 前端基础 | 2.0 | 27% |
| M3 打磨 | 1.5 | 20% |
| M4 测试与收尾 | 1.0 | 13% |
| **合计** | **7.5** | 100% |

### 10.3 代码骨架参考（关键文件最小形态）

**后端路由（`router/router.go`）——注意静态路由先注册：**

```go
func Setup(r *gin.Engine, h *handler.TodoHandler, m ...gin.HandlerFunc) {
    r.Use(m...)
    r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
    g := r.Group("/api/v1")
    g.POST("/todos", h.Create)
    g.GET("/todos", h.List)
    g.DELETE("/todos/completed", h.DeleteCompleted) // 先于 :id
    g.GET("/todos/:id", h.Get)
    g.PATCH("/todos/:id", h.Patch)
    g.DELETE("/todos/:id", h.Delete)
}
```

**前端 store 骨架（`stores/todo.ts`）——乐观更新模式：**

```ts
async toggle(todo: Todo) {
  const prev = todo.completed
  todo.completed = !prev                     // 乐观更新
  todo.completedAt = prev ? null : new Date().toISOString()
  try {
    const updated = await todoApi.patch(todo.id, { completed: !prev })
    this.replaceItem(updated)                // 以响应体合并
  } catch (e) {
    todo.completed = prev                    // 失败回滚
    todo.completedAt = prev ? todo.completedAt : null
    throw e
  }
}
```

### 10.4 参考文档

- PRD：`spec/0001-prd-spec.md`（需求、接口契约、设计规范——一切实现的唯一事实来源）；
- 本计划与 PRD 冲突时，以 PRD 为准并向 PRD 提出修订。

---

*本文档为 Todo List 的实施基线；任务状态推进建议同步维护（可勾选各阶段表格或引入看板）。任何范围/工期调整需更新本计划版本记录。*
