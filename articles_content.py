# -*- coding: utf-8 -*-
"""MyBlog 重新 mock 的真实可读文章内容（15 篇）"""

ARTICLES = [
    {
        "title": "Go 语言 GORM 实战：从入门到优雅操作 MySQL",
        "category": "后端开发",
        "tags": ["Go", "GORM", "MySQL", "后端", "数据库"],
        "abstract": "系统讲解 Go 生态中最流行的 ORM 框架 GORM，覆盖连接池配置、模型定义、事务与关联查询等核心用法，并结合真实项目代码演示如何写出优雅、可维护的数据库操作。",
        "content": """# Go 语言 GORM 实战：从入门到优雅操作 MySQL

## 为什么选择 GORM

在 Go 生态中，直接使用 `database/sql` 手写 SQL 虽然灵活，但业务复杂之后样板代码会急剧膨胀：字段扫描、类型转换、错误处理、事务管理，每一项都要重复编写。GORM 作为目前社区最流行的 ORM 框架，提供了链式查询、自动迁移、钩子函数、软删除等一揽子能力，让开发者把精力放回业务本身。

GORM 的核心设计理念是"约定优于配置"：表名、主键、时间字段都有合理的默认规则，只要你的模型定义符合约定，几乎不需要额外配置就能跑起来。

## 快速入门

### 连接数据库与连接池

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
    dsn := "root:root@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Warn), // 生产环境建议 Warn 级别
    })
    if err != nil {
        return nil, err
    }

    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(100)      // 最大连接数
    sqlDB.SetMaxIdleConns(10)       // 最大空闲连接数
    sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间
    return db, nil
}
```

连接池参数需要结合业务吞吐量反复压测调整，这里给出的是常见起点：并发写入密集的服务可以适当调大 `MaxOpenConns`，但不要让连接数超过数据库的 `max_connections` 上限。

### 模型定义

```go
type Article struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"size:200;not null;index" json:"title"`
    Content   string    `gorm:"type:longtext" json:"content"`
    Category  string    `gorm:"size:50;index" json:"category"`
    Tags      string    `gorm:"size:500" json:"tags"`
    AuthorID  uint      `json:"author_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

几个要点：

- `DeletedAt` 字段会自动启用**软删除**：删除操作只写 `deleted_at`，查询默认带上 `deleted_at IS NULL` 条件，数据不会真正丢失，也方便误删恢复。
- 常用的查询字段加上 `index` 标签，避免全表扫描。
- `parseTime=True` 是 MySQL 驱动连接串里最容易遗漏的参数，没有它 `time.Time` 字段会解析失败。

### 基础 CRUD

```go
// 创建
article := Article{Title: "GORM 实战", Category: "后端开发"}
result := db.Create(&article)
if result.Error != nil { /* 处理 */ }
fmt.Println(article.ID) // 主键回填

// 查询单条
var a Article
db.First(&a, 10)                    // 主键查询
db.Where("title LIKE ?", "%GORM%").First(&a)

// 列表 + 分页
var list []Article
db.Where("category = ?", "后端开发").
    Order("created_at DESC").
    Limit(10).Offset(20).
    Find(&list)

// 更新
db.Model(&a).Update("title", "新的标题")
db.Model(&a).Updates(map[string]any{"title": "t", "category": "c"})

// 删除（软删除）
db.Delete(&a)
```

## 事务：保证数据一致性

电商下单、发帖计数这类场景必须使用事务：

```go
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&post).Error; err != nil {
        return err // 自动回滚
    }
    // 更新帖子计数等
    if err := tx.Model(&User{}).Where("id = ?", uid).
        UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error; err != nil {
        return err
    }
    return nil
})
```

`db.Transaction` 会自动开启事务、提交或回滚，比手动 `Begin/Commit/Rollback` 更不容易漏掉回滚分支。

## 关联与预加载

GORM 支持 Has One / Has Many / Belongs To / Many To Many 四类关联。查询时用 `Preload` 避免 N+1 问题：

```go
type User struct {
    ID       uint
    Articles []Article `gorm:"foreignKey:AuthorID"`
}

var user User
db.Preload("Articles").First(&user, 1)
// 此时 user.Articles 已填充，只发两条 SQL
```

## 钩子函数

在模型生命周期节点插入业务逻辑，例如创建文章前自动生成摘要：

```go
func (a *Article) BeforeCreate(tx *gorm.DB) error {
    if a.Abstract == "" {
        a.Abstract = utils.MakeSummary(a.Content, 120)
    }
    return nil
}
```

## 小结

GORM 让 Go 的数据库操作变得清晰而克制：连接池管好并发、事务守住一致性、Preload 消除 N+1、软删除保留后悔药。配合 `logger` 观察慢查询，基本可以覆盖 90% 的业务场景。下一篇文章我们会聊缓存层，看看如何在 GORM 之上叠加 Redis，解决读多写少场景的性能问题。
""",
    },
    {
        "title": "Redis 缓存策略设计：穿透、击穿与雪崩的攻防",
        "category": "后端开发",
        "tags": ["Redis", "缓存", "架构", "后端"],
        "abstract": "深入分析缓存穿透、缓存击穿与缓存雪崩三大经典问题的成因与解决方案，给出布隆过滤器、互斥锁、逻辑过期、多级缓存等实用代码示例，帮你构建高可用的缓存层。",
        "content": """# Redis 缓存策略设计：穿透、击穿与雪崩的攻防

## 为什么需要缓存层

对于读多写少的互联网应用，数据库往往是第一个瓶颈。引入 Redis 作为缓存层后，热数据命中内存，数据库压力大幅下降。但缓存不是银弹——设计不当反而会引入新的故障。今天重点拆解缓存世界的三大经典问题：穿透、击穿、雪崩。

```mermaid
flowchart LR
    A[客户端请求] --> B{缓存命中?}
    B -- 是 --> C[直接返回]
    B -- 否 --> D{数据存在?}
    D -- 是 --> E[回源数据库<br/>写回缓存]
    D -- 否 --> F[穿透: 不存在数据<br/>反复打库]
    E --> C
    style F fill:#fde2e2,stroke:#d33
```

## 缓存穿透：查一个不存在的数据

**现象**：大量请求查询一个缓存和数据库中都不存在的 key（比如恶意构造的 `id=-1`），每次请求都穿透缓存直达数据库，数据库被打满。

### 方案一：缓存空值

查询不存在时，把一个短生命周期的空值（如 `"nil"`）写入缓存，TTL 设为 60 秒，后续请求直接命中空值不再打库。

```go
func GetUser(id int64) (*User, error) {
    key := fmt.Sprintf("user:%d", id)
    val, err := rdb.Get(ctx, key).Result()
    if err == redis.Nil {
        // 未命中
        user, dbErr := queryDB(id)
        if dbErr == ErrNotFound {
            rdb.Set(ctx, key, "nil", 60*time.Second) // 缓存空值
            return nil, ErrNotFound
        }
        rdb.Set(ctx, key, user, time.Hour)
        return user, nil
    }
    if val == "nil" {
        return nil, ErrNotFound
    }
    return decode(val)
}
```

### 方案二：布隆过滤器

在缓存之前加一层布隆过滤器：把所有合法 ID 预先生成到位数组，查询前先判断"这个 key 是否存在"。不存在的一定被拦下，存在的小概率误判（会多打一次库），但**不存在的请求不会再打到数据库**。

```go
// 初始化：把所有用户 ID 加入过滤器
bf := bloom.NewBloomFilter(1000000, 5)
for _, id := range allIDs {
    bf.Add(fmt.Sprintf("user:%d", id))
}

func GetUser(id int64) (*User, error) {
    if !bf.Contains(fmt.Sprintf("user:%d", id)) {
        return nil, ErrNotFound // 直接拦截
    }
    // 正常走缓存 -> 数据库
}
```

两种方案可以叠加使用：布隆过滤器拦截大部分非法请求，空值兜住剩余的边界情况。

## 缓存击穿：热点 key 过期瞬间

**现象**：某个**热点 key**（比如双十一的商品详情）在过期的瞬间，大量并发请求同时回源数据库，数据库瞬间被打垮。

### 方案：互斥锁（Mutex）

重建缓存时只允许一个请求回源，其余请求等待或自旋重试：

```go
func GetHot(key string) (string, error) {
    val, err := rdb.Get(ctx, key).Result()
    if err == nil {
        return val, nil
    }
    lockKey := "lock:" + key
    // 抢锁，带过期时间防止死锁
    ok, _ := rdb.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
    if ok {
        defer rdb.Del(ctx, lockKey)
        data := queryDB(key)              // 只有一个请求回源
        rdb.Set(ctx, key, data, time.Hour)
        return data, nil
    }
    time.Sleep(50 * time.Millisecond)
    return GetHot(key) // 自旋重试（加次数上限）
}
```

### 方案：逻辑过期

缓存不设置物理 TTL，而是把过期时间写进 value；读取时发现逻辑过期，则异步重建并更新缓存，请求方仍然返回旧数据。这种方案对"允许短暂脏读"的场景体验最好。

```go
type CacheValue struct {
    Data   string
    Expire int64 // 逻辑过期时间戳
}
// 读取时判断：expire 已过 -> 异步刷新，先返回旧值
```

## 缓存雪崩：大量 key 同时过期

**现象**：缓存中大量 key 在同一时间段集体过期（或 Redis 实例宕机），请求全部压到数据库，引发连锁故障。

### 对策

1. **过期时间加随机抖动**：`TTL = base + rand(0, 300s)`，打散过期时间点；
2. **多级缓存**：Redis 之上再加一层本地缓存（如 Go 的 `go-cache`），Redis 抖动时本地兜底；
3. **Redis 高可用**：主从 + Sentinel 哨兵，避免单点；
4. **服务降级**：缓存不可用时直接返回默认值或走异步队列，而不是让请求打穿数据库。

```mermaid
flowchart LR
    A[请求] --> B[本地缓存 go-cache]
    B -- 未命中 --> C[Redis]
    C -- 未命中 --> D[(MySQL)]
    C -- 不可用 --> E[降级: 返回默认值/熔断]
    D --> F[回填 C 与 B]
```

## 小结

| 问题 | 核心原因 | 推荐方案 |
|------|---------|---------|
| 穿透 | key 不存在 | 空值缓存 / 布隆过滤器 |
| 击穿 | 热点 key 过期瞬间并发回源 | 互斥锁 / 逻辑过期 |
| 雪崩 | 大量 key 同时失效 / 实例宕机 | TTL 抖动 + 多级缓存 + 高可用 |

缓存的本质是"用一致性换性能"，没有银弹。把上面几种策略组合起来，配合监控告警，才能让缓存层成为系统真正的加速器，而不是定时炸弹。
""",
    },
    {
        "title": "Elasticsearch 全文检索实战：博客搜索系统的设计",
        "category": "后端开发",
        "tags": ["Elasticsearch", "搜索", "Go", "架构"],
        "abstract": "以一个真实博客系统为例，讲解 ES 索引设计、分词器选择、中文搜索优化以及数据同步方案，帮助读者快速构建高性能的全文检索服务。",
        "content": """# Elasticsearch 全文检索实战：博客搜索系统的设计

## 为什么博客需要 ES

博客系统的搜索需求很朴素：用户输入关键词，希望毫秒级返回标题、正文里匹配的文章，还能按相关度排序。用 MySQL 的 `LIKE '%keyword%'` 在几万篇文章里做全表扫描，延迟轻松破秒，中文分词更是无从谈起。Elasticsearch 正是为这类"全文检索 + 聚合分析"场景设计的分布式搜索引擎。

## 索引设计

### Mapping：字段与分词器

```json
{
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart",
        "fields": { "keyword": { "type": "keyword" } }
      },
      "content": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "tags": { "type": "keyword" },
      "category": { "type": "keyword" },
      "created_at": { "type": "date", "format": "yyyy-MM-dd HH:mm:ss" },
      "views": { "type": "integer" }
    }
  }
}
```

要点：

- 中文场景必须装 **IK 分词器**（或更细的 smartcn）：`ik_max_word` 用于索引（尽可能切细），`ik_smart` 用于搜索（切最合理的词），召回率和准确率兼顾；
- `title` 增加 `keyword` 子字段，支持精确匹配和排序；
- `tags`/`category` 用 `keyword`，不参与分词，用于过滤（filter）和聚合（aggs）。

### 写入与查询

```go
// 写入（带 id 保证幂等）
doc := map[string]any{
    "title": "GORM 实战", "content": "...",
    "tags": []string{"Go", "GORM"}, "category": "后端开发",
    "created_at": "2026-09-01 10:00:00", "views": 0,
}
res, _ := es.Index("articles", esDocID, doc)

// 查询：多字段匹配 + 高亮
q := map[string]any{
    "query": map[string]any{
        "multi_match": map[string]any{
            "query":  keyword,
            "fields": []string{"title^3", "content"},
        },
    },
    "highlight": map[string]any{
        "fields": map[string]any{"content": map[string]any{}},
    },
}
```

`title^3` 表示标题字段权重是正文的 3 倍——标题命中比正文命中更相关，这是搜索排序里最常用的小技巧。

## 中文搜索的常见坑

1. **停用词与噪音**：IK 默认词库对"的、了、吗"等停用词处理有限，可以在自定义词典中维护业务黑名单；
2. **大小写与全半角**：`normalizer` 统一 lowercase，避免 "GORM" 和 "gorm" 查不到；
3. **搜索结果不刷新**：ES 的写操作默认近实时（refresh 间隔 1s），测试时容易误以为写入失败，读多写少的博客系统可把 `refresh_interval` 调到 30s 减负载。

## 数据同步：从 MySQL 到 ES

博客文章的主数据源是 MySQL，ES 只是检索副本。同步方案常用两条路：

```mermaid
flowchart LR
    A[(MySQL)] -->|1 双写| B[业务代码同步写 ES]
    A -->|2 Binlog| C[Canal 监听] --> D[ES]
    A -->|3 定时全量| E[Job 重建索引]
    B --> D
    D --> F[检索服务]
```

三种方案各有取舍：

- **业务双写**：最简单，但存在 MySQL 成功、ES 失败的不一致窗口，需要补偿任务兜底；
- **Canal 订阅 Binlog**：解耦可靠，适合团队已有中间件体系的场景，引入成本高；
- **定时全量重建**：数据量小（几万篇）时最稳，每晚重建一次索引即可。

## 检索服务的兜底逻辑

ES 不是高可用系统的主链路。检索服务要做好降级：ES 超时或异常时，降级为 MySQL `LIKE` 查询返回结果，保证搜索功能"慢但可用"，而不是直接 5xx。

```go
func Search(keyword string) []Article {
    hits, err := es.Search(keyword)
    if err != nil {
        log.Warn("es search failed, fallback to mysql: %v", err)
        return mysqlLikeSearch(keyword) // 降级
    }
    return hits
}
```

## 小结

一套合格的博客搜索 = 合理 mapping（IK 分词 + 权重）+ 可靠的同步链路 + 优雅的降级。ES 的学习曲线不低，但它换来的检索体验和扩展能力，在文章量上来之后会非常值得。本博客的搜索页就是这套方案的实际落地，欢迎在搜索框试一下中文检索效果。
""",
    },
    {
        "title": "Gin 框架中间件机制深度解析",
        "category": "后端开发",
        "tags": ["Go", "Gin", "中间件", "后端", "JWT", "认证"],
        "abstract": "拆解 Gin 中间件的洋葱模型原理，从源码角度分析其执行顺序与上下文传递机制，并手写 JWT 认证、日志记录、CORS 等生产可用的常用中间件。",
        "content": """# Gin 框架中间件机制深度解析

## 中间件：请求管线的积木

Gin 的中间件是构建 Web 服务的核心抽象：认证、日志、限流、跨域、恢复，这些横切关注点都可以用中间件优雅地插入请求处理链，而不用污染业务 handler。理解它，先从"洋葱模型"说起。

```mermaid
flowchart TD
    R[Request 请求] --> M1[中间件A<br/>before]
    M1 --> M2[中间件B<br/>before]
    M2 --> H[Handler<br/>业务逻辑]
    H --> M2b[中间件B<br/>after]
    M2b --> M1b[中间件A<br/>after]
    M1b --> S[Response 响应]
```

请求像洋葱一样**从外到内**穿过所有中间件，到达 handler 后再**从内到外**返回。中间件在 `c.Next()` 之前的代码是"入"阶段，之后的代码是"出"阶段。

## c.Next() 与 c.Abort() 的源码语义

```go
// Gin 核心循环（简化）
func (c *Context) Next() {
    c.index++
    for c.index < int8(len(c.handlers)) {
        c.handlers[c.index](c)
        c.index++
    }
}
```

- `c.Next()`：调用链中的下一个 handler，返回后继续执行本中间件 `Next()` 之后的代码（出阶段）；
- `c.Abort()`：把 `c.index` 直接跳到末尾，**终止后续 handler 执行**（但已注册的 handler 仍会执行完），常用于认证失败直接返回 401；
- `c.Set() / c.Get()`：在中间件间传递数据（如认证中间件写入 `userID`，业务 handler 读取），这是 Gin 上下文的核心用法。

## 手写一个 JWT 认证中间件

```go
func JWTAuth(secret []byte) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "未登录"})
            return
        }
        claims, err := parseJWT(strings.TrimPrefix(token, "Bearer "), secret)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"code": 401, "msg": "登录已过期"})
            return
        }
        c.Set("user_id", claims.UserID) // 写入上下文
        c.Next()
    }
}

// 路由组中使用
r := gin.Default()
api := r.Group("/api")
api.Use(JWTAuth(jwtSecret))
{
    api.GET("/article/list", ListArticles)   // 受保护
    api.GET("/article/:id", GetArticle)      // 受保护
}
```

这样 `/api` 组下所有路由自动带上认证，业务 handler 里 `c.GetInt64("user_id")` 就能拿到当前用户。

## 日志中间件与性能埋点

```go
func AccessLog() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        latency := time.Since(start)
        log.Printf("[%s] %s %d %s",
            c.Request.Method, c.Request.URL.Path,
            c.Writer.Status(), latency.Round(time.Microsecond))
    }
}
```

## CORS 中间件：前后端分离的必修课

```go
func Cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(204) // 预检请求直接返回
            return
        }
        c.Next()
    }
}
```

## 中间件的执行顺序与依赖

`r.Use(A, B)` 的注册顺序就是执行顺序：A 先入、B 后入、B 先出、A 后出。依赖关系要按此设计——比如"恢复中间件（Recovery）必须注册在最外层，才能兜住内层 panic"。

```go
r.Use(gin.Recovery()) // 最外层：捕获 panic
r.Use(Logger())       // 记录请求
r.Use(Cors())         // 跨域
r.Use(RateLimit(100)) // 限流
```

## 小结

中间件把"横切关注点"从业务代码中彻底剥离。理解洋葱模型和 `Next/Abort` 的语义，你就掌握了 Gin 的骨架。生产实践中建议至少备好：Recovery、Logger、JWT 认证、CORS、限流、请求 ID 追踪这六个中间件。
""",
    },
    {
        "title": "Vue3 + Vite 从零搭建企业级前端工程",
        "category": "前端开发",
        "tags": ["Vue", "Vite", "TypeScript", "前端", "工程化"],
        "abstract": "介绍如何基于 Vue3 Composition API 与 Vite 搭建可维护的前端工程，包含目录规范、Pinia 状态管理、路由权限与代码规范等最佳实践，适合想要建立工程化认知的开发者。",
        "content": """# Vue3 + Vite 从零搭建企业级前端工程

## 为什么是 Vue3 + Vite

Vue3 的 Composition API 让逻辑复用从 mixin 的泥潭中解脱出来；Vite 基于 esbuild 的原生 ESM 构建，冷启动速度比 Webpack 时代快一个数量级。两者组合是当下 Vue 生态新建项目的默认选择。

## 工程初始化与目录规范

```bash
npm create vite@latest my-app -- --template vue-ts
cd my-app && npm install
```

一个可持续演进的目录结构：

```mermaid
flowchart TD
    SRC[src] --> API[api/ 接口层<br/>按模块拆分]
    SRC --> ASSETS[assets/ 静态资源]
    SRC --> COMP[components/ 公共组件]
    SRC --> VIEWS[views/ 页面]
    SRC --> STORES[stores/ Pinia 状态]
    SRC --> ROUTER[router/ 路由与守卫]
    SRC --> UTILS[utils/ 工具函数]
    SRC --> TYPES[types/ 全局类型]
    API --> HTTP[http.ts 请求封装]
```

**规范**：页面组件放 `views`，跨页面复用的放 `components`，单页独有组件就近放在页面目录下；`api` 层统一收口所有请求，业务代码禁止直接写 `axios`。

## 请求封装与拦截器

```ts
// api/http.ts
import axios from 'axios'

const http = axios.create({ baseURL: '/api', timeout: 10_000 })

http.interceptors.request.use(config => {
    const token = localStorage.getItem('token')
    if (token) config.headers.Authorization = `Bearer ${token}`
    return config
})

http.interceptors.response.use(
    res => res.data,
    err => {
        if (err.response?.status === 401) {
            // 统一跳登录
            router.push('/login')
        }
        return Promise.reject(err)
    }
)
```

## Pinia：更现代的状态管理

```ts
// stores/user.ts
export const useUserStore = defineStore('user', () => {
    const userInfo = ref<UserInfo | null>(null)
    const isLogin = computed(() => !!userInfo.value)

    async function fetchUserInfo() {
        userInfo.value = await api.getUserInfo()
    }
    return { userInfo, isLogin, fetchUserInfo }
})
```

Composition 风格（setup store）比 Options 风格更贴合 Vue3 的心智模型，类型推导也更好。

## 路由权限：前置守卫

```ts
// router/index.ts
router.beforeEach((to, from, next) => {
    const userStore = useUserStore()
    if (to.meta.requiresAuth && !userStore.isLogin) {
        next({ path: '/login', query: { redirect: to.fullPath } })
        return
    }
    next()
})
```

把 `meta.requiresAuth` 标在需要登录的路由定义上，权限逻辑集中在守卫里，页面组件保持纯净。

## 代码规范：ESLint + Prettier + Husky

```bash
npm install -D eslint prettier husky lint-staged
```

- ESLint 负责规则（Vue3 推荐配置 + TS）；
- Prettier 负责格式化，两者用 `eslint-config-prettier` 消解冲突；
- Husky + lint-staged 在 `pre-commit` 阶段只检查暂存文件，保持提交干净。

```json
// .lintstagedrc
{ "*.{vue,ts}": ["eslint --fix", "prettier --write"] }
```

## 环境变量与多环境构建

```ts
// .env.development
VITE_API_BASE=/api
// .env.production
VITE_API_BASE=https://api.example.com

// 使用时
const base = import.meta.env.VITE_API_BASE
```

Vite 约定：只有 `VITE_` 前缀的变量会暴露给客户端代码，其余环境变量不会泄露到打包产物。

## 小结

工程化的本质不是工具越多越好，而是**边界清晰 + 约定统一**：目录分层、请求收口、状态集中、权限前置、规范自动化。这套骨架搭好后，新成员加入项目只需要半天就能上手开发，这正是企业级工程的价值所在。
""",
    },
    {
        "title": "现代前端工程化：从 Webpack 到 Vite 的迁移实录",
        "category": "前端开发",
        "tags": ["前端", "Vite", "Webpack", "工程化", "性能优化"],
        "abstract": "对比 Webpack 与 Vite 的构建原理，分享一个中大型项目迁移到 Vite 过程中的问题与解决方案，记录构建速度从分钟级到秒级的真实数据与踩坑清单。",
        "content": """# 现代前端工程化：从 Webpack 到 Vite 的迁移实录

## 背景：为什么必须迁移

我们维护的一个中大型后台项目，Webpack 5 冷启动需要 **38 秒**，热更新单次 2~4 秒，开发体验越来越差。随着代码量增长到 300+ 页面、1000+ 组件，团队成员每天大量时间花在等待编译上。

Vite 的核心思路完全不同：**开发阶段不做打包**，利用浏览器原生 ESM 按需加载，冷启动只编译真正被用到的模块。

```mermaid
flowchart LR
    subgraph Webpack
        A1[全部模块] --> A2[打包成 Bundle] --> A3[浏览器加载]
    end
    subgraph Vite
        B1[浏览器请求模块] --> B2[esbuild 按需预构建] --> B3[直接返回 ESM]
    end
```

## 迁移步骤

### 第一步：基础设施替换

```bash
npm uninstall webpack webpack-cli webpack-dev-server
npm install -D vite @vitejs/plugin-vue
```

### 第二步：配置文件迁移

Webpack 的 resolve.alias 到 Vite：

```ts
// vite.config.ts
export default defineConfig({
    plugins: [vue()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
    },
    server: {
        port: 80,
        proxy: {
            '/api': { target: 'http://localhost:8080', changeOrigin: true },
        },
    },
})
```

### 第三步：环境变量

Webpack 的 `process.env.NODE_ENV` 在 Vite 中是 `import.meta.env`：

```ts
// 迁移前
const api = process.env.VUE_APP_API_BASE
// 迁移后
const api = import.meta.env.VITE_API_BASE
```

**坑**：项目里大量 `process.env` 直引用，迁移时写一个兼容垫片最省事——不过这只是过渡方案，最终还是要改成 `import.meta.env`。

## 踩坑清单

### 动态 import 的变量问题

Vite 不支持完全动态的 `import(variable)`，必须让路径可静态分析：

```ts
// 报错：Vite 无法解析
import(`../views/${pageName}.vue`)
// 正确：用 import.meta.glob 预注册
const modules = import.meta.glob('../views/**/*.vue')
const loader = modules[`../views/${pageName}.vue`]
const Comp = await loader()
```

### 第三方库的 CJS 兼容

部分老库只提供 CommonJS 导出，Vite 预构建（optimizeDeps）会自动处理，但如果库在运行时才 require，需要手动 `optimizeDeps.include`：

```ts
optimizeDeps: {
    include: ['some-legacy-lib'],
}
```

### SVG / 静态资源

Webpack 的 `require('@/assets/x.svg')` 写法需要改成 ESM 导入：

```ts
import icon from '@/assets/logo.svg'   // Vite 直接返回 URL
```

## 迁移后的真实收益

| 指标 | Webpack | Vite | 提升 |
|------|---------|------|------|
| 冷启动 | 38s | 1.8s | **21x** |
| HMR 单次 | 2~4s | <100ms | **30x+** |
| 生产构建 | 96s | 58s | 1.7x |

开发体验的提升是全方位的：改一行代码几乎瞬间可见，配合组件级 HMR，调试 React 组件状态、排查样式问题都变得流畅。

## 小结

Vite 不是银弹：大型单体应用的生产构建仍依赖 Rollup，部分生态工具（如 Storybook 早期版本）兼容性有坑。但如果你的团队受困于 Webpack 的开发编译慢，Vite 的迁移收益几乎立竿见影。记住一个原则——**先让开发链路跑起来，再逐步清理兼容垫片**，迁移风险就会可控。
""",
    },
    {
        "title": "Docker Compose 部署前后端分离项目的完整指南",
        "category": "运维部署",
        "tags": ["Docker", "部署", "Nginx", "运维", "微服务"],
        "abstract": "从零讲解如何用 Docker Compose 编排 MySQL、Redis、Elasticsearch 与前后端服务，包含镜像构建、健康检查、数据卷持久化与 Nginx 网关配置，实现一键启动的现代化部署方案。",
        "content": """# Docker Compose 部署前后端分离项目的完整指南

## 为什么用 Docker Compose

一套典型的前后端分离项目至少包含：MySQL、Redis、Elasticsearch、后端 API、前端静态站点、Nginx 网关。逐个手动部署意味着成百上千条命令、无尽的版本冲突。Docker Compose 用一份 YAML 描述全部服务，一条 `docker compose up -d` 全部拉起，环境一致性也得到保证——**本机跑什么，服务器就跑什么**。

## 服务拓扑

```mermaid
flowchart LR
    U[用户浏览器] --> N[Nginx 网关 :80]
    N -->|/api| B[后端 API :8080]
    N -->|静态资源| F[前端静态文件]
    B --> M[(MySQL :3306)]
    B --> R[(Redis :6379)]
    B --> E[(Elasticsearch :9200)]
```

## 第一步：编写后端 Dockerfile

```dockerfile
# 多阶段构建：编译与运行分离，镜像更小
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/uploads ./uploads
EXPOSE 8080
CMD ["./server"]
```

多阶段构建的好处：最终镜像只保留二进制和必要资源，体积从 1GB+ 降到几十 MB。

## 第二步：编写 docker-compose.yml

```yaml
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: blog_db
    volumes:
      - mysql-data:/var/lib/mysql          # 数据卷持久化
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 5s
      retries: 10
    networks: [app]

  redis:
    image: redis:7-alpine
    volumes:
      - redis-data:/data
    networks: [app]

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.17.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - ES_JAVA_OPTS=-Xms512m -Xmx512m
    volumes:
      - es-data:/usr/share/elasticsearch/data
    networks: [app]

  backend:
    build: ./server
    depends_on:
      mysql: { condition: service_healthy }   # 等 MySQL 就绪再启动
      redis: { condition: service_started }
      elasticsearch: { condition: service_started }
    environment:
      DB_DSN: root:root@tcp(mysql:3306)/blog_db?charset=utf8mb4&parseTime=True
      REDIS_ADDR: redis:6379
      ES_ADDR: http://elasticsearch:9200
    networks: [app]

  nginx:
    image: nginx:1.27-alpine
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf:ro
      - ./web/dist:/usr/share/nginx/html:ro
    ports:
      - "80:80"
    depends_on: [backend]
    networks: [app]

volumes:
  mysql-data:
  redis-data:
  es-data:

networks:
  app:
```

几个容易踩的坑：

1. **容器间通信用服务名**：后端连数据库用 `mysql:3306` 而不是 `127.0.0.1`，容器名即 DNS；
2. **`depends_on` 只保证启动顺序，不保证就绪**：配合 `healthcheck + condition: service_healthy` 才能让后端真正等到 MySQL 可连接；
3. **数据卷持久化**：不挂卷的话，容器一删数据全没，生产事故级别的错误；
4. **ES 内存限制**：容器里 ES 默认按宿主机内存分配堆，必须用 `ES_JAVA_OPTS` 显式限制，否则小内存服务器直接 OOM。

## 第三步：Nginx 网关配置

```nginx
server {
    listen 80;
    server_name example.com;

    # 前端静态资源
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;   # SPA 路由回退
    }

    # 后端 API 反向代理
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 一键启动与日常运维

```bash
docker compose up -d            # 启动全部
docker compose logs -f backend  # 看后端日志
docker compose ps               # 查看健康状态
docker compose down             # 停止（保留数据卷）
docker compose up -d --build    # 改代码后重建
```

## 小结

Docker Compose 把"部署"从一门手艺变成了一份配置。它对中小项目的价值在于：**可复现、可迁移、可回滚**。生产环境建议再叠加一层 Traefik/Nginx 做 TLS 终结，并把镜像推到私有仓库，服务器上只跑 `docker compose pull && up -d`，发布就完成了。
""",
    },
    {
        "title": "MySQL 索引优化实战：慢查询排查与调优",
        "category": "后端开发",
        "tags": ["MySQL", "索引", "性能优化", "数据库", "后端"],
        "abstract": "通过真实慢查询案例，讲解 EXPLAIN 解读、索引选择原则、覆盖索引与最左前缀匹配等优化技巧，附性能对比数据，帮你建立 MySQL 性能优化的完整方法论。",
        "content": """# MySQL 索引优化实战：慢查询排查与调优

## 一个真实的慢查询案例

线上告警：某列表接口 P99 延迟从 80ms 飙到 2.3s。排查发现一条 SQL 每次执行 1.9s，表里只有 120 万行，但每天都在增长：

```sql
SELECT * FROM articles
WHERE category = '后端开发'
  AND created_at > '2026-08-01'
  AND tags LIKE '%GORM%'
ORDER BY created_at DESC
LIMIT 20;
```

## 第一步：打开慢查询日志定位元凶

```sql
-- 查看当前慢查询配置
SHOW VARIABLES LIKE 'slow_query_log%';
SET GLOBAL slow_query_log = ON;
SET GLOBAL long_query_time = 1;   -- 超过 1 秒记录

-- 查看慢日志
SHOW GLOBAL STATUS LIKE 'Slow_queries';
```

生产环境建议把 `long_query_time` 设为 0.5~1s，配合 `pt-query-digest` 定期分析。

## 第二步：EXPLAIN 读懂执行计划

```sql
EXPLAIN SELECT ... /* 上面的 SQL */;
```

关键列解读：

| 列 | 本例值 | 含义 |
|----|--------|------|
| type | ALL | **全表扫描**，问题根源 |
| key | NULL | 没有使用任何索引 |
| rows | 1200000 | 预计扫描 120 万行 |
| Extra | Using where; Using filesort | 还要额外排序 |

`type=ALL + key=NULL` 意味着 MySQL 逐行扫了整张表，慢是必然的。

## 第三步：设计索引

### 复合索引与最左前缀

```sql
ALTER TABLE articles ADD INDEX idx_cat_time (category, created_at);
```

`category` 等值匹配 + `created_at` 范围排序，恰好构成最左前缀。执行计划立刻变为 `type=ref, key=idx_cat_time, rows≈20000`，再配合 filesort 消除（索引本身有序，`created_at DESC` 直接倒序扫索引）：

```sql
EXPLAIN SELECT ... -- 现在 type=ref, key=idx_cat_time, Using index condition
```

### 覆盖索引：榨干最后一滴性能

```sql
-- 如果只需要 id 和 title
ALTER TABLE articles ADD INDEX idx_cat_time_title (category, created_at, title);

SELECT id, title FROM articles
WHERE category = '后端开发' ORDER BY created_at DESC LIMIT 20;
-- Extra: Using index   ← 完全不需要回表，理论最快
```

覆盖索引让查询的所有字段都包含在索引里，MySQL 无需回聚簇索引取数据，这是查询优化的天花板。

### 为什么 LIKE 没有用到索引

`tags LIKE '%GORM%'` 前导通配符导致索引失效（无法利用 B+ 树的有序性）。解决：拆出 `article_tags` 关联表做等值匹配，或引入全文索引（`FULLTEXT`）支持中文分词检索。这也是博客搜索最终选择 Elasticsearch 的原因之一。

## 索引选择的通用原则

1. **选择性优先**：`SELECT COUNT(DISTINCT col)/COUNT(*)` 越大越好，选择性低于 10% 的字段别单建索引；
2. **区分读写**：高并发写入表少建索引（每次写都要维护索引树），读多写少表大胆建；
3. **覆盖即最优**：能用覆盖索引就不要回表；
4. **控制冗余**：`(a, b)` 索引已经覆盖 `(a)` 的场景，别再单建 `a` 索引；
5. **警惕隐式转换**：`WHERE phone = 13800000000`（字段是 varchar）会导致索引失效，类型必须一致。

## 验证优化效果

| 阶段 | 扫描行数 | 延迟 |
|------|---------|------|
| 优化前（全表扫描） | 1,200,000 | 1.9s |
| 加复合索引后 | 20,000 | 45ms |
| 覆盖索引后 | 20,000（不回表） | 22ms |

## 小结

MySQL 优化是"证据驱动"的：**慢日志找问题 → EXPLAIN 读执行计划 → 设计索引 → 压测验证**。四步循环下来，绝大多数慢查询都能在半小时内解决。记住：索引不是越多越好，它是空间换时间的取舍，建之前先问自己"这个查询真的高频吗"。
""",
    },
    {
        "title": "Nginx 反向代理与负载均衡配置详解",
        "category": "运维部署",
        "tags": ["Nginx", "负载均衡", "运维", "架构", "网络安全"],
        "abstract": "讲解 Nginx 作为反向代理的核心配置，包括静态资源缓存、gzip 压缩、HTTPS 证书配置与多节点负载均衡策略，并给出生产环境的安全加固清单。",
        "content": """# Nginx 反向代理与负载均衡配置详解

## Nginx 在架构中的位置

Nginx 是互联网架构里的"守门人"：对外暴露 80/443 端口，把请求按规则分发到后端的业务服务。相比让后端直接对公网，Nginx 带来的收益是全方位的——统一入口、静态资源高效服务、负载均衡、TLS 终结、安全防护。

```mermaid
flowchart LR
    U[用户] --> N[Nginx :443 网关]
    N --> S1[后端节点 1 :8080]
    N --> S2[后端节点 2 :8080]
    N --> S3[后端节点 3 :8080]
    N -->|静态资源| F[磁盘/缓存]
```

## 基础反向代理配置

```nginx
upstream backend_servers {
    # 负载均衡策略默认是轮询（round-robin）
    server 10.0.0.11:8080 weight=3;   # weight 越高权重越大
    server 10.0.0.12:8080 weight=1;
    server 10.0.0.13:8080 backup;      # 备份节点，前两个挂了才启用
}

server {
    listen 80;
    server_name example.com;

    location /api/ {
        proxy_pass http://backend_servers;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_connect_timeout 5s;
        proxy_read_timeout 30s;
    }
}
```

## 负载均衡策略怎么选

| 策略 | 指令 | 适用场景 |
|------|------|---------|
| 轮询 | （默认） | 后端无状态、性能均匀 |
| 加权轮询 | `weight=3` | 机器性能有差异 |
| IP 哈希 | `ip_hash` | 需要会话粘滞（session 本地化） |
| 最少连接 | `least_conn` | 请求处理时长差异大 |

无状态后端（JWT 认证的 API）优先用轮询/加权轮询；有状态的 session 场景才考虑 `ip_hash`，更现代的方案是彻底无状态化，把粘滞需求干掉。

## 静态资源与缓存

```nginx
location ~* \\.(js|css|png|jpg|svg|woff2)$ {
    root /usr/share/nginx/html;
    expires 30d;                 # 强缓存 30 天
    add_header Cache-Control "public, immutable";
    access_log off;
}
```

配合前端构建时给文件名加 hash（Vite 默认行为），静态资源可以放心缓存 30 天，用户二次访问几乎零等待。

## gzip 压缩：4 行配置省一半流量

```nginx
gzip on;
gzip_comp_level 6;
gzip_min_length 1k;
gzip_types text/plain text/css application/json application/javascript image/svg+xml;
```

JSON API 和 JS/CSS 是压缩收益最大的对象，压缩后传输体积通常减少 60%~80%。现代浏览器都支持 Brotli，也可以把 `gzip` 换成 `brotli` 模块获得更高压缩率。

## HTTPS：TLS 终结与自动跳转

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    ssl_certificate     /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / { proxy_pass http://backend_servers; }
}
# HTTP 全部 301 到 HTTPS
server { listen 80; server_name example.com; return 301 https://$host$request_uri; }
```

证书用 Let's Encrypt 免费签发，`certbot --nginx` 一条命令完成签发与自动续期。`ssl_protocols` 只保留 TLSv1.2+，把 TLS1.0/1.1 彻底关掉。

## 安全加固清单

```nginx
# 防 SQL 注入 / XSS 的关键响应头
add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-XSS-Protection "1; mode=block" always;

# 隐藏版本号
server_tokens off;

# 简单限流：每 IP 每秒 10 个请求
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
location /api/ {
    limit_req zone=api_limit burst=20 nodelay;
    proxy_pass http://backend_servers;
}
```

`limit_req` 是防止接口被刷的性价比之王，10 行配置就能拦住绝大多数脚本攻击。

## 小结

Nginx 的配置哲学是"声明式"：每个指令解决一个明确问题。把**反向代理 → 静态缓存 → gzip → HTTPS → 限流加固**这条链路走通，一个可用性 99.9% 的网关就搭好了。排障时记住三板斧：`nginx -t` 校验配置、`nginx -s reload` 平滑重载、`access.log` 里找 5xx 和 499。
""",
    },
    {
        "title": "JWT 双 Token 认证机制：设计与 Go 实现",
        "category": "后端开发",
        "tags": ["JWT", "认证", "安全", "Go", "微服务"],
        "abstract": "从安全角度设计 Access Token + Refresh Token 双令牌机制，讲解 Redis 黑名单、多点登录控制与令牌轮换策略，附完整的 Go 实现代码。",
        "content": """# JWT 双 Token 认证机制：设计与 Go 实现

## 单 Token 的隐患

最简单的 JWT 方案：登录后发一个 Access Token，有效期 2 小时，过期就重新登录。问题在于——**JWT 一旦签发就无法在服务端主动作废**（无状态）。用户改密码、封禁账号、设备被盗，旧 Token 依然有效，直到自然过期。这就是单 Token 方案最大的安全缺口。

## 双 Token 机制

```mermaid
sequenceDiagram
    participant C as 客户端
    participant B as 后端
    participant R as Redis

    C->>B: 用户名/密码
    B->>R: 存储 RefreshToken(用户ID, 过期7天)
    B-->>C: 返回 AccessToken(2h) + RefreshToken(7d)
    C->>B: 请求(携带 AccessToken)
    B-->>C: 校验通过，返回数据
    C->>B: AccessToken 过期(401)
    C->>B: 提交 RefreshToken
    B->>R: 校验 RefreshToken 是否存在且未过期
    B-->>C: 签发新的 AccessToken(+新 RefreshToken)
```

核心思路：**Access Token 短命（分钟~小时级），Refresh Token 长命（天级）但必须可作废**。把 Refresh Token 存进 Redis，服务端就能精确控制"这个会话还能不能续期"。

## Go 实现

### 签发 Token

```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateTokenPair(userID uint, username string) (access, refresh string, err error) {
    now := time.Now()
    // Access Token：有效期 2 小时
    access, err = jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
        UserID: userID, Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(now),
        },
    }).SignedString(secret)
    if err != nil {
        return "", "", err
    }
    // Refresh Token：有效期 7 天，携带类型声明
    refresh, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "type":    "refresh",
        "exp":     now.Add(7 * 24 * time.Hour).Unix(),
    }).SignedString(secret)
    return
}
```

### 校验与续期

```go
func RefreshToken(oldRefresh string) (newAccess, newRefresh string, err error) {
    claims, err := parseRefresh(oldRefresh)
    if err != nil {
        return "", "", err
    }
    // 关键：检查 Redis 中是否存在该会话
    key := fmt.Sprintf("refresh:%s", oldRefresh)
    if rdb.Get(ctx, key).Err() == redis.Nil {
        return "", "", ErrSessionExpired // 已被注销/轮换
    }
    // 轮换：删除旧 Refresh Token，签发新对
    rdb.Del(ctx, key)
    return GenerateTokenPair(claims.UserID, claims.Username)
}
```

## 安全要点

### 1. 登录/登出/改密主动作废

```go
// 登出：删掉 Redis 里的 refresh token，Access Token 等它自然过期
rdb.Del(ctx, "refresh:"+refreshToken)

// 改密码：把该用户所有 refresh token 清掉（强制重新登录）
rdb.Del(ctx, "user_sessions:"+userID) // 配合 refresh 键带用户前缀
```

### 2. 多点登录控制

- **限制单端**：登录时把该用户旧的 Refresh Token 全部删除，新 Token 生效，旧设备下次续期直接失败；
- **设备列表管理**：每个会话分配 `session_id`，后台"下线指定设备"只需删对应键。

### 3. 令牌轮换（Rotation）

每次刷新都签发**全新的 Refresh Token**，旧 Refresh Token 立即作废。好处：即使 Refresh Token 泄露，攻击者也只能用一次，泄露窗口最小化。

### 4. 其他细节

- Access Token 不存服务端（无状态、可水平扩展），Refresh Token 存 Redis（有状态、可作废）——两者职责分开；
- `secret` 用 32 字节以上随机值，放环境变量或密钥管理服务，**绝不进代码仓库**；
- 敏感操作（改密码、绑手机）要求重新验证密码，别只靠 Token 权限。

## 小结

双 Token 机制在"无状态可扩展"和"可作废可控"之间取得了平衡：短命的 Access Token 保护主链路性能，长命但可控的 Refresh Token 兜底安全。加上轮换、黑名单和多点登录控制，这套方案足以支撑中大型应用的认证体系。本博客的登录系统正是基于此设计，欢迎在登录后观察网络请求里的令牌刷新流程。
""",
    },
    {
        "title": "Kubernetes 入门：从容器到集群编排",
        "category": "运维部署",
        "tags": ["Kubernetes", "云原生", "Docker", "微服务", "部署"],
        "abstract": "以实践视角讲解 Kubernetes 的核心概念：Pod、Deployment、Service、Ingress，包含资源清单编写、滚动更新与集群架构理解，帮助开发者迈出云原生第一步。",
        "content": """# Kubernetes 入门：从容器到集群编排

## 为什么需要 K8s

Docker 解决了"应用怎么打包"的问题，但一个应用要跑 3 个副本、更新时滚动升级、节点挂了自动调度、流量按域名分发——这些"编排"问题 Docker 本身并不擅长。Kubernetes（K8s）就是容器世界的操作系统：它管理一组服务器（集群），自动决定容器跑在哪、挂了怎么恢复、流量怎么进。

## 集群架构

```mermaid
flowchart TD
    subgraph 控制面
        API[API Server]
        Sched[Scheduler]
        CM[Controller Manager]
    end
    subgraph 工作节点1
        K1[Kubelet] --> P1[Pod]
        K1 --> P2[Pod]
    end
    subgraph 工作节点2
        K2[Kubelet] --> P3[Pod]
    end
    API --> K1
    API --> K2
```

- **控制面**：集群的"大脑"，API Server 是所有操作的唯一入口；
- **工作节点**：真正跑应用的机器，Kubelet 负责执行控制面的指令、维护 Pod 状态。

## 核心对象

### Pod：最小调度单元

Pod 是一个或多个容器的集合，共享网络与存储。通常一个 Pod 里放一个主容器，再加可选的 Sidecar（日志采集、流量代理）。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: blog-api
  labels:
    app: blog-api
spec:
  containers:
    - name: api
      image: myblog/server:1.0.0
      ports:
        - containerPort: 8080
      resources:
        requests: { cpu: 250m, memory: 256Mi }
        limits:   { cpu: "1",   memory: 512Mi }
```

### Deployment：声明期望状态

你**不直接创建 Pod**，而是声明"我要 3 个副本、镜像版本 1.0.0"，Deployment 负责把现实调整到期望状态——这就是声明式 API 的精髓。

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: blog-api
spec:
  replicas: 3
  selector:
    matchLabels: { app: blog-api }
  template:
    metadata:
      labels: { app: blog-api }
    spec:
      containers:
        - name: api
          image: myblog/server:1.0.1   # 升级版本
          readinessProbe:              # 就绪探针：通过才接流量
            httpGet: { path: /health, port: 8080 }
            initialDelaySeconds: 5
```

### 滚动更新与回滚

```bash
# 升级镜像（默认滚动更新：逐批替换，探针失败自动暂停）
kubectl set image deployment/blog-api api=myblog/server:1.0.1

# 查看更新进度
kubectl rollout status deployment/blog-api

# 出问题？一条命令回滚
kubectl rollout undo deployment/blog-api
```

滚动更新 + 就绪探针 = 发布期间服务不中断。这是 K8s 对运维最大的价值之一。

### Service 与 Ingress：流量入口

Pod 的 IP 是临时的，Service 提供稳定的访问入口和负载均衡：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: blog-api-svc
spec:
  selector: { app: blog-api }
  ports:
    - port: 80
      targetPort: 8080
```

集群外部访问，再上一层 Ingress（按域名/路径路由到不同 Service）：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: blog-ingress
spec:
  rules:
    - host: example.com
      http:
        paths:
          - path: /api
            pathType: Prefix
            backend:
              service: { name: blog-api-svc, port: { number: 80 } }
```

## ConfigMap 与 Secret：配置与密钥

```bash
# 配置与镜像解耦
kubectl create configmap blog-config --from-literal=LOG_LEVEL=info
kubectl create secret generic db-secret --from-literal=DB_PASSWORD=root
```

容器内通过环境变量或挂载使用，改配置不必重新构建镜像。

## 入门建议

1. **先手动跑通**：用 `kind` 或 `minikube` 在本地起一个单节点集群，亲手 `kubectl apply` 一个 Deployment；
2. **理解三个心智模型**：声明式期望状态、控制面/工作节点分工、探针（存活/就绪/启动）；
3. **别急着上生产**：生产级集群还要处理存储（PVC）、命名空间、RBAC、Ingress Controller 高可用，建议先在开发环境把核心对象玩熟。

## 小结

K8s 的学习曲线陡，但它的核心思想其实很朴素：**你描述想要的状态，系统负责达成并维持**。从 Deployment 起步，逐步掌握 Service、Ingress、滚动更新，你就拥有了现代云原生基础设施的通用语言。本博客的部署方案目前使用 Docker Compose 足够，文章量再大一个量级时，K8s 会是自然的下一步。
""",
    },
    {
        "title": "2026 个人技术博客搭建心得：从零到上线的完整记录",
        "category": "生活随笔",
        "tags": ["博客", "经验", "随笔", "工作日常"],
        "abstract": "记录从选型、设计到上线一个个人技术博客的全过程，分享极简编辑风的设计思考、部署踩坑与内容创作心得，给同样想建站的朋友一份真实可参考的经验。",
        "content": """# 2026 个人技术博客搭建心得：从零到上线的完整记录

## 为什么突然想建博客

工作第五年，笔记散落在语雀、本地 Markdown 和收藏夹里，想回头找一篇自己写过的解决方案，往往要翻半天。知识不沉淀、不公开，就只是"经验"而不是"资产"。所以 2026 年的第一个大决定：**搭一个属于自己的技术博客**，把踩过的坑、读过的书、想明白的事都写下来。

## 选型：从极简出发

市面上的博客方案很多：Hexo、Hugo、VuePress、WordPress、以及各种"开箱即用"的 SaaS。我的选择标准很朴素：

1. 内容是我的，数据能导出（拒绝被平台绑架）；
2. 样式可控，要符合我对"阅读体验"的执念；
3. 部署成本低，一台小服务器足够。

最终选了自研：前端 Vue3 + Vite（本博客就是），后端 Go + MySQL + Elasticsearch，部署在 Docker Compose 上。**自己写一遍系统，比用一百个现成主题都更能理解博客应该长什么样**。

## 设计：极简编辑风的执念

技术博客的核心价值是文字本身，所以设计上做"减法"：

- 单栏时间线列表，砍掉轮播和侧栏，读者注意力 100% 落在标题和内容上；
- 正文阅读宽度控制在 65ch 左右，这是排版学验证的"黄金行宽"，眼睛扫读最舒服；
- 颜色只用一组克制的中性色（zinc 系）+ 一个强调蓝，深色浅色两套主题；
- 衬线字体只留给文章标题，正文用无衬线，代码用等宽字体——三种字体各司其职。

```css
/* 核心设计 token 的一小部分 */
:root {
  --bg: #FAFAFA;
  --text-primary: #18181B;
  --accent: #2563EB;
  --font-serif: "Source Serif Pro", "Noto Serif SC", Georgia, serif;
}
```

设计没有堆砌任何装饰，但每一个间距、每一道分割线都经过斟酌。**好的设计是看不见的设计**——读者只觉得"读起来舒服"，不会意识到背后有设计。

## 部署踩坑实录

上线过程并非一帆风顺，最有价值的三个坑：

1. **Elasticsearch 内存**：第一次部署，ES 直接吃满 2G 内存导致整机卡死。解决：`ES_JAVA_OPTS=-Xms512m -Xmx512m` 显式限制堆内存；
2. **SPA 路由刷新 404**：直接访问 `/article/123` 返回 404，因为 Nginx 没配 `try_files` 回退到 `index.html`；
3. **HTTPS 证书**：Let's Encrypt 免费证书 + `certbot --nginx` 自动续期，半小时搞定，别再用自签名证书吓浏览器。

## 内容创作：坚持比技巧重要

上线第一周热情高涨，第三周就开始找借口断更。后来给自己定了两条简单规则，一直用到现在：

- **每周至少一篇**，题材不限（技术/读书/随笔），关键是形成节奏；
- **每篇必须解决一个具体问题**，哪怕是"怎么给 Nginx 配 gzip"这种小问题——对别人有参考价值，对自己是知识复利。

三个月下来，最意外的收获不是流量，而是**写作迫使我把模糊的知识重新梳理成结构化的表达**。很多以为"会了"的东西，动笔才发现漏洞百出——写作是最好的学习方式，没有之一。

## 小结

建站这件事，技术选型只占 20% 的难度，剩下 80% 是**长期主义的考验**：持续写、持续改、持续优化。如果你也在犹豫要不要建博客，我的建议是：先有一个能发文章的壳子（哪怕是最简单的 Hexo），然后立刻开始写第一篇。**内容永远比技术重要**——本博客会继续更新，欢迎常来。
""",
    },
    {
        "title": "阅读《代码整洁之道》后的工程实践笔记",
        "category": "读书笔记",
        "tags": ["读书笔记", "代码规范", "随笔"],
        "abstract": "结合《代码整洁之道》的核心观点，讨论命名、函数设计、注释与重构在真实项目中的应用，附重构前后对比案例，记录读书之后的工程实践心得。",
        "content": """# 阅读《代码整洁之道》后的工程实践笔记

## 为什么重读这本书

《代码整洁之道》（Clean Code）是软件行业最经典的"价值观书籍"之一。第一次读的时候觉得"道理都懂"，工作几年、维护过几个没人敢动的老项目之后重读，才发现那些朴素的建议背后全是血泪教训。这篇笔记不罗列全书目录，只写**对我工程实践真正产生影响的五个点**。

## 一、命名：代码可读性的第一杠杆

书里反复强调：**命名是程序员最重要的工作**。可读性差的代码，再好的注释也救不回来。

重构前：

```python
def get_data(a, b):
    lst = []
    for i in range(len(a)):
        if a[i].status == b:
            lst.append(a[i])
    return lst
```

重构后：

```python
def filter_active_articles(articles, status):
    return [a for a in articles if a.status == status]
```

变量名 `articles`、`status` 让函数意图一目了然，甚至不需要看实现。给自己立了一条规矩：**如果写变量名时需要思考超过 2 秒，说明名字没起好**。

## 二、函数：短小 + 单一职责

"函数的第一规则是要短小，第二规则是还要更短小。" 一个函数只做一件事，是划分边界的起点。

重构前（一个函数干了三件事）：

```go
func processOrder(order Order) error {
    // 1. 校验
    if order.Amount <= 0 { return errors.New("invalid amount") }
    // 2. 扣库存 + 落库
    db.UpdateStock(order)
    db.Save(order)
    // 3. 发通知
    email.Send(order.CustomerEmail, "订单已确认")
    return nil
}
```

重构后（每个函数一个职责）：

```go
func ProcessOrder(order Order) error {
    if err := validateOrder(order); err != nil { return err }
    if err := persistOrder(order); err != nil { return err }
    return notifyCustomer(order)
}
```

拆分的收益：每个步骤可以独立测试、独立复用、独立失败处理。

## 三、注释：用代码表达意图，而不是解释行为

书里有个著名的观点："**注释是一种失败，因为我们找不到不用注释就能表达自己的方式**。" 注释应该解释"为什么"，而不是复述"做了什么"。

```go
// 坏注释：复述行为
// 循环遍历用户列表，检查状态是否为 1
for _, u := range users {
    if u.Status == 1 { ... }
}

// 好注释：解释为什么
// 状态 1 = 待激活；超过 72 小时未激活的账号按僵尸号清理，避免发信浪费
if u.Status == 1 && time.Since(u.CreatedAt) > 72*time.Hour { ... }
```

实践建议：**先尝试用命名和结构表达，实在表达不了的"为什么"才写注释**。

## 四、重构：随时进行的小步重构

书里的重构不是一次性的"大动干戈"，而是**每时每刻的微小改进**：改一个名字、拆一个函数、删一段重复代码。配合测试保护，小步重构几乎零风险。

我的工作流：

```mermaid
flowchart LR
    A[发现坏味道] --> B[写/跑测试]
    B --> C[做一次小重构]
    C --> D{测试通过?}
    D -- 是 --> E[提交]
    D -- 否 --> C
```

坏味道清单（书中的经典）：过长函数、过长参数列表、重复代码、魔法数字、注释掉的死代码。每看到一处就顺手处理一处，积少成多。

## 五、错误处理：别吞异常

"吞掉异常就是撒谎——告诉调用方一切正常，其实并没有。" 这一条我踩过最深的坑：

```go
// 坏味道：吞掉错误
resp, _ := http.Get(url)
// 好习惯：显式处理
resp, err := http.Get(url)
if err != nil {
    log.Error("请求失败", "url", url, "err", err)
    return nil, fmt.Errorf("fetch %s: %w", url, err)
}
```

配合 `errors.Is/As` 做错误分类，日志里永远能追溯到根因。

## 小结

《代码整洁之道》与其说教你怎么写代码，不如说教你怎么**对代码负责**：命名认真、函数短小、注释克制、随时重构、错误显式。读完最大的变化是 code review 时不再只说"这里逻辑有问题"，而是能说出"这里命名有歧义、这个函数做了两件事"——这些才是让团队代码质量持续变好的日常力量。
""",
    },
    {
        "title": "AI 辅助编程实践：从 Copilot 到智能体的开发新范式",
        "category": "后端开发",
        "tags": ["人工智能", "开发工具", "开源", "编程语言", "前端"],
        "abstract": "梳理 AI 编程工具从代码补全到自主智能体的演进，分享提示词工程、AI Code Review、测试生成与智能体工作流的实践经验，以及落地时需要注意的风险边界。",
        "content": """# AI 辅助编程实践：从 Copilot 到智能体的开发新范式

## 两年间的范式转移

2024 年 AI 编程助手还停留在"代码补全"，2026 年的今天，主流的工具已经能理解整个代码仓库、跨文件修改、自动跑测试、自主修复失败。这场变化的本质是：**编程从"写代码"变成了"描述意图 + 审查产出"**。

```mermaid
flowchart LR
    A[开发者] -->|自然语言描述需求| B[AI 助手]
    B -->|检索仓库上下文| C[代码库]
    B -->|生成代码/测试| D[修改建议]
    D --> E[开发者 Review]
    E -->|批准| F[提交+跑测试]
    F -->|失败| B
```

## 我日常使用的 AI 工作流

### 1. 补全与对话：解放机械劳动

写 CRUD 接口、ORM 模型、重复的样板代码时，AI 补全能把时间从几十分钟压缩到几分钟。关键是**喂好上下文**：把相关的接口定义、表结构、命名规范贴进对话，产出质量天差地别。

### 2. 测试生成：最被低估的能力

让 AI 根据函数签名和业务逻辑生成单元测试，是目前性价比最高的用法：

```
根据以下 Go 函数生成表驱动测试，覆盖：正常路径、空输入、超限输入、错误分支。
函数签名: func ValidateOrder(order Order) error
约束: 使用 testify，测试命名 TestXxx。
```

AI 生成的测试未必全部合理，但**测试骨架 + 边界枚举**的思路能显著提升覆盖率，我再人工补关键断言。

### 3. AI Code Review：第二双眼睛

PR 提交前，让 AI 以资深 reviewer 的视角检查：命名一致性、错误处理遗漏、潜在空指针、SQL 注入风险。它能稳定发现人类容易漏掉的低级问题——虽然高级设计问题仍需人判断。

```mermaid
flowchart LR
    A[提交 PR] --> B[AI Review: 命名/错误处理/安全]
    B --> C{发现问题?}
    C -- 是 --> D[开发者修复]
    C -- 否 --> E[人工 Review 设计/架构]
    D --> B
```

### 4. 智能体（Agent）：把任务交给会用工具的程序

最新的范式：把"重构这个模块并保证测试通过"这种完整任务交给智能体，它自己会读代码、改文件、跑测试、迭代修复。我的使用原则：**低风险任务放手（样板代码、测试补全），高风险任务必须全程把关**（核心业务逻辑、支付、权限）。

## 提示词工程的三个实战技巧

1. **给足约束**：技术栈、命名规范、目标语言、禁止事项，写清楚比写多更重要：

```
用 Go + GORM 实现文章分页查询，字段包含 title/category/tags，按 created_at 倒序。
禁止引入新依赖，错误统一返回 *AppError，函数命名用动词开头。
```

2. **示例驱动**：给 1~2 个"期望的输入输出"或"现有代码风格片段"，AI 的模仿能力比理解抽象描述强得多；

3. **迭代而非重写**：一次生成不满意，用"修改：把 X 改成 Y，因为 Z"的方式迭代，比重新描述整个需求更稳定。

## 风险边界：AI 编程的红线

- **幻觉代码**：AI 会自信地生成不存在的 API、过期的依赖版本，编译不过/运行时才暴露。**必须跑测试验证，不能盲信**；
- **安全漏洞**：AI 生成的 SQL 拼接、命令执行代码可能引入注入风险，安全敏感代码必须人工审查；
- **版权与许可**：训练语料来源复杂，商用项目注意核对生成代码的许可证问题；
- **思维同质化**：长期依赖 AI 会让开发者失去"从零设计"的能力，复杂架构问题建议先自己想清楚，再用 AI 加速实现。

## 小结

AI 不会取代程序员，但**会用 AI 的程序员会取代不会用的**。它把编程的重心从"怎么实现"推向"要什么、怎么验收"。保持判断力、守住安全红线、持续把重复劳动交给工具——这是 2026 年开发者的新基本功。本博客的很多工具型代码（ORM 模型、测试骨架）也有 AI 的参与，欢迎来讨论你的 AI 工作流。
""",
    },
    {
        "title": "2026 生活随笔：工作日常、美食与旅行的记录",
        "category": "生活随笔",
        "tags": ["工作日常", "美食", "旅行", "生活随笔", "数码"],
        "abstract": "记录 2026 年的生活切片：远程工作的节奏、周末做饭的治愈感、一次说走就走的旅行，以及陪伴日常的数码装备，分享普通人的平凡日子。",
        "content": """# 2026 生活随笔：工作日常、美食与旅行的记录

## 工作：远程与节奏

2026 年，远程办公已经成为常态。我的日常：上午 10 点开始，先花 15 分钟过一遍邮件和 IM 消息，然后进入一天中最宝贵的大块时间——**写代码、写博客、看论文**。下午两点前尽量完成所有需要高度专注的工作，下午留给会议、Code Review 和沟通。

远程工作最大的敌人不是孤独，而是**边界模糊**。我的应对：物理隔离（一个单独的书房）、时间隔离（18 点后不回工作消息）、仪式感（开工时泡一杯咖啡，收工时合上笔记本）。

## 美食：周末的厨房实验

工作日靠外卖和简餐对付，周末的厨房成了我的"精神角落"。这一年学会了三样拿手菜：

- **番茄牛腩**：牛腩焯水后煸炒出油脂，番茄炒出沙再下锅，高压锅 40 分钟，配米饭绝了；
- **蒜香黄油虾**：虾开背去线，黄油小火融化，蒜末炒香，下虾大火 2 分钟，撒黑胡椒和欧芹；
- **葱油拌面**：小葱熬葱油是关键，焦而不糊，拌面时加一勺猪油，香气直接拉满。

做饭的治愈感在于：**它是少数"投入一定有产出"的事**。代码可能跑不通，但一锅牛腩不会骗你。

## 旅行：一次说走就走的短途

年初和几个老同学去了趟川西。高原的天气变化比代码里的状态机还快：上午还是晴空，下午就飘雪。折多山垭口海拔 4298 米，下车走两步就喘，但看到云海翻涌的那一刻，觉得一切都值。

旅行教会我的事：**计划赶不上变化的时候，享受变化**。那些"没按攻略走"的意外——路边的小馆子、客栈老板推荐的野路子景点——反而成了记忆里最亮的部分。

## 数码：陪伴日常的装备

作为写代码的人，对数码产品有种天然的兴趣。2026 年在用的一套：

| 装备 | 用途 | 使用感受 |
|------|------|---------|
| MacBook Pro 14" | 主力开发 | M 系列芯片，编译 Go 项目几乎秒级 |
| 机械键盘 | 写代码 | 茶轴，段落感适中，码字不累 |
| 降噪耳机 | 专注 | 开启降噪那一刻，世界安静了 |
| 墨水屏阅读器 | 看书 | 比手机护眼，通勤路上读完了三本书 |

装备没有追求顶配，够用就好——**工具是服务于生活的，别让生活服务于工具**。

## 小结

这一年工作依旧忙碌，但学会了在忙碌里留白：一顿认真做的饭、一次没计划的旅行、一个关掉消息提示的下午。技术博客写了十几篇，生活随笔这是第一篇——希望以后能多记录一些"代码之外"的东西。日子平凡，但值得被记住。
""",
    },
]

