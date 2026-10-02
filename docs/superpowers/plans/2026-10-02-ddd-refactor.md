# MyBlog Server DDD 重构计划

> **For agentic workers:**
>
>  REQUIRED SUB-SKILL: Use superpowers:executing-plans（或 subagent-driven-development）按阶段执行本计划。本计划为阶段级路线图，每个 Phase 批准后需细化为 bite-size 任务再执行；阶段使用 checkbox（
>
> `- [ ]`
>
> ）语法跟踪。

**Goal:** 将 `server/` 从 gin-vue-admin 风格的三层贫血架构（api/service/model + global 全局变量）重构为 DDD 分层架构：模块之间只依赖抽象（接口），形成高内聚、低耦合、基础设施可替换的代码结构，且对外 HTTP 行为完全不变。

**Architecture:** 采用经典 DDD 四层（interface → application → domain → infrastructure）与 go-zero 式 "按业务模块（限界上下文）分包" 的组织方式结合。domain 层零外部依赖（只依赖标准库与自身实体），所有数据访问与外部系统（MySQL/Redis/ES/ 七牛 / 高德 / 热搜 / 日历）以 Port 接口暴露给 application，由 infrastructure 层实现；依赖方向恒指 domain。全局变量 `global.DB/Redis/ESClient/Log` 全部改为构造注入。

**Tech Stack:** Go 1.22.5（module 名保持 `server`）、Gin、GORM/MySQL、go-redis、go-elasticsearch/v8、七牛 SDK、JWT、robfig/cron。**禁止引入 go-zero 或任何 DDD/DI 框架依赖**（如 go-zero、wire、dig、go-kratos），依赖注入用手写构造组装（bootstrap）。

**Spec:** 本文档即 Spec。依据：`server/` 现有代码（api/service/model/router/middleware/task/utils/initialize/global）、`docs/` 架构图、用户偏好（go-zero 风格、Mermaid 图表、零废话）。

## Global Constraints



* 禁止新增第三方框架依赖（go-zero/wire/dig/kratos 等一律不引入）；仅可使用 Go 标准库与现存依赖。

* HTTP 对外行为完全兼容：路由路径、方法、请求 / 响应 JSON 结构、状态码、错误文案均不得变化。

* module 名保持 `server`；迁移后包路径形如 `server/internal/...`。

* 每个 Phase 结束必须 `go build ./...`、`go vet ./...`、`go test ./...` 全绿，且服务可启动、关键路径可手动验证。

* 前端 `web/`、`deploy/`、脚本不在本次重构范围；只保证后端 API 兼容。

* 领域逻辑与基础设施代码必须分层放置，任何 domain 包不得 import `gorm`、`gin`、`go-elasticsearch` 等基础设施库。

* 现有数据库表结构、ES 索引结构不变。



***

## 1. 现状诊断（As-Is）

### 1.1 现有分层与规模



| 目录                                                    | 职责                                               | 规模 / 说明                                       |
| ----------------------------------------------------- | ------------------------------------------------ | --------------------------------------------- |
| `api/`                                                | Gin Handler（11 个模块）                              | user 435 行、article 252 行等                     |
| `service/`                                            | 业务逻辑（贫血事务脚本）                                     | 约 2400 行，article 451 行、forum 424 行、user 281 行 |
| `model/database`                                      | 贫血 GORM 实体                                       | 16 个表结构                                       |
| `model/request,response,other,appTypes,elasticsearch` | DTO / 枚举 / ES 文档结构                               | —                                             |
| `router/` + `initialize/router.go`                    | 路由注册                                             | public/private/admin 三组 + JWT/Admin 中间件       |
| `middleware/`                                         | jwt、admin、logger、login\_record                   | —                                             |
| `config/ core/ global/`                               | 配置加载、日志、全局变量                                     | `global.DB/Redis/ESClient/Log/Config`         |
| `initialize/`                                         | DB/Redis/ES/Cron/ 路由初始化                          | —                                             |
| `utils/`                                              | jwt、bcrypt、email、upload、pagination、hotSearch 爬虫等 | 通用与业务混杂                                       |
| `task/`                                               | cron 任务（浏览量同步、热搜、日历）                             | 直接调用 service + global                         |

### 1.2 核心耦合点（重构要消灭的对象）



1. **全局变量横穿所有层**：`global.DB`、`global.Redis`、`global.ESClient`、`global.Log`、`global.Config` 被 api/service/task/middleware 直接引用 —— 隐式依赖、无法测试、基础设施无法替换。

