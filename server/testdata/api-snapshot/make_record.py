#!/usr/bin/env python3
"""将 curl 捕获的响应组装为快照记录 JSON: {name, method, path, status, body}"""
import json
import sys
import os

name, method, path, code, body_file = sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4], sys.argv[5]
with open(body_file, 'rb') as f:
    body = f.read().decode('utf-8', 'replace')
rec = {"name": name, "method": method, "path": path, "status": int(code), "body": body}
print(json.dumps(rec, ensure_ascii=False))
