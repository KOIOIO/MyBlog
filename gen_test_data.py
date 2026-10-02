#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""生成 MyBlog 测试数据：MySQL SQL + ES bulk JSON"""
import json, random, uuid
from datetime import datetime, timedelta

random.seed(20260928)

NOW = datetime.now()
PASSWORD_HASH = "$2a$10$lKSsQ6V3aSCpx05.2FY7w.oapNLUPCz12vjyKfdRprQ1AyOOT8LPi"  # admin123456

# ---------- 素材 ----------
COVERS = [
    "/uploads/image/0c1964e42f5d0cfd52939f741acf43c8-20250226145548.jpg",
    "/uploads/image/15a5c678022fa2abb94c1b838438ae67-20250226150243.jpg",
    "/uploads/image/1f434888fffbce478109f4898cf35fb6-20250226103045.jpg",
    "/uploads/image/524f9042cb70cc3d03a859a6d01b39e4-20250226150254.jpg",
    "/uploads/image/5c61eb33d5fae4e19a6e3c2598334dfa-20250226150300.jpg",
    "/uploads/image/7a13da6b74ad6f7d354b3cbe7669c28f-20250226100230.png",
    "/uploads/image/8b00af3a6a9e7ac0e0ce7be3e5a5d897-20250226150223.jpg",
    "/uploads/image/8e75a6d4e4b737fcd17ab285f535ea8b-20250226103352.jpg",
    "/uploads/image/a3c59f20727634d4d34722b349a7cbd9-20250226150311.jpg",
    "/uploads/image/cc12bb3c2da4e86ea64370df0235c672-20250226150212.jpg",
    "/uploads/image/db7281566adb2f7eb84d5209cd46a5e7-20250226150249.jpg",
    "/uploads/image/dfe3b636ca868aa32b7c3996c7558454-20250226150235.jpg",
]

ARTICLES = [
    ("Go 语言 GORM 实战：从入门到优雅地操作 MySQL", "backend", ["Go", "GORM", "MySQL", "后端"], "本文系统讲解 Go 生态中最流行的 ORM 框架 GORM，涵盖连接池配置、事务、关联查询与软删除等核心用法，并结合实际项目代码演示如何写出优雅的数据库操作。"),
    ("Vue3 + Vite 从零搭建企业级前端工程", "frontend", ["Vue", "Vite", "TypeScript", "前端"], "介绍如何基于 Vue3 Composition API 与 Vite 搭建可维护的前端工程，包含目录规范、Pinia 状态管理、路由权限与自动化部署等最佳实践。"),
    ("Redis 缓存策略设计：穿透、击穿与雪崩的攻防", "backend", ["Redis", "缓存", "架构", "后端"], "深入分析缓存穿透、缓存击穿与缓存雪崩三大经典问题的成因与解决方案，给出布隆过滤器、互斥锁、多级缓存等实用代码示例。"),
    ("Elasticsearch 全文检索实战：博客搜索系统的设计", "backend", ["Elasticsearch", "搜索", "Go", "架构"], "以一个真实博客系统为例，讲解 ES 索引设计、分词器选择、中文搜索优化以及数据同步方案，帮助读者快速构建高性能全文检索服务。"),
    ("Gin 框架中间件机制深度解析", "backend", ["Go", "Gin", "中间件", "后端"], "拆解 Gin 中间件的洋葱模型原理，从源码角度分析其执行顺序与上下文传递机制，并手写 JWT 认证、日志记录等常用中间件。"),
    ("Docker 部署前后端分离项目的完整指南", "devops", ["Docker", "部署", "Nginx", "运维"], "从零讲解如何用 Docker Compose 编排 MySQL、Redis、Elasticsearch 与前后端服务，实现一键启动的现代化部署方案。"),
    ("2026 年个人技术博客搭建心得与踩坑记录", "life", ["博客", "随笔", "经验"], "记录从选型到上线一个个人博客的全过程，分享遇到的坑、性能优化思路以及内容创作的心得体会。"),
    ("阅读《代码整洁之道》后的工程实践笔记", "reading", ["读书笔记", "代码规范", "随笔"], "结合书中核心观点，讨论命名、函数设计、注释与重构在真实项目中的应用，附个人实践案例。"),
    ("Nginx 反向代理与负载均衡配置详解", "devops", ["Nginx", "负载均衡", "运维", "架构"], "讲解 Nginx 作为反向代理的核心配置，包括静态资源缓存、gzip、HTTPS 证书配置与多节点负载均衡策略。"),
    ("JWT 双 Token 认证机制的设计与实现", "backend", ["JWT", "认证", "安全", "Go"], "从安全角度设计 Access Token + Refresh Token 双令牌机制，讲解黑名单、多点登录控制与令牌轮换策略，附完整 Go 实现。"),
    ("现代前端工程化：从 Webpack 到 Vite 的迁移实战", "frontend", ["前端", "Vite", "Webpack", "工程化"], "对比 Webpack 与 Vite 的构建原理，分享一个中大型项目迁移到 Vite 过程中的问题与解决方案，构建速度提升 5 倍的真实数据。"),
    ("MySQL 索引优化实战：慢查询排查与调优", "backend", ["MySQL", "索引", "性能优化", "后端"], "通过真实慢查询案例，讲解 EXPLAIN 解读、索引选择原则、覆盖索引与最左前缀匹配等优化技巧，附性能对比数据。"),
]