2. **贫血模型**：`model/database` 是纯数据 + GORM tag，业务规则（如注册校验、分类 / 标签计数联动、评论树构建）全部散落在 service 事务脚本中。

3. **service 之间横向耦合**：`ServiceGroupApp.XxxService` 互调（user→jwt/gaode、website→hotSearch/calendar、comment→comment\_helps 递归），无显式接口边界。

4. **外部系统无防腐层**：ES client、七牛 SDK、高德 HTTP、热搜爬虫在 service 中被直接操作，替换基础设施 = 改业务代码。

5. **DTO 与技术层捆绑**：request/response 与 Gin、ES typedapi 类型混用，domain 概念未独立。

6. **cron 任务穿透业务**：task 直接 `service.ServiceGroupApp...` + `global.ESClient`。



***

## 2. 目标架构（To-Be）

### 2.1 目录结构（DDD 四层 × 限界上下文分包）



```
server/
├── cmd/server/main.go            # 入口：加载配置 → bootstrap 组装 → 启动 HTTP/Cron
├── internal/
│   ├── config/                   # 配置加载与校验（替代 core/config，注入式，不再有 global.Config）
│   ├── common/                   # 跨 BC 共享的纯工具与 Port 定义（无业务语义）
│   │   ├── jwt/  crypto/  email/  page/  errs/  logger/
│   │   └── storage/              # FileStoragePort（本地/七牛共用接口）
│   ├── bootstrap/                # 基础设施初始化 + 手工 DI 组装（替代 initialize + global）
│   ├── interface/
│   │   ├── http/
│   │   │   ├── router.go         # 路由注册（替代 initialize/router.go）
│   │   │   ├── middleware/       # jwt、admin、logger、login_record（改造为注入式）
│   │   │   └── handler/          # 各 BC 的 Gin Handler（替代 api/，按 BC 分包）
│   │   │       ├── user/  article/  comment/  forum/  image/
│   │   │       ├── website/  advertisement/  friendlink/  feedback/  config/
│   │   └── cron/                 # cron 装配 + 各 BC 定时任务（替代 task/）
│   ├── application/              # 各 BC 应用服务（用例编排/事务/DTO 转换）
│   │   ├── auth/  user/  article/  comment/  forum/
│   │   └── image/  website/  advertisement/  friendlink/  feedback/  config/
│   ├── domain/                   # 领域层：实体、值对象、领域服务、仓储/外部 Port 接口
│   │   ├── auth/  user/  article/  comment/  forum/  image/
│   │   ├── website/  advertisement/  friendlink/  feedback/
│   │   └── shared/               # 跨 BC 共享领域概念（分页、RoleID/Register 等枚举）
│   └── infrastructure/           # Port 实现（可整体替换）
│       ├── mysql/                # GORM 仓储实现（user/article/comment/forum/...）
│       ├── redis/                # 浏览量计数、JWT 会话、验证码
│       ├── es/                   # 文章搜索 + 索引管理（替代 es_index service）
│       ├── storage/              # qiniu、local 上传实现
│       ├── geo/                  # 高德 API 实现（IP 定位/天气）
│       └── hotsearch/            # 百度/知乎/B站/抖音/快手/头条热搜爬虫
└── （旧目录 api/ service/ model/ router/ middleware/ task/ utils/ global/ 在 Phase 5 全部删除）
```

### 2.2 分层依赖规则（依赖倒置核心）



* **domain**：不 import 任何基础设施库（gorm/gin/es/redis/qiniu），只依赖标准库与 `domain/shared`。实体可带行为（充血），规则收敛到实体 / 领域服务。

* **application**：import domain（接口）+ common；编排用例、管理事务边界、做 DTO ↔ 实体转换。不碰 gin/gorm。

* **interface**：import application + common；只做 HTTP 协议适配与参数绑定。

* **infrastructure**：import domain（接口）+ common；实现仓储与外部系统适配。

* **依赖方向**：interface → application → domain ← infrastructure（依赖全部指向 domain）。

### 2.3 限界上下文划分（BC）



