# hengshui-tablet-video-rpc（MyBlog-ES 博客）云服务器部署手册

> 目标：把本地这套 
>
> **Go 后端 + Vue3 前端 + MySQL + Redis + Elasticsearch**
>
>  的博客系统，连同现有数据，完整迁移到云服务器，并用 1Panel 统一管理。



***

## 0. 部署架构总览（先理解，再动手）



```
&#x20;                      公网用户

&#x20;                         │  HTTPS (80/443)

&#x20;                         ▼

&#x20;                 DNS: 你的域名 A 记录 → 服务器 IP

&#x20;                         ▼

&#x20;        ┌─────────────────────────────────────┐

&#x20;        │         1Panel 管理的 OpenResty      │   ← 1Panel「网站」自动安装

&#x20;        │  静态文件: /opt/myblog/web-dist      │     托管前端 dist，自动续期 SSL

&#x20;        │  反向代理: /api、/uploads → 127.0.0.1:8080

&#x20;        └──────────────┬──────────────────────┘

&#x20;                       │ 127.0.0.1:8080（仅本机，不暴露公网）

&#x20;                       ▼

&#x20;        ┌─────────────────────────────────────┐

&#x20;        │   backend 容器 (Go:8080)            │  ← 1Panel「编排」构建启动

&#x20;        │   挂载 config.yaml / uploads / logs │

&#x20;        └───────┬───────────┬───────────┬────┘

&#x20;                │           │           │

&#x20;                ▼           ▼           ▼

&#x20;             MySQL        Redis      Elasticsearch

&#x20;          (blog\_db)    (缓存,可重建)   (article\_index)
```

**核心原理（为什么要这样分层）：**



* **DNS + 反向代理**：域名只负责「把流量引到哪台服务器」。服务器上真正对外收流量的进程是 Nginx（1Panel 用 OpenResty），它同时做两件事：直接返回前端静态文件（快），把 `/api`、`/uploads` 转发给后端 Go 进程（动态数据）。这样前端、后端、数据库不需要各自暴露公网端口，攻击面最小。

* **为什么后端只绑&#x20;**`127.0.0.1:8080`：OpenResty 和后端都在同一台机器上，通过回环地址通信即可；公网只留 80/443，防止别人直接打你的 API 或数据库。

* **为什么 MySQL/Redis/ES 用容器服务名（**`mysql`**、**`redis`**、**`es`**）而不是 IP**：同一 compose 网络内，服务名即域名，由 Docker 内置 DNS 解析；IP 会随容器重建变化，服务名不变。

* **数据迁移的本质**：MySQL 是「逻辑导出」`mysqldump`（把表结构和数据写成 SQL）；Elasticsearch 用项目自带的 `--es-export`（滚动查询导出 JSON）；上传文件就是打包目录。三者都是纯数据，跟运行环境无关，拿到新环境再导入即可。

* **Redis 不需要迁移**：它存的是缓存、会话、JWT 黑名单等可重建数据（黑名单在服务启动时从数据库加载）。新环境 Redis 是空的，不影响正确性。



***

## 1. 服务器准备



| 项目 | 要求                                                        | 说明                                                 |
| -- | --------------------------------------------------------- | -------------------------------------------------- |
| 系统 | 64 位 Linux（Ubuntu 20.04+/22.04、Debian 11+、CentOS 7.9+ 均可） | 1Panel 官方支持                                        |
| 内存 | **≥ 4G 推荐，最低 2G**                                         | ES (JVM 512M) + MySQL + Redis + Go 后端 + 面板，2G 会比较紧 |
| 磁盘 | ≥ 40G                                                     | ES 数据 + MySQL + 日志                                 |
| 端口 | 22、80、443、1Panel 面板端口                                     | 云厂商安全组 + 1Panel 防火墙都要放行                            |

**安全组放行**（在云厂商控制台操作，别漏）：



* 22（SSH）、80（HTTP）、443（HTTPS）、1Panel 面板端口（安装时终端会显示，一般是个 5 位随机端口）。



***

## 2. 安装 1Panel

SSH 登录服务器后执行（国内服务器用国内源）：



```
curl -sSL https://resource.1panel.cn/quick\_start.sh -o quick\_start.sh

sudo bash quick\_start.sh
```

> 原理：脚本会检测系统、安装 Docker 引擎和 1Panel 本体（面板本身也跑在 Docker 里），并随机生成面板端口 / 安全入口。安装结束时
>
> **记下**
>
> 终端打印的「面板地址、用户名、密码」。

安装后：



* 浏览器访问 `https://服务器IP:面板端口` 登录（第一次需要创建安全入口）。

* 1Panel 的「网站」功能在你创建第一个网站时会自动安装 OpenResty，不用提前装。



***

