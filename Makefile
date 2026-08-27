# =====================================================================
# Todo List 生命周期管理（唯一入口）
# 目标一览：make help
# =====================================================================

# ---- 配置 ----
SERVER_DIR  := server
WEB_DIR     := web
BIN         := bin/todolist-server
PORT        ?= 8080
DB_PATH     ?= data/todolist.db
WEB_DIST    ?= web/dist

.PHONY: help setup dev dev-api dev-web build run test test-web test-api lint fmt seed db-reset clean ci

help:            ## 显示全部目标与说明
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup:           ## 安装前后端依赖
	cd $(SERVER_DIR) && go mod download
	cd $(WEB_DIR) && pnpm install

dev:             ## 一键启动开发环境（前端 :5173 + 后端 :8080，Ctrl+C 全部退出）
	@trap 'kill 0' INT TERM EXIT; \
	$(MAKE) dev-api & $(MAKE) dev-web & wait

dev-api:         ## 启动后端（有 air 则热重载，否则 go run）
	@command -v air >/dev/null 2>&1 && (cd $(SERVER_DIR) && DB_PATH=../$(DB_PATH) air --build.cmd "go build -o ../$(BIN) ./cmd/api" --build.bin "../$(BIN)") || (cd $(SERVER_DIR) && DB_PATH=../$(DB_PATH) go run ./cmd/api)

dev-web:         ## 启动前端（Vite dev server）
	cd $(WEB_DIR) && pnpm dev

build:           ## 构建前后端产物
	cd $(SERVER_DIR) && go build -o ../$(BIN) ./cmd/api
	cd $(WEB_DIR) && pnpm build

run:             ## 生产模式运行（先 build；后端托管 web/dist）
	cd $(SERVER_DIR) && GIN_MODE=release DB_PATH=../$(DB_PATH) WEB_DIST=../$(WEB_DIST) ../$(BIN)

test:            ## 后端单元测试（含覆盖率）
	cd $(SERVER_DIR) && go test ./... -cover

test-web:        ## 前端单元测试（Vitest）
	cd $(WEB_DIR) && pnpm test

test-api:        ## REST Client 接口测试（VS Code REST Client / HttpYac）
	@echo "方式A：在 VS Code 打开 server/test.rest 逐条执行"
	@echo "方式B（CI）：httpyac send server/test.rest --output short"

lint:            ## 静态检查（golangci-lint + ESLint）
	cd $(SERVER_DIR) && golangci-lint run ./...
	cd $(WEB_DIR) && pnpm lint

fmt:             ## 代码格式化（gofmt + Prettier）
	cd $(SERVER_DIR) && gofmt -w .
	cd $(WEB_DIR) && pnpm format

seed:            ## 注入演示数据
	cd $(SERVER_DIR) && DB_PATH=../$(DB_PATH) go run ./cmd/seed

db-reset:        ## 重置数据库（开发用）
	rm -f $(DB_PATH)

ci:              ## 本地 CI：格式/静态检查/测试/构建 全流程
	bash scripts/ci.sh

clean:           ## 清理构建产物与数据库
	rm -rf bin web/dist $(DB_PATH)