| BC            | 来源代码                                                               | 聚合 / 实体                          | 关键领域职责                                          | 依赖的 Port（外）                                                |
| ------------- | ------------------------------------------------------------------ | -------------------------------- | ----------------------------------------------- | ---------------------------------------------------------- |
| auth（支撑）      | service/jwt、middleware/jwt、jwt\_blacklist                          | JwtBlacklist、Token               | 签发 / 校验、黑名单、Redis 会话                            | BlacklistStore（redis）                                      |
| user（核心）      | service/user、model/user、login                                      | User                             | 注册、邮箱登录、找回密码、资料、冻结、登录记录                         | AuthPort、UserRepository、LoginRecordRepository              |
| article（核心）   | service/article\*、article\_category/tag/like、elasticsearch/article | Article、Category、Tag、ArticleLike | CRUD、ES 搜索、浏览量（redis 增量→ES）、分类 / 标签计数事务、置顶 / 点赞 | ArticleRepository、EsArticlePort、ViewCounter（redis）         |
| comment       | service/comment\*、model/comment                                    | Comment（树）                       | 树形加载、级联删除                                       | CommentRepository                                          |
| forum         | service/forum、forum\_post/comment/like                             | ForumPost、ForumComment、ForumLike | 帖子 / 回复 / 点赞                                    | ForumRepository                                            |
| image         | service/image、model/image                                          | Image                            | 上传、分类、存储策略                                      | ImageRepository、FileStoragePort                            |
| website（聚合）   | service/website、service/config、model/website                       | WebsiteConfig                    | 站点配置、热搜 / 日历 / 天气聚合展示                           | HotSearchPort、CalendarPort、GeoPort、WebsiteConfigRepository |
| advertisement | service/advertisement                                              | Advertisement                    | 广告位管理                                           | AdvertisementRepository                                    |
| friendlink    | service/friend\_link、footer\_link                                  | FriendLink                       | 友链 / 页脚链接                                       | FriendLinkRepository                                       |
| feedback      | service/feedback                                                   | Feedback                         | 留言反馈                                            | FeedbackRepository                                         |
| es-index（支撑）  | service/es\_index                                                  | EsIndex                          | ES 索引生命周期管理                                     | EsAdminPort                                                |

> 说明：forum_comment/forum_like 归属 forum BC；comment BC 只处理文章评论。calendar 由 6tail lunar 库支撑，以 CalendarPort 暴露。



***

## 3. Port / Repository 接口清单（依赖抽象的核心交付物）

每个 BC 的 domain 目录内定义如下接口（示意签名，最终以代码为准）：



```
// domain/article/repository.go
type ArticleRepository interface {
    FindByID(ctx context.Context, id string) (*Article, error)
    ExistsByTitle(ctx context.Context, title string) (bool, error)
    Create(ctx context.Context, a *Article) error
    Update(ctx context.Context, a *Article) error
    DeleteByIDs(ctx context.Context, ids []string) error
    SetTop(ctx context.Context, id string, top bool) error
    UpdateCategoryCount(ctx context.Context, tx Tx, oldCat, newCat string) error
    UpdateTagsCount(ctx context.Context, tx Tx, oldTags, newTags []string) error
    Page(ctx context.Context, cond PageCond) ([]Article, int64, error)
}

type EsArticlePort interface {        // 外部系统：Elasticsearch
    Search(ctx context.Context, req SearchCond) ([]ArticleDoc, int64, error)
    IndexDoc(ctx context.Context, doc ArticleDoc) error
    DeleteDoc(ctx context.Context, id string) error
    UpdateViews(ctx context.Context, id string, delta int) error
}

type ViewCounter interface {          // 外部系统：Redis 浏览量增量计数
    Add(ctx context.Context, id string, n int) error
    Snapshot(ctx context.Context) (map[string]int, error)
    Clear(ctx context.Context) error
}
```



```
// domain/user/repository.go
type UserRepository interface {
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByUUID(ctx context.Context, uuid string) (*User, error)
    Create(ctx context.Context, u *User) error
    Update(ctx context.Context, u *User) error
}

// domain/auth/port.go
type AuthPort interface {
    GenerateToken(ctx context.Context, u *User) (string, error)
    Parse(ctx context.Context, token string) (Claims, error)
    Blacklist(ctx context.Context, token string) error
    IsBlacklisted(ctx context.Context, token string) (bool, error)
}
```



```
// domain/website/port.go（聚合展示域依赖的外部 Port）
type HotSearchPort interface { GetBySource(ctx context.Context, source string) ([]HotItem, error) }
type CalendarPort  interface { GetByDate(ctx context.Context, date string) (CalendarInfo, error) }
type GeoPort       interface { LocationByIP(ctx context.Context, ip string) (GeoInfo, error)
                               WeatherByAdcode(ctx context.Context, adcode string) (Weather, error) }
type FileStoragePort interface { Upload(ctx context.Context, f File, dir string) (string, error)
                                 Delete(ctx context.Context, path string) error }
```

