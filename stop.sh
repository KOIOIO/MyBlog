#!/usr/bin/env bash
# MyBlog-ES 一键停止脚本：停前端/后端进程并停止 Docker 容器
set -uo pipefail

echo "==> 停止后端/前端进程"
pkill -f "myblog-server" 2>/dev/null && echo "[backend] stopped" || echo "[backend] not running"
pkill -f "vite" 2>/dev/null && echo "[frontend] stopped" || echo "[frontend] not running"

echo "==> 停止 Docker 容器"
for c in myblog-mysql myblog-redis myblog-es; do
  docker stop "$c" >/dev/null 2>&1 && echo "[$c] stopped" || echo "[$c] not running"
done
echo "done"
