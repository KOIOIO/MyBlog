#!/usr/bin/env bash
# hengshui-tablet-video-rpc 本地数据导出脚本（在 Mac 上运行）
# 导出内容：MySQL(blog_db) + Elasticsearch(article_index) + 本地上传文件(uploads)
# 前提：Docker 可用、本机 Go 环境可用；建议先运行 ./start.sh 让三个容器处于运行状态。
set -euo pipefail

PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BACKUP_DIR="$PROJECT_DIR/deploy/backup"
mkdir -p "$BACKUP_DIR"

# 0. 确保本地三个基础设施容器在运行（导出工具需要它们在线）
for c in myblog-mysql myblog-redis myblog-es; do
  docker start "$c" >/dev/null 2>&1 || true
done

echo "==> 1/3 导出 MySQL (blog_db)"
docker exec myblog-mysql sh -c 'exec mysqldump -uroot -proot --single-transaction --databases blog_db' > "$BACKUP_DIR/blog_db.sql"
echo "    -> $BACKUP_DIR/blog_db.sql ($(wc -c < "$BACKUP_DIR/blog_db.sql" | tr -d ' ') bytes)"

echo "==> 2/3 导出 Elasticsearch (article_index)"
(
  cd "$PROJECT_DIR/server"
  go build -o /tmp/myblog-server .
  /tmp/myblog-server --es-export
)
mv "$PROJECT_DIR/server"/es_*.json "$BACKUP_DIR/es_data.json"
echo "    -> $BACKUP_DIR/es_data.json"

echo "==> 3/3 打包本地上传文件 (uploads)"
tar -czf "$BACKUP_DIR/uploads.tar.gz" -C "$PROJECT_DIR/server" uploads
echo "    -> $BACKUP_DIR/uploads.tar.gz"

echo ""
echo "导出完成，产物都在 $BACKUP_DIR"
echo "下一步：把整个 deploy/ 目录上传到服务器 /opt/myblog/（scp -r 或 rsync）"