**接口归属规则**：Repository / 领域服务接口放对应 BC 的 `domain/`；纯外部系统（ES / 七牛 / 高德 / 热搜 / 日历 /redis 会话）以 Port 形式放消费方 BC 的 domain，实现放 `infrastructure/`。



***

## 4. 迁移路线图（绞杀者模式：新结构与旧结构并存，逐 BC 替换）

> 每个 Phase 独立可交付、可回滚（git 分支 + 保留旧目录至 Phase 5）。

### Phase 0：基座化（无业务行为变化）

**范围**：目录骨架、配置注入、基础设施实例化、DI 组装、路由骨架。



* [x] 创建 `cmd/server/main.go`、`internal/config`、`internal/bootstrap`、`internal/common/{logger,jwt,crypto,email,page,errs}` 骨架

* [x] `internal/bootstrap` 完成 DB/Redis/ES 实例化并持有句柄（替代 `global.DB/Redis/ESClient` 的写入点；旧代码过渡期仍写 global，由 bootstrap 赋值）

* [x] 配置从 `global.Config` 改为 `config.Config` 实例经构造传递；`global.Config` 仅保留兼容层

* [x] 建立 `internal/interface/http/router.go` 新骨架，按 BC 注册空路由组（不挂业务），旧路由保持原样

* [x] `go build ./... && go vet ./... && go test ./...` 全绿；服务可启动

* [x] commit：`refactor(ddd): bootstrap skeleton with constructor injection`

**验收**：应用启动、访问既有接口行为与重构前一致；全局变量写入点收敛到 bootstrap。

### Phase 1：auth + user BC（核心域第一刀，依赖最独立）

**范围**：auth（JWT / 黑名单）与 user 完整迁移。



* [x] `domain/auth`：JwtBlacklist 实体、AuthPort 接口

* [x] `domain/user`：User 实体（含 `Freeze`、`RoleID` 校验等行为）、UserRepository、LoginRecordRepository 接口

* [x] `infrastructure/redis`：JWT 黑名单 / 会话实现；`infrastructure/mysql/user.go`：User 仓储实现

* [x] `application/auth`：Token 用例；`application/user`：Register/EmailLogin/ForgotPassword/UserCard/LoginLog/Freeze 用例（事务边界在此层）

* [x] `interface/http/handler/user`：迁移 user 全部路由；`interface/http/middleware`：JWT/Admin 改为从 AuthPort 注入

* [x] 旧 `api/user.go`、`service/user.go`、`service/jwt.go` 中 user/jwt 部分停止被引用（先不删除，防回滚）

* [x] 单测：User 实体规则、auth token 黑名单流程（stub 仓储）；gin 冒烟测试覆盖 user 接口

* [x] commit：`refactor(ddd): migrate auth & user bounded contexts`

**验收**：user 全部接口行为不变（注册 / 登录 / 登出 / 个人卡片 / 冻结）；`go test ./...` 通过。

### Phase 2：article BC（最复杂，含 ES 与事务）

**范围**：文章 CRUD、搜索、浏览量、分类 / 标签计数、置顶、点赞。



* [x] `domain/article`：Article/Category/Tag/ArticleLike 实体、ArticleRepository、EsArticlePort、ViewCounter 接口

* [x] `infrastructure/mysql/article.go`：仓储实现（含 Create/Delete 时分类与标签计数的**事务联动**，事务边界暴露给 application）；`infrastructure/es`：搜索 / 索引 / 浏览量回写实现；`infrastructure/redis`：浏览量计数实现

* [x] `application/article`：Create（事务）/Update/Delete/Top/Search/Get/View 用例

* [x] `interface/http/handler/article`：迁移 article 路由；`interface/cron`：浏览量同步任务改注入 article 用例

* [x] 单测：article 实体规则、创建事务联动（stub/mock 仓储）、ES 搜索查询构建（不连真实 ES，验证 query 结构）

* [x] commit：`refactor(ddd): migrate article bounded context with es & view counter`

**验收**：文章创建→分类 / 标签计数一致、搜索行为不变、浏览量 redis→ES 回写任务正常；旧 `service/article*.go` 停止被引用。

### Phase 3：comment + forum BC



* [x] `domain/comment`：Comment 树实体（LoadChildren / 级联删除行为收敛）、CommentRepository