## 3. 域名解析与备案（内网服务器硬前置）

**如果你的域名还没有 ICP 备案，或服务器位于中国大陆：**

> 原理：按工信部要求，域名解析到境内服务器并用 80/443 提供 Web 服务，必须完成 ICP 备案；云厂商会在未备案访问时拦截（显示备案提示页）。所以这是
>
> **上线前必须完成**
>
> 的一步，通常 1~2 周。备案期间可先把后续所有准备工作做完。



1. 在域名服务商（阿里云 / 腾讯云 / 西部数码等）控制台，添加解析记录：

* 主机记录 `@` → 记录类型 `A` → 记录值 `你的服务器公网IP`

* 主机记录 `www` → 记录类型 `A` → 记录值 `你的服务器公网IP`

1. 去服务器厂商备案系统提交 ICP 备案（按提示填主体信息，域名需先实名认证）。

2. 若域名已在其他厂商备案过，解析到新厂商服务器时做「接入备案」。

3. 备案通过后，把 `deploy/config.yaml.prod` 里 `website.icp_filing`、`public_security_filing` 填上真实备案号，页脚会正确展示。

> 提示：你本地配置里 QQ 登录已从代码中移除（
>
> `config/conf_qq.go`
>
>  已删除），部署配置模板也不再包含 qq 段，无需处理旧的回调地址。



***

## 4. 上传项目到服务器

在 **Mac 终端**执行。只需要后端源码 + 部署目录（前端 dist 后面单独构建上传）：



```
\# 假设服务器 IP 为 1.2.3.4，用户名 ubuntu（按你的实际改）

scp -r \~/projects/MyBlog/server root@1.2.3.4:/opt/myblog/server

scp -r \~/projects/MyBlog/deploy root@1.2.3.4:/opt/myblog/deploy
```

> 原理说明：
>
> `deploy/`
>
>  是部署根目录，compose 里的相对路径（
>
> `./mysql-data`
>
> 、
>
> `./uploads`
>
> 、
>
> `./logs`
>
> 、
>
> `./backup`
>
> ）都会落在 
>
> `/opt/myblog/deploy/`
>
>  下，数据文件集中在同一处，方便备份和查看。
> 备选：如果 GitHub 仓库（
>
> `git@github.com:KOIOIO/MyBlog.git`
>
> ）已是最新，也可在服务器 
>
> `git clone`
>
>  后，把 
>
> `deploy/`
>
>  目录拷贝过去。注意：
>
> **不要把含密钥的&#x20;**
>
> `config.yaml`
>
> **&#x20;提交到公开仓库**
>
> 。



***

## 5. 修改配置（服务器上操作）

### 5.1 密码：`.env`



```
cd /opt/myblog/deploy

cp .env.example .env

vi .env        # 把 MYSQL\_ROOT\_PASSWORD 改成强密码
```

### 5.2 后端配置：`config.yaml.prod`

按文件内「必改 / 从本地复制」标注修改：



* `mysql.password` ← 与 `.env` 的 `MYSQL_ROOT_PASSWORD` **保持一致**；

* `es.url`、`redis.address`、`mysql.host` 保持 `es:9200` / `redis:6379` / `mysql:3306`（容器服务名，别改成 127.0.0.1）；

* JWT / SMTP / 七牛 / 高德密钥：从本地 `server/config.yaml` **原样复制**（保持 JWT 密钥不变，用户已登录的 token 依然有效，不用重新登录）；

* `website.*`：按你的新域名 / 备案号修改。

> 原理：后端在启动目录读取固定文件名 
>
> `config.yaml`
>
> ，容器里工作目录是 
>
> `/app`
>
> ，compose 把 
>
> `config.yaml.prod`
>
>  挂载成 
>
> `/app/config.yaml`
>
> （只读）。改配置 = 改宿主机文件，重启容器即生效，不用重新编译。



***

## 6. 用 1Panel 编排启动基础设施 + 后端

**前置：ES 需要调高宿主机的内存映射上限**



```
sudo sysctl -w vm.max\_map\_count=262144

echo "vm.max\_map\_count=262144" | sudo tee /etc/sysctl.d/99-elasticsearch.conf
```

> 原理：Elasticsearch 的 Lucene 底层用 
>
> `mmap`
>
>  映射索引文件，内核默认每个进程最多 65530 个映射区，ES 会因不足直接拒绝启动。这个参数是宿主机内核参数，不是容器配置。

**在 1Panel 面板操作：**



1. 左侧「容器」→「编排」→「创建编排」；

2. 名称填 `myblog`，选择服务器上的编排文件 `/opt/myblog/deploy/docker-compose.yml`（或把文件内容粘贴进去）；

3. 确认后 1Panel 执行 `docker compose up -d --build`：

