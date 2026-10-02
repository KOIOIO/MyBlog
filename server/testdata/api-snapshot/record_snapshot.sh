#!/usr/bin/env bash
# 录制重构前 HTTP 接口快照，用于 DDD 重构逐阶段回归对照。
# 用法: bash record_snapshot.sh [base_url] [out_dir]
set -u
BASE="${1:-http://127.0.0.1:8080}"
OUT="${2:-$(cd "$(dirname "$0")" && pwd)/captured}"
mkdir -p "$OUT"

# 公共公开接口（GET，带真实样本数据）
public_gets=(
  "GET /api/user/card?uuid=fbd5364d-bb15-11f1-b230-16b7a2303b52"
  "GET /api/article/4"
  "GET /api/article/search?query=Gin"
  "GET /api/article/category"
  "GET /api/article/tags"
  "GET /api/comment/2"
  "GET /api/comment/15"
  "GET /api/comment/new"
  "GET /api/feedback/new"
  "GET /api/forum/list"
  "GET /api/forum/detail?id=3"
  "GET /api/forum/tags"
  "GET /api/advertisement/info"
  "GET /api/friendLink/info"
  "GET /api/website/logo"
  "GET /api/website/title"
  "GET /api/website/info"
  "GET /api/website/carousel"
  "GET /api/website/news"
  "GET /api/website/calendar"
  "GET /api/website/footerLink"
)
# 公共 POST（无效请求体，捕获绑定错误行为）
public_posts=(
  "POST /api/base/captcha {}"
  "POST /api/base/sendEmailVerificationCode {}"
  "POST /api/user/forgotPassword {}"
  "POST /api/user/register {}"
  "POST /api/user/login {}"
)
# 登录后私有接口（无 token → 期望 NoAuth）
private_noauth=(
  "POST /api/user/logout {}"
  "PUT /api/user/resetPassword {}"
  "GET /api/user/info"
  "PUT /api/user/changeInfo {}"
  "POST /api/user/avatar {}"
  "GET /api/user/weather"
  "GET /api/user/chart?date=7"
  "POST /api/article/like {}"
  "GET /api/article/isLike"
  "GET /api/article/likesList"
  "POST /api/comment/create {}"
  "DELETE /api/comment/delete {}"
  "GET /api/comment/info"
  "POST /api/feedback/create {}"
  "GET /api/feedback/info"
  "POST /api/forum/publish {}"
  "POST /api/forum/upload {}"
  "POST /api/forum/like {}"
  "POST /api/forum/comment {}"
  "GET /api/forum/manageList"
  "GET /api/forum/manageComments"
  "DELETE /api/forum/delete {}"
  "DELETE /api/forum/comment {}"
)
# 管理接口（无 token → 期望 NoAuth/Forbidden）
admin_noauth=(
  "GET /api/user/list"
  "GET /api/user/loginList"
  "PUT /api/user/freeze {}"
  "PUT /api/user/unfreeze {}"
  "POST /api/article/create {}"
  "DELETE /api/article/delete {}"
  "PUT /api/article/update {}"
  "PUT /api/article/setTop {}"
  "GET /api/article/list"
  "GET /api/comment/list"
  "DELETE /api/feedback/delete {}"
  "PUT /api/feedback/reply {}"
  "GET /api/feedback/list"
  "POST /api/image/upload {}"
  "DELETE /api/image/delete {}"
  "GET /api/image/list"
  "POST /api/advertisement/create {}"
  "DELETE /api/advertisement/delete {}"
  "PUT /api/advertisement/update {}"
  "GET /api/advertisement/list"
  "POST /api/friendLink/create {}"
  "DELETE /api/friendLink/delete {}"
  "PUT /api/friendLink/update {}"
  "GET /api/friendLink/list"
  "POST /api/website/addCarousel {}"
  "PUT /api/website/cancelCarousel {}"
  "POST /api/website/createFooterLink {}"
  "DELETE /api/website/deleteFooterLink {}"
  "GET /api/config/website"
  "GET /api/config/system"
  "GET /api/config/email"
  "GET /api/config/qiniu"
  "GET /api/config/jwt"
  "GET /api/config/gaode"
  "PUT /api/config/website {}"
  "PUT /api/config/system {}"
  "PUT /api/config/email {}"
  "PUT /api/config/qiniu {}"
  "PUT /api/config/jwt {}"
  "PUT /api/config/gaode {}"
)

capture() {
  local method="$1" path="$2" body="$3" rec="$4"
  local tmp="$OUT/.tmp_body"
  local code
  if [ "$method" = "GET" ]; then
    code=$(curl -s -o "$tmp" -w "%{http_code}" "$BASE$path")
  else
    if [ -n "$body" ]; then
      code=$(curl -s -o "$tmp" -w "%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$body" "$BASE$path")
    else
      code=$(curl -s -o "$tmp" -w "%{http_code}" -X "$method" "$BASE$path")
    fi
  fi
  python3 "$(dirname "$0")/make_record.py" "$rec" "$method" "$path" "$code" "$tmp" > "$OUT/$rec.json"
}

i=0
for item in "${public_gets[@]}"; do
  i=$((i+1))
  path=$(echo "$item" | cut -d' ' -f2)
  capture "GET" "$path" "" "public_get_$i"
done
i=0
for item in "${public_posts[@]}"; do
  i=$((i+1))
  method=$(echo "$item" | cut -d' ' -f1); path=$(echo "$item" | cut -d' ' -f2); body=$(echo "$item" | cut -d' ' -f3)
  capture "$method" "$path" "$body" "public_post_$i"
done
i=0
for item in "${private_noauth[@]}"; do
  i=$((i+1))
  method=$(echo "$item" | cut -d' ' -f1); path=$(echo "$item" | cut -d' ' -f2); body=$(echo "$item" | cut -d' ' -f3)
  capture "$method" "$path" "$body" "private_noauth_$i"
done
i=0
for item in "${admin_noauth[@]}"; do
  i=$((i+1))
  method=$(echo "$item" | cut -d' ' -f1); path=$(echo "$item" | cut -d' ' -f2); body=$(echo "$item" | cut -d' ' -f3)
  capture "$method" "$path" "$body" "admin_noauth_$i"
done
rm -f "$OUT/.tmp_body"
echo "SNAPSHOT DONE: $(ls "$OUT"/*.json | wc -l | tr -d ' ') records in $OUT"
