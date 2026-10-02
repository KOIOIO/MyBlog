# MyBlog 阿里云 ECS + 1Panel 部署手把手指南

> 服务代号：hengshui-tablet-video-rpc（MyBlog：Go+Gin 后端 + Vue3 前端 + MySQL8 + Redis7 + ES8.17）
> 更新日期：2026-09-29



***

## 0. 当前进度一览（2026-09-29 已核实）



| 步骤                          | 状态        | 说明                                                              |
| --------------------------- | --------- | --------------------------------------------------------------- |
| 服务器连接（Workbench CLI）        | ✅ 已完成     | 实例 `i-f8z9zmvj3iow727wjmpd`，免密可连                                |
| 1Panel 安装                   | ✅ 已完成     | 面板信息见第 3 节，**登录后请立即改密码**                                        |
| 项目文件上传 `/opt/myblog`        | ✅ 已完成     | server/ 与 deploy/ 均已就位                                          |
| `.env` + `config.yaml.prod` | ✅ 已完成     | 密钥已从本地复制，DB 密码已生成（见附录 B）                                        |
| ES 内核参数 `vm.max_map_count`  | ✅ 已完成     | 已设置并持久化                                                         |
| 基础设施镜像拉取                    | ✅ 已完成     | mysql:8.0 /redis:7-alpine /es:8.17.0 已在服务器                      |
| 本地数据导出                      | ✅ 已完成     | 产物在 `deploy/backup/`（blog\_db.sql、es\_data.json、uploads.tar.gz） |
| **安全组放行端口**                 | ⬜ **待你做** | 见第 1 步                                                          |
| **DNS 解析改 A 记录**            | ⬜ **待你做** | 见第 2 步                                                          |
| **启动整栈编排**                  | ⬜ 待执行     | 见第 4 步（backend 镜像尚未构建）                                          |
| **数据上传与导入**                 | ⬜ 待执行     | 见第 5 步                                                          |
| **前端构建上传**                  | ⬜ 待执行     | 见第 6 步                                                          |
| **1Panel 建站 + HTTPS**       | ⬜ 待执行     | 见第 7 步                                                          |
| 验收                          | ⬜ 待执行     | 见第 8 步                                                          |

> 重要提示：如果你在编辑器里打开本文件，
> **请关闭后重新打开一次**
> ，确保命令没有被编辑器自动转义（正常命令里不应出现
> `\$`
> 、
> `\_`
> 这样的反斜杠）。粘贴命令到终端后，先扫一眼有没有多余的
> `\`
> 。



***

## 服务器环境速记



| 项目    | 值                                                         |
| ----- | --------------------------------------------------------- |
| 公网 IP | `47.120.64.37`                                            |
| 实例 ID | `i-f8z9zmvj3iow727wjmpd`                                  |
| 地域    | 华南 2（河源）A・`cn-heyuan`                                     |
| 系统    | Ubuntu 24.04・2 核 4G                                       |
| 部署目录  | `/opt/myblog`（server/ 源码 + deploy/ 编排）                    |
| 数据卷目录 | `/opt/myblog/deploy`（mysql-data、uploads、logs、backup 都在这里） |



***

## 第 0 步：连接服务器（已配好，复习用）

本地 Mac 已装 Workbench CLI，直接免密连接：



```
workbench connect -i i-f8z9zmvj3iow727wjmpd
```



* 这是**交互式登录**，进去后可以粘贴长命令、看实时输出。

* 单独执行一条命令（非交互）用：



```
workbench exec -i i-f8z9zmvj3iow727wjmpd -c "命令" --timeout 60
```



* `workbench exec` **默认只有 30 秒超时**，跑编译 / 拉镜像这类长命令必须加 `--timeout 600`。



***

## 第 1 步：阿里云安全组放行端口（网页操作）

**做什么**：给服务器放行 1Panel 面板端口、网站端口。

**为什么**：安全组是阿里云的 "防火墙"，不放行端口等于门外上锁，浏览器和域名都进不来。



1. 打开 阿里云控制台 → ECS → 实例 → 找到 `i-f8z9zmvj3iow727wjmpd`

2. 点「安全组」→ 所在安全组 →「管理规则」→「入方向」→「手动添加」

3. 放行下面 4 条（授权对象都填 `0.0.0.0/0`）：



| 协议  | 端口范围  | 用途                       |
| --- | ----- | ------------------------ |
| TCP | 30630 | 1Panel 面板（你的面板端口，见第 3 节） |
| TCP | 22    | SSH（保险起见放行）              |
| TCP | 80    | 网站 HTTP                  |
| TCP | 443   | 网站 HTTPS                 |

> 别一次全开完也行，但
> **80/443 必须在第 7 步申请证书前放行**
> ，30630 在登录面板前放行。



***

## 第 2 步：DNS 解析（网页操作）

**做什么**：把域名 `@` 和 `www` 两条 A 记录**改成** `47.120.64.37`。

**为什么**：现在这两条记录还指向旧服务器 `47.97.117.113`。注意是**修改**，不是**新增**—— 同一主机名存在两条 A 记录会做轮询，一半流量仍会跑到旧服务器。



1. 打开 阿里云控制台 → 域名 → 解析设置

2. 找到主机记录 `@` 的 A 记录 → 修改 → 记录值填 `47.120.64.37` → 确定

3. 找到主机记录 `www` 的 A 记录 → 修改 → 记录值填 `47.120.64.37` → 确定

4. 其余 `_dnsauth` TXT、comodoca CNAME 是旧证书的验证记录，留删皆可（不影响）

5. 验证（等 1\~10 分钟生效）：



```
dig +short 你的域名

