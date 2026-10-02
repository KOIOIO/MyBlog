#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""导入 15 篇真实文章到 ES，并同步 MySQL 统计表（article_categories / article_tags / images）。
浏览量、评论数、点赞数全部为 0。"""
import json, sys, random, subprocess
from datetime import datetime, timedelta

sys.path.insert(0, "/Users/xiaoyuwang/projects/MyBlog")
from articles_content import ARTICLES

random.seed(20260928)
NOW = datetime.now()

# 1. 封面路径：已按序下载为 01.jpg ~ 15.jpg（对应文章 1~15）
covers = [f"/uploads/image/{i:02d}.jpg" for i in range(1, 16)]
assert len(covers) == len(ARTICLES), f"封面数 {len(covers)} != 文章数 {len(ARTICLES)}"

# 2. 生成时间：近 60 天随机分布（固定 seed 可复现）
def rand_time(days_back=(5, 60)):
    d = random.randint(days_back[0], days_back[1])
    h, m, s = random.randint(0, 23), random.randint(0, 59), random.randint(0, 59)
    return (NOW - timedelta(days=d, hours=h, minutes=m, seconds=s)).strftime("%Y-%m-%d %H:%M:%S")

# 3. 组装 ES bulk
bulk = []
for i, a in enumerate(ARTICLES, 1):
    created = rand_time()
    doc = {
        "created_at": created,
        "updated_at": created,
        "cover": covers[i - 1],
        "title": a["title"],
        "keyword": a["title"],
        "category": a["category"],
        "tags": a["tags"],
        "abstract": a["abstract"],
        "content": a["content"],
        "views": 0,
        "comments": 0,
        "likes": 0,
    }
    bulk.append(json.dumps({"index": {"_index": "article_index", "_id": str(i)}}, ensure_ascii=False))
    bulk.append(json.dumps(doc, ensure_ascii=False))

with open("/tmp/new_articles_bulk.json", "w", encoding="utf-8") as f:
    f.write("\n".join(bulk) + "\n")

# 4. 统计分类与标签
cat_count, tag_count = {}, {}
for a in ARTICLES:
    cat_count[a["category"]] = cat_count.get(a["category"], 0) + 1
    for t in a["tags"]:
        tag_count[t] = tag_count.get(t, 0) + 1

# 5. 组装 MySQL SQL
sql = ["SET FOREIGN_KEY_CHECKS = 0;", "USE blog_db;", "START TRANSACTION;"]
for cat, num in cat_count.items():
    sql.append(f"INSERT INTO article_categories (category, number) VALUES ('{cat}', {num});")
for tag, num in tag_count.items():
    sql.append(f"INSERT INTO article_tags (tag, number) VALUES ('{tag}', {num});")
for i, a in enumerate(ARTICLES, 1):
    name = covers[i - 1].split("/")[-1]
    t = rand_time()
    sql.append(f"INSERT INTO images (created_at, updated_at, name, url, category, storage) "
               f"VALUES ('{t}', '{t}', '{name}', '{covers[i - 1]}', 3, 0);")
sql.append("COMMIT;")
sql.append("SET FOREIGN_KEY_CHECKS = 1;")

with open("/tmp/new_articles_meta.sql", "w", encoding="utf-8") as f:
    f.write("\n".join(sql))

print(f"articles={len(ARTICLES)} covers={len(covers)}")
print(f"categories={cat_count}")
print(f"tags={sorted(tag_count)} ({len(tag_count)})")
print("ES bulk -> /tmp/new_articles_bulk.json")
print("SQL     -> /tmp/new_articles_meta.sql")