* [x] `domain/forum`：ForumPost/ForumComment/ForumLike 实体、ForumRepository

* [x] `infrastructure/mysql`：comment/forum 仓储；`application/comment|forum` 用例（树加载、级联删除事务、点赞幂等）

* [x] handler 与 cron 迁移；单测覆盖树构建与级联删除

* [x] commit：`refactor(ddd): migrate comment & forum bounded contexts`

### Phase 4：image /website/advertisement /friendlink/feedback + 外部防腐层



* [x] `infrastructure/storage`：FileStoragePort 实现（local/qiniu，替代 utils/upload）；`domain/image`、application/handler 迁移

* [x] `infrastructure/geo`（高德）、`infrastructure/hotsearch`（多平台爬虫）、CalendarPort 实现；`domain/website`（配置聚合 + 消费以上 Port）、application/handler 迁移；`service/gaode.go`、`service/calendar.go`、`service/hot_search.go` 停用

* [x] advertisement、friendlink、feedback、config BC 迁移（轻量）

* [x] `internal/interface/cron`：热搜、日历定时任务改注入

* [x] commit：`refactor(ddd): migrate support contexts & external adapters`

**验收**：应用启动、访问既有接口行为与重构前一致；全局变量写入点收敛到 bootstrap。

### Phase 1：auth + user BC（核心域第一刀，依赖最独立）

**范围**：auth（JWT / 黑名单）与 user 完整迁移。



* [x] `domain/auth`：JwtBlacklist 实体、AuthPort 接口

* [x] `domain/user`：User 实体（含 `Freeze`、`RoleID` 校验等行为）、UserRepository、LoginRecordRepository 接口

* [x] `infrastructure/redis`：JWT 黑名单 / 会话实现；`infrastructure/mysql/user.go`：User 仓储实现

* [x] `application/auth`：Token 用例；`application/user`：Register/EmailLogin/ForgotPassword/UserCard/LoginLog/Freeze 用例（事务边界在此层）

* [x] `interface/http/handler/user`：迁移 user 全部路由；`interface/http/middleware`：JWT/Admin 改为从 AuthPort 注入

* [x] 旧 `api/user.go`、`service/user.go`、`service/jwt.go` 中 user/jwt 部分停止被引用（先不删除，防回滚）

* [x] 单测：User 实体规则、auth token 黑名单流程（stub 仓储）；gin 冒烟测试覆盖 user 接口

* [x] commit：`refactor(ddd): migrate auth & user bounded contexts`

**验收**：user 全部接口行为不变（注册 / 登录 / 登出 / 个人卡片 / 冻结）；`go test ./...` 通过。

### Phase 2：article BC（最复杂，含 ES 与事务）

**范围**：文章 CRUD、搜索、浏览量、分类 / 标签计数、置顶、点赞。



* [x] `domain/article`：Article/Category/Tag/ArticleLike 实体、ArticleRepository、EsArticlePort、ViewCounter 接口

* [x] `infrastructure/mysql/article.go`：仓储实现（含 Create/Delete 时分类与标签计数的**事务联动**，事务边界暴露给 application）；`infrastructure/es`：搜索 / 索引 / 浏览量回写实现；`infrastructure/redis`：浏览量计数实现

* [x] `application/article`：Create（事务）/Update/Delete/Top/Search/Get/View 用例

* [x] `interface/http/handler/article`：迁移 article 路由；`interface/cron`：浏览量同步任务改注入 article 用例

* [x] 单测：article 实体规则、创建事务联动（stub/mock 仓储）、ES 搜索查询构建（不连真实 ES，验证 query 结构）

* [x] commit：`refactor(ddd): migrate article bounded context with es & view counter`

**验收**：文章创建→分类 / 标签计数一致、搜索行为不变、浏览量 redis→ES 回写任务正常；旧 `service/article*.go` 停止被引用。

### Phase 3：comment + forum BC



* [x] `domain/comment`：Comment 树实体（LoadChildren / 级联删除行为收敛）、CommentRepository

* [x] `domain/forum`：ForumPost/ForumComment/ForumLike 实体、ForumRepository

* [x] `infrastructure/mysql`：comment/forum 仓储；`application/comment|forum` 用例（树加载、级联删除事务、点赞幂等）

* [x] handler 与 cron 迁移；单测覆盖树构建与级联删除

* [x] commit：`refactor(ddd): migrate comment & forum bounded contexts`

### Phase 4：image /website/advertisement /friendlink/feedback + 外部防腐层



