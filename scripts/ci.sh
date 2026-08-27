#!/usr/bin/env bash
# ============================================================
# Todo List 本地 CI：格式检查 → 静态检查 → 测试 → 构建
# 用法：make ci（或 bash scripts/ci.sh）
# ============================================================
set -euo pipefail

export PATH="$PATH:/usr/local/go/bin:/home/cxh/go/bin"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> [1/7] gofmt 检查"
UNFORMATTED="$(cd "$ROOT/server" && gofmt -l .)"
if [ -n "$UNFORMATTED" ]; then
  echo "以下文件需要 gofmt 格式化："
  echo "$UNFORMATTED"
  exit 1
fi

echo "==> [2/7] Prettier 检查"
(cd "$ROOT/web" && pnpm exec prettier --check .)

echo "==> [3/7] go vet"
(cd "$ROOT/server" && go vet ./...)

echo "==> [4/7] golangci-lint"
(cd "$ROOT/server" && golangci-lint run ./...)

echo "==> [5/7] ESLint"
(cd "$ROOT/web" && pnpm lint)

echo "==> [6/7] 后端测试（含覆盖率）"
(cd "$ROOT/server" && go test ./... -cover)

echo "==> [7/7] 前端测试"
(cd "$ROOT/web" && pnpm test)

echo "==> 构建"
make -C "$ROOT" build

echo "============================================"
echo "CI 全部通过 ✅"
