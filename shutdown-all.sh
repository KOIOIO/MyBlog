#!/usr/bin/env bash
# MyBlog 一键关闭脚本：停后端/前端进程 + 停 Docker 容器 + 强制释放项目端口 + 残留检查
# 用法: ./shutdown-all.sh
set -uo pipefail

echo "==> 1/3 停止后端/前端进程"
# 兼容新旧二进制名（myblog-server / blog-serv）
pkill -f "myblog-server" 2>/dev/null && echo "[backend] stopped" || echo "[backend] not running"
pkill -f "blog-serv" 2>/dev/null && echo "[backend legacy] stopped" || true
pkill -f "vite" 2>/dev/null && echo "[frontend] stopped" || echo "[frontend] not running"

echo "==> 2/3 停止 Docker 容器"
for c in myblog-mysql myblog-redis myblog-es; do
  docker stop "$c" >/dev/null 2>&1 && echo "[$c] stopped" || echo "[$c] not running"
done

echo "==> 3/3 强制释放项目端口并检查残留"
for port in 8080 80; do
  if lsof -tiTCP:$port -sTCP:LISTEN >/dev/null 2>&1; then
    lsof -tiTCP:$port -sTCP:LISTEN | xargs kill 2>/dev/null
    echo "[port $port] force released"
  else
    echo "[port $port] released"
  fi
done
sleep 1
if pgrep -f "myblog-server|blog-serv|vite" >/dev/null 2>&1; then
  echo "[warning] residual process detected"
  pgrep -fl "myblog-server|blog-serv|vite"
else
  echo "[ok] no residual process"
fi
echo "done"