dig +short www.你的域名
```

两条都输出 `47.120.64.37` 就对了。



***

## 第 3 步：登录 1Panel 并改默认密码

**面板信息（安装时生成的，务必保存）：**



| 项目   | 值                                      |
| ---- | -------------------------------------- |
| 面板地址 | `http://47.120.64.37:30630/0cd9cdd4d7` |
| 用户名  | `wwy`                                  |
| 密码   | `af225bfcb1`                           |



1. 浏览器打开面板地址，登录

2. 右上角头像 → 修改密码 → 换一个只有你知道的强密码

3. 顺手确认：面板里「容器」页面能看到 Docker 已经在跑

> 面板地址里的
> `/0cd9cdd4d7`
> 是 "安全入口"，不要删掉它访问，否则会被拦截。



***

## 第 4 步：启动整栈（服务器上执行）

**做什么**：用 docker compose 把 MySQL、Redis、ES、后端四个服务全部拉起来。

**为什么**：文件已经全部就位（/opt/myblog），这一步会构建后端镜像（约 2\~5 分钟）并启动全部容器。

在 connect 会话里执行（推荐，能看到实时日志）：



```
cd /opt/myblog/deploy && docker compose up -d --build
```

看到类似输出就成功：



```
\\\[+] Running 5/5

\&#x20;✔ Container deploy-mysql-1    Started

\&#x20;✔ Container deploy-redis-1    Started

\&#x20;✔ Container deploy-es-1       Started

\&#x20;✔ Container deploy-backend-1  Started
```

> 如果是在
> `workbench exec`
> 里跑（非交互），必须加超时参数：



```
workbench exec -i i-f8z9zmvj3iow727wjmpd -c "cd /opt/myblog/deploy && docker compose up -d --build" --timeout 600
```

**验证四个服务都健康：**



```
cd /opt/myblog/deploy && docker compose ps
```

四个容器都显示 `Up`、`healthy`（MySQL/ES 健康检查通过）即成功。

再单独验证后端接口通了：



```
curl -s http://127.0.0.1:8080/api/health
```

> 没有 /api/health 就试
> `curl -s http://127.0.0.1:8080/`
> ，或看后端日志：



