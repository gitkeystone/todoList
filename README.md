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

打开 http://localhost:5173 即可使用。

## 常用命令

```bash
make help    # 查看全部生命周期目标
make test    # 后端单元测试（含覆盖率）
make test-web# 前端 Vitest
make test-api# REST Client 接口测试（VS Code 中运行 server/test.rest）
make build && make run   # 生产构建与运行
```

## 文档

| 文档 | 说明 |
| --- | --- |
| `spec/0001-prd-spec.md` | 需求与设计文档（PRD） |
| `spec/0002-implementation-plan.md` | 实施计划（M0~M4） |
