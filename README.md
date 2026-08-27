# Todo List

一个极简、Apple 风格的**单用户待办事项应用**（无用户系统，本地 SQLite 存储）。

- **后端**：Go + Gin + GORM + SQLite3（`glebarez/sqlite`，纯 Go 无 CGO）
- **前端**：Vue 3 + TypeScript + Vite + Tailwind CSS + shadcn-vue
- **工程化**：Makefile 管理生命周期；REST Client（`server/test.rest`）验证 REST API

## 快速开始

前置条件：Go ≥ 1.22、Node ≥ 20、pnpm ≥ 9、GNU Make ≥ 4.0

```bash
make setup   # 安装前后端依赖
make dev     # 一键启动开发环境（前端 :5173，后端 :8080）
```

打开 http://localhost:5173 即可使用：输入回车创建、点击勾选完成、悬停删除、分段筛选、清除已完成、右上角切换浅色/深色/跟随系统。

> 演示数据：`make seed` 注入 12 条示例待办；`make db-reset` 重置数据库。

## 常用命令（`make help` 查看全部）

| 目标 | 说明 |
| --- | --- |
| `make setup` | 安装前后端依赖 |
| `make dev` / `dev-api` / `dev-web` | 启动开发环境（双端 / 仅后端 / 仅前端） |
| `make build` | 构建后端二进制 `bin/todolist-server` 与前端 `web/dist` |
| `make run` | 生产模式运行（后端单进程托管前端产物） |
| `make test` / `test-web` | 后端测试（含覆盖率）/ 前端 Vitest |
| `make test-api` | REST Client 接口测试（VS Code 中运行 `server/test.rest`） |
| `make lint` / `fmt` | 静态检查 / 格式化 |
| `make ci` | 本地 CI：格式检查 → 静态检查 → 测试 → 构建（`scripts/ci.sh`） |
| `make seed` / `db-reset` / `clean` | 演示数据 / 重置数据库 / 清理产物 |

## REST Client 测试

`server/test.rest` 包含全部 23 个接口用例（TC-01~TC-23，对应 PRD §10.3）：

1. 启动后端：`make dev-api`（或 `make dev`）；
2. 在 VS Code 中打开 `server/test.rest`，安装 [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) 插件后逐条点击 "Send Request"；
3. 或在 CI 中用 HttpYac 批量执行：`httpyac send server/test.rest --output short`。

用例覆盖：健康检查、创建（含空/超长标题）、列表（筛选/搜索/分页/排序/异常参数）、单查、完成/取消完成、改标题、删除、清除已完成、路由优先级。

## 目录结构

```text
├── Makefile / scripts/ci.sh     # 生命周期与 CI
├── spec/                        # 0001 需求设计 / 0002 实施计划
├── server/                      # Go 后端（cmd / internal / test.rest）
└── web/                         # Vue 前端（src/api·stores·features·components）
```

## 文档

| 文档 | 说明 |
| --- | --- |
| `spec/0001-prd-spec.md` | 需求与设计文档（PRD） |
| `spec/0002-implementation-plan.md` | 实施计划（M0~M4） |