```
cd /opt/myblog/deploy && docker compose logs --tail 50 backend
```

**常见报错与处理：**



| 现象                    | 原因与处理                                                      |
| --------------------- | ---------------------------------------------------------- |
| `vm.max_map_count` 报错 | 已提前设置，若仍报，执行 `sudo sysctl -w vm.max_map_count=262144`      |
| 拉镜像很慢 / 失败            | 国内网络访问 docker.elastic.co 慢，可重试；1Panel 装好后会在「容器→仓库」里自带加速镜像源 |
| backend 容器反复重启        | 看日志定位；多半是 MySQL/ES 还没就绪，等 30 秒再 `docker compose ps`        |
| 端口冲突                  | `docker compose ps` 看是哪个端口被占，或 `ss -lntp` 排查               |



***

## 第 5 步：上传数据备份并导入（Mac 上操作）

**做什么**：把本机导出的 MySQL、ES、uploads 三件套搬到服务器并灌进去。

**为什么**：本地容器里的数据不会自动上云，必须导出→传输→导入三步。

### 5.1 上传备份（本机执行）

备份已在 `deploy/backup/` 下（已生成）：



```
workbench upload /Users/xiaoyuwang/projects/MyBlog/deploy/backup/blog\\\_db.sql /opt/myblog/deploy/backup/blog\\\_db.sql -f -i i-f8z9zmvj3iow727wjmpd

workbench upload /Users/xiaoyuwang/projects/MyBlog/deploy/backup/es\\\_data.json /opt/myblog/deploy/backup/es\\\_data.json -f -i i-f8z9zmvj3iow727wjmpd

workbench upload /Users/xiaoyuwang/projects/MyBlog/deploy/backup/uploads.tar.gz /opt/myblog/deploy/backup/uploads.tar.gz -f -i i-f8z9zmvj3iow727wjmpd
```

> `-f`
> 表示覆盖同名文件，避免工具停下来问 "是否覆盖"。

### 5.2 导入 MySQL（服务器上执行）



```
cd /opt/myblog/deploy && docker compose exec -T mysql sh -c 'exec mysql -uroot -p"\\\$MYSQL\\\_ROOT\\\_PASSWORD"' < backup/blog\\\_db.sql
```

看到没有报错、回到命令行即成功。验证：



```
cd /opt/myblog/deploy && docker compose exec -T mysql sh -c 'exec mysql -uroot -p"\\\$MYSQL\\\_ROOT\\\_PASSWORD" -e "SHOW TABLES FROM blog\\\_db;"'
```

能看到文章、分类等表名就成功。

### 5.3 导入 Elasticsearch（服务器上执行）

后端 CLI 自带导入命令，一步重建索引 + 灌数据：



```
cd /opt/myblog/deploy && docker compose exec -T backend ./myblog-server --es-import /app/backup/es\\\_data.json
```

> 原理：这条命令会先删掉旧索引再重建（保证索引结构与当前代码一致），然后批量写入。看到
> `Successfully imported`
> 即成功。导入前确认容器里
> `backup/`
> 挂载生效（compose 已配置
> `./backup:/app/backup`
> ）。

### 5.4 解压上传文件（服务器上执行）



```
cd /opt/myblog/deploy && tar -xzf backup/uploads.tar.gz
```

验证：



```
ls /opt/myblog/deploy/uploads && ls /opt/myblog/deploy/uploads/image | head
```

能看到图片文件就成功。



***

## 第 6 步：构建并上传前端（Mac 上操作）

**做什么**：把 Vue3 前端打成静态文件，上传到服务器，供 1Panel 托管。

**为什么**：前端构建产物（dist/）就是纯静态网页，1Panel 用 Nginx 直接托管它，接口请求再反代到后端。

### 6.1 本地构建



```
cd /Users/xiaoyuwang/projects/MyBlog/web

npm install

npm run build
```

构建产物在 `web/dist/`。打包：