CATEGORY_LABEL = {
    "backend": "后端开发",
    "frontend": "前端开发",
    "devops": "运维部署",
    "life": "生活随笔",
    "reading": "读书笔记",
}

USER_NAMES = ["山间清风", "代码旅人", "夜航星", "南山南", "起风了", "月下独酌", "追光者", "小满", "拾光者"]
USER_EMAILS = ["user1@test.com", "user2@test.com", "user3@test.com", "user4@test.com", "user5@test.com",
               "user6@test.com", "user7@test.com", "user8@test.com", "user9@test.com"]

COMMENT_TEXTS = [
    "写得太好了，正好解决了我最近遇到的问题，感谢分享！",
    "博主讲得很清楚，尤其是中间件那部分，看懂了。",
    "mark 一下，周末照着实践一遍。",
    "请问这里如果数据量大的话性能怎么样？",
    "已三连，期待下一篇！",
    "这个方法很实用，比我之前用的方案优雅多了。",
    "有一点不太明白，能详细讲讲吗？",
    "收藏了，正好最近在做类似的项目。",
    "思路很清晰，跟着做了一遍成功了。",
    "补充一点，新版框架里这个接口有变化，可以更新下。",
    "写得不错，学到了很多，感谢作者。",
    "代码风格很舒服，已关注。",
    "这个问题困扰我很久了，终于找到答案。",
    "建议加个实战案例的完整源码，会更好理解。",
]

def rand_time(days_back=(5, 60)):
    d = random.randint(days_back[0], days_back[1])
    h, m, s = random.randint(0, 23), random.randint(0, 59), random.randint(0, 59)
    return (NOW - timedelta(days=d, hours=h, minutes=m, seconds=s)).strftime("%Y-%m-%d %H:%M:%S")

def rand_mysql_time(days_back=(5, 60)):
    return (NOW - timedelta(days=random.randint(days_back[0], days_back[1]),
                            hours=random.randint(0, 23), minutes=random.randint(0, 59),
                            seconds=random.randint(0, 59))).strftime("%Y-%m-%d %H:%M:%S")

# ---------- 1. 用户 ----------
users = []  # (id, uuid, username, email)
sql_users = []
for i in range(9):
    uid = i + 2  # id 1 是管理员
    u = str(uuid.uuid4())
    users.append((uid, u, USER_NAMES[i], USER_EMAILS[i]))
    sql_users.append(
        f"INSERT INTO users (id, created_at, updated_at, uuid, username, password, email, avatar, address, role_id, register, freeze) "
        f"VALUES ({uid}, '{rand_mysql_time()}', '{rand_mysql_time()}', '{u}', '{USER_NAMES[i]}', '{PASSWORD_HASH}', "
        f"'{USER_EMAILS[i]}', '/image/avatar.jpg', '河南省郑州市', 1, 0, 0);"
    )

# ---------- 2. 文章 (ES) ----------
es_docs = []
category_count = {}
tag_count = {}
article_meta = []  # (id, cover, title, category, tags)
for idx, (title, cat, tags, abstract) in enumerate(ARTICLES, start=1):
    aid = str(idx)
    cover = COVERS[(idx - 1) % len(COVERS)]
    created = rand_time((5, 60))
    views = random.randint(200, 8000)
    content = f"""# {title}

## 引言

这是一篇关于 **{title}** 的实战分享文章，内容基于真实项目经验整理。

## 核心要点

- 从零开始的完整实践路径
- 常见问题与解决方案
- 性能与工程化的思考

## 正文

本文通过详细的代码示例和项目实践，深入讲解了相关技术方案的设计思路与落地过程。
在实际项目中，我们需要注意边界条件的处理，以及不同场景下的取舍。

### 示例代码

```go
func main() {{
    fmt.Println("Hello, MyBlog!")
}}
```

## 总结

希望这篇文章对你有帮助，欢迎在评论区交流讨论。
"""
    es_docs.append({
        "created_at": created,
        "updated_at": created,
        "cover": cover,
        "title": title,
        "keyword": title,
        "category": CATEGORY_LABEL[cat],
        "tags": tags,
        "abstract": abstract,
        "content": content,
        "views": views,
        "comments": 0,  # 后面回填
        "likes": 0,     # 后面回填
    })
    category_count[cat] = category_count.get(cat, 0) + 1
    for t in tags:
        tag_count[t] = tag_count.get(t, 0) + 1
    article_meta.append((aid, cover, title, CATEGORY_LABEL[cat], tags))