* 拉取 mysql:8.0 /redis:7 /elasticsearch:8.17.0 镜像；

* **首次构建后端镜像**（自动拉 golang 基础镜像 → 编译 Go 二进制 → 生成 alpine 运行镜像），约几分钟，可在编排详情「日志」里看进度；

* 按依赖顺序启动：MySQL/Redis/ES 健康检查通过后才启动 backend（`depends_on: condition: service_healthy`）。

启动后验证（SSH 或 1Panel 终端）：



```
cd /opt/myblog/deploy

docker compose ps          # 四个服务都应为 running / healthy

curl http://127.0.0.1:8080/api/...   # 后端可达（具体接口见下方验收）
```



***

## 7. 数据迁移（核心：把本地数据搬过去）

### 7.1 本地导出（Mac 上执行）



```
cd \~/projects/MyBlog

chmod +x deploy/migrate-export.sh

./deploy/migrate-export.sh
```

产物在 `deploy/backup/`：



* `blog_db.sql`：MySQL 全库（表结构 + 数据）；

* `es_data.json`：ES 文章索引全文数据（用项目自带 `--es-export` 滚动导出）；

* `uploads.tar.gz`：本地上传的图片 / 头像等文件。

> 原理：
>
> `mysqldump --single-transaction`
>
>  在导出期间不锁表、保证一致性快照；ES 导出用 Scroll 滚动接口分批拉取，避免一次性大结果集压垮内存。

### 7.2 上传备份到服务器



```
scp -r \~/projects/MyBlog/deploy/backup root@1.2.3.4:/opt/myblog/deploy/
```

### 7.3 服务器导入



```
cd /opt/myblog/deploy

\# ① MySQL：导入 SQL（管道直接喂给容器内的 mysql 客户端）

docker compose exec -T mysql sh -c 'exec mysql -uroot -p"\$MYSQL\_ROOT\_PASSWORD"' < backup/blog\_db.sql

\# ② Elasticsearch：先建索引，再导入数据

\#    --es 若提示索引已存在并问 y/n，输入 y（会删掉空索引重建）

echo y | docker compose exec -T backend ./myblog-server --es

docker compose exec -T backend ./myblog-server --es-import /app/backup/es\_data.json

\# ③ 上传文件：解压到部署目录（后端容器把 ./uploads 挂载为 /app/uploads，立即可见）

tar -xzf backup/uploads.tar.gz -C .
```

验证导入：



```
docker compose exec mysql sh -c 'exec mysql -uroot -p"\$MYSQL\_ROOT\_PASSWORD" -e "use blog\_db; show tables;"'

docker compose exec es curl -s http://127.0.0.1:9200/article\_index/\_count

ls uploads/image | head
```

> 原理：容器重启后数据不会丢，因为 MySQL 数据在 
>
> `./mysql-data`
>
>  绑定目录、ES 在 
>
> `esdata`
>
>  命名卷、上传文件在 
>
> `./uploads`
>
> ；它们都独立于容器生命周期。将来容器重建，数据原地保留。



***

## 8. 构建并上传前端

前端在 **Mac 本地构建**（不需要在服务器装 Node）：



```
cd \~/projects/MyBlog/web

npm run build          # 产物在 web/dist/

\# 若 type-check 报 TS 错误且不影响功能，可用：npm run build-only（仅 vite 构建）
```

上传并放置到网站根目录：



```
scp -r \~/projects/MyBlog/web/dist root@1.2.3.4:/opt/myblog/web-dist
```

> 原理：Vite 构建时把 
>
> `VITE_BASE_API=/api`
>
>  编译进产物，所以线上前端请求的接口地址就是
>
> **同源的&#x20;**
>
> `/api`
>
> —— 这正是 Nginx 反向代理要转发给后端的路径。开发时的 
>
> `VITE_SERVER_URL=127.0.0.1:8080`
>
>  只影响 vite dev 代理，不影响生产包。



***

## 9. 1Panel 创建网站 + 反向代理 + HTTPS

在 1Panel 面板操作：



1. 左侧「网站」→「创建网站」→ 选「静态网站」：

* 主域名：`你的域名`（建议同时勾选 `www.你的域名`）；

* 网站目录：`/opt/myblog/web-dist`；

* 创建时 1Panel 自动安装 / 启动 OpenResty。

1. 进入该网站的「设置」→「反向代理」，添加两条（都指向同一后端，1Panel 会生成 location 规则）：

* 路径 `/api` → `http://127.0.0.1:8080`

* 路径 `/uploads` → `http://127.0.0.1:8080`

1. 「证书」→ 创建证书 → ACME（Let's Encrypt）：

* 首选 **DNS 验证**（在 1Panel 里配置你的域名服务商 API 密钥，可自动添加解析记录完成验证，最稳）；