```
tar -czf /tmp/web-dist.tar.gz -C /Users/xiaoyuwang/projects/MyBlog/web dist
```

### 6.2 上传并解压



```
workbench upload /tmp/web-dist.tar.gz /opt/web-dist.tar.gz -i i-f8z9zmvj3iow727wjmpd
```



```
workbench exec -i i-f8z9zmvj3iow727wjmpd -c "mkdir -p /opt/myblog/web-dist && tar -xzf /opt/web-dist.tar.gz -C /opt/myblog/web-dist --strip-components=1 && rm -f /opt/web-dist.tar.gz && ls /opt/myblog/web-dist | head" --timeout 60
```

看到 `index.html`、`assets` 等就成功。

> 如果前端请求后端用的不是相对路径（
> `/api`
> ），检查
> `web/.env.production`
> 里的
> `VITE_API_BASE`
> ，确保线上走
> `/api`
> （与后端
> `router_prefix: api`
> 对应）。



***

## 第 7 步：1Panel 建站 + 反向代理 + HTTPS（网页操作）

**做什么**：在 1Panel 里创建网站，把前端目录托管起来，把 `/api` 和 `/uploads` 转发给后端，再申请 HTTPS 证书。

**为什么**：这是 1Panel 的核心价值 ——Nginx 托管静态站 + 反代后端 + 自动续期证书，全部面板化操作。

### 7.1 创建网站



1. 面板左侧「网站」→「创建网站」→ 选「**静态网站**」

2. 域名：填你的主域名（如 `www.你的域名`，加不加 `www` 看你 DNS 里留了哪条；两条都留就建两个站点或做跳转）

3. 主目录：`/opt/myblog/web-dist`

4. 其余默认，提交

### 7.2 配置反向代理（把接口转发给后端）



1. 网站列表 → 点击刚建的网站 →「配置文件」（Nginx 配置文件）

2. 在 `server { }` 块内（`location /` 之后）添加两段：



```
location /api {

\&#x20;   proxy\\\_pass http://127.0.0.1:8080;

\&#x20;   proxy\\\_set\\\_header Host \\\$host;

\&#x20;   proxy\\\_set\\\_header X-Real-IP \\\$remote\\\_addr;

\&#x20;   proxy\\\_set\\\_header X-Forwarded-For \\\$proxy\\\_add\\\_x\\\_forwarded\\\_for;

\&#x20;   proxy\\\_set\\\_header X-Forwarded-Proto \\\$scheme;

}

location /uploads {

\&#x20;   proxy\\\_pass http://127.0.0.1:8080;

\&#x20;   proxy\\\_set\\\_header Host \\\$host;

}
```



1. 保存后点「重载」（或「重启」）让配置生效

2. 上传图片大小限制：在同一配置文件 `http` 或 `server` 块里加一行：



```
client\\\_max\\\_body\\\_size 20m;
```

> 原理说明：前端静态文件由 Nginx 直接返回；
> `/api/**`
> 和
> `/uploads/**`
> 的请求被 Nginx 转发到
> `127.0.0.1:8080`
> 上的后端容器（compose 里 backend 只绑了本机回环地址，外部无法直连，安全）。
> `$host`
> 保证后端收到正确的域名。

### 7.3 申请 HTTPS 证书



1. 网站列表 → 该网站 →「HTTPS」→「申请证书」

2. 类型选 Let's Encrypt，验证方式选「**DNS 验证**」（阿里云域名）

3. DNS 服务商选「阿里云」，填入你的 AccessKey ID / Secret（没有的话在阿里云 RAM 里建一个子账号，只授权 `AliyunDNSFullAccess` 就行；你现有的 Workbench 子账号也可加这个权限）

4. 勾选自动续期，提交

5. 申请成功后，把网站「HTTPS」开关打开

### 7.4 打开强制 HTTPS（可选）

面板网站设置里开启「强制 HTTPS」，让 http 自动跳转 https。



***

## 第 8 步：验收（浏览器 + 终端）



