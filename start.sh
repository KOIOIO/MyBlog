#!/usr/bin/env bash
# MyBlog-ES 一键启动脚本：基础设施(Docker) + 后端(Go) + 前端(Vite)
# 用法: ./start.sh  或  ./stop.sh
set -euo pipefail

PROJECT_DIR="/Users/xiaoyuwang/projects/MyBlog"
SERVER_DIR="$PROJECT_DIR/server"
WEB_DIR="$PROJECT_DIR/web"
BACKEND_BIN="/tmp/myblog-server"
BACKEND_LOG="/tmp/myblog-server.log"
WEB_LOG="/tmp/myblog-web.log"

ensure_container() {
  local name="$1"; shift
  if docker ps -a --format '{{.Names}}' | grep -qx "$name"; then
    docker start "$name" >/dev/null
  else
    docker run -d --name "$name" "$@" >/dev/null
  fi
  echo "[$name] started"
}

wait_ready() {
  local url="$1" name="$2" i=0
  until curl -s -o /dev/null "$url"; do
    i=$((i+1))
    if [ "$i" -gt 60 ]; then echo "[$name] not ready after 60s"; exit 1; fi
    sleep 2
  done
  echo "[$name] ready"
}

echo "==> 1/4 启动基础设施 (Docker)"
ensure_container myblog-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=blog_db mysql:8.0
ensure_container myblog-redis -p 6379:6379 redis:7-alpine
ensure_container myblog-es -p 9200:9200 -e discovery.type=single-node -e xpack.security.enabled=false -e ES_JAVA_OPTS="-Xms512m -Xmx512m" docker.elastic.co/elasticsearch/elasticsearch:8.17.0
wait_ready "http://127.0.0.1:9200/" "elasticsearch"

echo "==> 2/4 启动后端 (Go)"
cd "$SERVER_DIR"
if ! lsof -iTCP:8080 -sTCP:LISTEN >/dev/null 2>&1; then
  go build -o "$BACKEND_BIN" .
  nohup "$BACKEND_BIN" > "$BACKEND_LOG" 2>&1 &
  sleep 3
  echo "[backend] started, log: $BACKEND_LOG"
else
  echo "[backend] already running on :8080"
fi

echo "==> 3/4 启动前端 (Vite)"
cd "$WEB_DIR"
if ! lsof -iTCP:80 -sTCP:LISTEN >/dev/null 2>&1; then
  [ -d node_modules ] || npm install
  nohup npm run dev > "$WEB_LOG" 2>&1 &
  sleep 3
  echo "[frontend] started, log: $WEB_LOG"
else
  echo "[frontend] already running on :80"
fi

echo "==> 4/4 完成"
echo "前端: http://localhost:80   后端: http://127.0.0.1:8080"
echo "管理员登录: 邮箱 2652777599@qq.com / 密码 admin123456 (请及时修改)"