* 或 **HTTP 验证**（要求 80 端口已放行、A 记录已生效）；

* 签发后把证书关联到网站，开启「强制 HTTPS」，1Panel 会自动续期。

1. 可选：网站「设置」→ 伪静态 / 缓存默认即可；确认 `client_max_body_size` 不小于上传限制（后端配置 `upload.size=20`，即 20MB）。

> 原理：HTTPS 证书验证的本质是「证明你控制这个域名」——DNS 验证要求你在域名的 DNS 记录里放一段 TXT 值；HTTP 验证要求 80 端口能返回指定文件。1Panel 的 ACME 客户端把这两件事自动化了，到期前自动续期。



***

## 10. 验收清单



* [ ] 浏览器访问 `https://你的域名`，前端正常加载（logo、样式、图片）；

* [ ] 文章列表能打开，搜索能出结果（ES 数据已导入）；

* [ ] 文章配图 / 头像正常显示（uploads 已迁移）；

* [ ] 用管理员账号登录（本地是 `2652777599@qq.com`，登录后**立即改密码**）；

* [ ] 后台发布一篇文章，带图片上传，确认图片落到 `/opt/myblog/deploy/uploads/image/`；

* [ ] 直接访问 `http://服务器IP:8080` 应**无法**打开（后端只绑了 127.0.0.1，符合预期）；

* [ ] `docker compose ps` 全部 healthy。



***

## 11. 日常运维

### 备份（建议做）

在 1Panel「计划任务」里加两条每日任务：



* Shell：`cd /opt/myblog/deploy && docker compose exec -T mysql sh -c 'exec mysqldump -uroot -p"$MYSQL_ROOT_PASSWORD" --databases blog_db' > backup/daily_$(date +%F).sql`（保留 N 天）；

* Shell：ES 数据用 `docker compose exec -T backend ./myblog-server --es-export` 导出（产物在容器 /app，用 `docker compose cp backend:/app/es_*.json backup/` 拷出）；

* 目录备份：`tar -czf backup/uploads_$(date +%F).tar.gz uploads`。

> 备份文件定期下载到本地或对象存储，别只留在服务器上（服务器故障 = 备份一起丢）。

### 更新（改代码后）



```
\# 后端：本地重新打包上传 server/，然后

cd /opt/myblog/deploy && docker compose build backend && docker compose up -d backend

\# 前端：本地 npm run build 后

scp -r dist root@1.2.3.4:/opt/myblog/web-dist

\# 配置：改 config.yaml.prod 后

docker compose restart backend
```

### 看日志



* 1Panel 编排详情 → 服务 → 日志；

* 或 `docker compose logs -f backend`。



***

## 12. 常见问题（FAQ）



| 现象             | 原因                        | 解决                                                                                          |
| -------------- | ------------------------- | ------------------------------------------------------------------------------------------- |
| ES 容器反复重启 / 退出 | 宿主机 `vm.max_map_count` 不足 | 执行第 6 节前置命令并重启服务器                                                                           |
| ES 健康检查一直失败    | 镜像内无 curl（个别版本）           | 把 healthcheck 改为：\`test: \["CMD-SHELL","exec 3<>/dev/tcp/127.0.0.1/9200 && head -1 <&3      |
| 2G 内存机器太卡      | MySQL+ES + 后端挤内存          | 调小 `ES_JAVA_OPTS` 为 `-Xms256m -Xmx256m`；给 MySQL 加 `command: --innodb-buffer-pool-size=256M` |
| 访问域名显示 502     | 反向代理指向的后端没起来              | `docker compose ps`；`curl http://127.0.0.1:8080` 自测                                         |
| 图片 404         | uploads 没迁移或路径不对          | 检查 `/opt/myblog/deploy/uploads/image` 是否有文件；确认反代了 `/uploads`                                |
| 搜索无结果          | ES 索引没建 / 没导入             | 重跑 7.3 第②步，`curl http://127.0.0.1:9200/article_index/_count`                                |
| 备案没通过前访问被拦截    | 云厂商备案拦截                   | 等待备案通过；期间可先用 IP 加端口临时联调（不改生产入口）                                                             |
| 上传文件失败         | OpenResty 默认 body 大小限制    | 网站设置里调大 `client_max_body_size`（≥ 20m）                                                       |



***

## 附：文件清单



```
deploy/

├── docker-compose.yml     # 四服务编排：mysql/redis/es/backend

├── .env.example           # 密码模板 → 复制为 .env

├── config.yaml.prod       # 生产后端配置模板

├── migrate-export.sh      # Mac 本地数据导出脚本

├── DEPLOY.md              # 本手册

└── backup/                # 导出/备份产物（运行后生成）

server/Dockerfile          # 后端多阶段构建镜像
```