# ---------- 3. 点赞 (article_likes) ----------
likes = []
for aid, *_ in article_meta:
    n = random.randint(2, 6)
    likers = random.sample(users, n)
    for uid, u, uname, uemail in likers:
        likes.append((aid, uid, rand_mysql_time()))

# ---------- 4. 评论 (comments) ----------
comments = []  # (id, article_id, user_uuid, content, p_id, created_at)
cid = 1
for aid, *_ in article_meta:
    n = random.randint(2, 5)
    for _ in range(n):
        user = random.choice(users)
        comments.append((cid, aid, user[1], random.choice(COMMENT_TEXTS), None, rand_mysql_time((1, 50))))
        cid += 1
# 回复（约 1/3 评论带一条回复）
replies = []
for c in comments[:]:
    if random.random() < 0.35:
        reply_user = random.choice(users)
        replies.append((cid, c[1], reply_user[1], f"回复：{random.choice(COMMENT_TEXTS)}", c[0], rand_mysql_time((1, 30))))
        cid += 1
comments.extend(replies)

# ---------- 回填 ES 计数 ----------
like_count = {}
comment_count = {}
for aid, *_ in likes:
    like_count[aid] = like_count.get(aid, 0) + 1
for c in comments:
    comment_count[c[1]] = comment_count.get(c[1], 0) + 1
for i, doc in enumerate(es_docs):
    aid = str(i + 1)
    doc["likes"] = like_count.get(aid, 0)
    doc["comments"] = comment_count.get(aid, 0)

# ---------- 输出 SQL ----------
sql = ["SET FOREIGN_KEY_CHECKS = 0;", "USE blog_db;", "START TRANSACTION;"]
sql += sql_users

for cat, num in category_count.items():
    sql.append(f"INSERT INTO article_categories (category, number) VALUES ('{CATEGORY_LABEL[cat]}', {num});")
for tag, num in tag_count.items():
    sql.append(f"INSERT INTO article_tags (tag, number) VALUES ('{tag}', {num});")

# images（封面）
for aid, cover, title, cat, tags in article_meta:
    name = cover.split("/")[-1]
    sql.append(f"INSERT INTO images (created_at, updated_at, name, url, category, storage) "
               f"VALUES ('{rand_mysql_time()}', '{rand_mysql_time()}', '{name}', '{cover}', 3, 0);")

for aid, uid, t in likes:
    sql.append(f"INSERT INTO article_likes (created_at, updated_at, article_id, user_id) VALUES ('{t}', '{t}', '{aid}', {uid});")

for cid, aid, uuid_, content, p_id, t in comments:
    p = f"{p_id}" if p_id else "NULL"
    sql.append(f"INSERT INTO comments (id, created_at, updated_at, article_id, p_id, user_uuid, content) "
               f"VALUES ({cid}, '{t}', '{t}', '{aid}', {p}, '{uuid_}', '{content.replace(chr(39), chr(39)+chr(39))}');")

sql.append("COMMIT;")
sql.append("SET FOREIGN_KEY_CHECKS = 1;")

with open("/tmp/test_data.sql", "w", encoding="utf-8") as f:
    f.write("\n".join(sql))

# ---------- 输出 ES bulk ----------
bulk = []
for i, doc in enumerate(es_docs):
    bulk.append(json.dumps({"index": {"_index": "article_index", "_id": str(i + 1)}}, ensure_ascii=False))
    bulk.append(json.dumps(doc, ensure_ascii=False))
with open("/tmp/articles_bulk.json", "w", encoding="utf-8") as f:
    f.write("\n".join(bulk) + "\n")

print(f"users={len(users)} articles={len(es_docs)} likes={len(likes)} comments={len(comments)}")
print(f"categories={category_count}")
print(f"tags={tag_count}")
print("SQL -> /tmp/test_data.sql")
print("ES   -> /tmp/articles_bulk.json")