| 检查项  | 方法                                    | 预期                                            |
| ---- | ------------------------------------- | --------------------------------------------- |
| 网页打开 | 浏览器访问 `https://你的域名`                  | 首页正常显示、样式 / 图片齐全                              |
| 接口正常 | 浏览器访问 `https://你的域名/api/article/list` | 返回 JSON 文章列表                                  |
| 搜索可用 | 页面搜索框搜一个文章关键词                         | 能搜出结果（ES 数据已导入）                               |
| 登录可用 | 用管理员账号登录                              | 成功进入后台                                        |
| 图片上传 | 后台传一张图                                | 上传成功，图片地址为 `/uploads/...`                     |
| 备案号  | 页面底部                                  | 若显示占位符，去 `config.yaml.prod` 填真实备案号后重启 backend |

管理员账号：`2652777599@qq.com` / `admin123456`（**上线后立即改密**）。

验证后端日志无异常：



```
cd /opt/myblog/deploy && docker compose logs --tail 50 backend
```



***

## 附录 A：常用命令速查



| 场景       | 命令                                                              |
| -------- | --------------------------------------------------------------- |
| 连服务器（交互） | `workbench connect -i i-f8z9zmvj3iow727wjmpd`                   |
| 查实例      | `workbench list ecs -r cn-heyuan`                               |
| 执行单条命令   | `workbench exec -i i-f8z9zmvj3iow727wjmpd -c "命令" --timeout 60` |
| 上传文件     | `workbench upload 本地文件 服务器路径 -f -i i-f8z9zmvj3iow727wjmpd`      |
| 查看编排状态   | `cd /opt/myblog/deploy && docker compose ps`                    |
| 看后端日志    | `cd /opt/myblog/deploy && docker compose logs -f backend`       |
| 重启后端     | `cd /opt/myblog/deploy && docker compose restart backend`       |
| 停止整栈     | `cd /opt/myblog/deploy && docker compose down`                  |
| 面板信息     | `http://47.120.64.37:30630/0cd9cdd4d7` 用户 `wwy` 密码 `af225bfcb1` |

## 附录 B：关键配置信息（已替你生成，请妥善保存）



| 项目                 | 值                                                                 |
| ------------------ | ----------------------------------------------------------------- |
| MySQL root 密码（服务器） | `42bbc41631556955aaa864ea9a70baef`                                |
| 本地数据库密码（macOS 容器）  | `root`                                                            |
| 管理员账号              | `2652777599@qq.com` / `admin123456`                               |
| 配置文件位置             | `/opt/myblog/deploy/.env` 与 `/opt/myblog/deploy/config.yaml.prod` |

## 附录 C：常见问题

**Q：命令粘贴后报&#x20;**`syntax error near unexpected token`**？**

A：命令里被混入了反斜杠（如 `\$(curl`）。从本文复制后，先检查有没有 `\$`、`\_`、`\~`；也可以手动重新输入一遍 `$(...)` 部分。

**Q：域名解析改完多久生效？**

A：一般 1\~10 分钟；用 `dig +short 你的域名` 确认输出的是 `47.120.64.37`。

**Q：HTTPS 证书申请失败（DNS 验证）？**

A：确认 (1) DNS 解析已生效到新 IP；(2) RAM 子账号有 `AliyunDNSFullAccess` 权限；(3) 申请时域名和解析记录在同一个阿里云账号下。

**Q：网页能开但接口 502？**

A：多半是反代没生效或后端没起来。依次检查：`docker compose ps` 是否 healthy → `curl -s http://127.0.0.1:8080/api/article/list` 是否返回 JSON → Nginx 配置是否已重载。

**Q：想用 1Panel 管理编排？**

A：目前容器是用 docker compose 命令行启动的，1Panel 的「容器→编排」里看不到它属于正常现象（不影响使用）。想纳入 1Panel 管理，可先把 compose 停掉，再在 1Panel 创建编排粘贴同一个 compose 文件启动。