* [x] `infrastructure/storage`：FileStoragePort 实现（local/qiniu，替代 utils/upload）；`domain/image`、application/handler 迁移

* [x] `infrastructure/geo`（高德）、`infrastructure/hotsearch`（多平台爬虫）、CalendarPort 实现；`domain/website`（配置聚合 + 消费以上 Port）、application/handler 迁移；`service/gaode.go`、`service/calendar.go`、`service/hot_search.go` 停用

* [x] advertisement、friendlink、feedback、config BC 迁移（轻量）

* [x] `internal/interface/cron`：热搜、日历定时任务改注入

* [x] commit：`refactor(ddd): migrate support contexts & external adapters`

### Phase 5：收尾清理与全量回归



* [ ] `middleware/logger.go`、`login_record` 迁移收口；`common/logger` 注入替换 `global.Log`

* [ ] 删除旧目录：`api/ service/ model/ router/ middleware/ task/ utils/ global/`（`utils` 中已迁入 common/infrastructure 的部分先核对）

* [ ] `main.go` 迁至 `cmd/server/main.go`，根目录不再有业务代码

* [ ] 全量回归：`go build ./... && go vet ./... && go test ./...`；HTTP 冒烟脚本覆盖全部既有接口路径（对照重构前抓取的响应快照）

* [ ] 文档：更新 `Readme.md`/`record.md` 目录结构说明

* [ ] commit：`refactor(ddd): remove legacy layers & finalize`

**验收**：仓库中不再有 global 全局变量与旧三层目录；行为与重构前一致。



***

## 5. 测试策略



| 层              | 测试方式                            | 说明                               |
| -------------- | ------------------------------- | -------------------------------- |
| domain         | 纯单元测试                           | 实体规则、领域服务，无外部依赖                  |
| application    | 单元测试（stub 仓储 / Port）            | 手写 stub 实现接口，不引入 mock 框架（零新依赖约束） |
| infrastructure | 集成测试（真实 MySQL/Redis，可选 ES 本地容器） | 仓储 CRUD、事务联动                     |
| interface      | gin httptest 冒烟                 | 每 Phase 后跑既有接口路径，与重构前响应对照        |
| 回归             | 脚本快照对比                          | Phase 0 时录制全部接口响应样例，逐 Phase 复核   |



***

## 6. 风险与对策



| 风险                                    | 对策                                                            |
| ------------------------------------- | ------------------------------------------------------------- |
| global 引用量大，一次性替换风险高                  | Phase 0 先建注入骨架，旧代码过渡期仍走 global，各 Phase 逐个 BC 消灭引用，Phase 5 才删除 |
| ES 查询 /script 逻辑复杂，抽取 Port 时行为漂移      | EsArticlePort 方法逐一对应现有方法，用重构前录制的查询快照做对照测试                     |
| article 创建 / 删除的分类标签计数事务、comment 级联删除 | 事务边界放 application 用例；仓储提供 tx 语义，回归单测锁定一致性                     |
| cron 任务直连 service/global              | cron handler 统一改注入 application 用例，Phase 2/4 分步迁移              |
| DDD 过度设计（为实体加行为导致行为变化）                | 只迁移规则归属，不改变既有判定逻辑；实体行为先做 "纯搬运"，评审后再谈建模优化                      |
| 长周期多分支回滚复杂                            | 每 Phase 独立 commit + 保留旧目录至 Phase 5，任何阶段可整体回退                  |



***

## 7. 里程碑



| 阶段      | 交付物                      | 预计范围         |
| ------- | ------------------------ | ------------ |
| Phase 0 | 注入骨架 + 目录基线              | 纯结构          |
| Phase 1 | auth + user BC           | 约 3\~4 个仓库文件 |
| Phase 2 | article BC（ES/redis/ 事务） | 工作量最大        |
| Phase 3 | comment + forum BC       | 树 / 级联       |
| Phase 4 | 支撑域 + 外部防腐层              | 多为搬运         |
| Phase 5 | 清理旧层 + 全量回归              | 删代码          |



***

## 8. 执行方式（批准后二选一）

**1. Subagent-Driven（推荐）** — 每个 Phase 派发独立子代理执行，阶段间人工评审，迭代快、上下文干净。

**2. Inline Execution** — 当前会话按 executing-plans 分批执行，带检查点。

> 本计划为阶段级路线图；选定执行方式后，将按 Phase 细化为 bite-size 任务（含测试与提交步骤）再实施。