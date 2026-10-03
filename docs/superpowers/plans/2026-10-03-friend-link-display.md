# 友链展示重构 Implementation Plan

> **执行方式：**
>
>  本计划由 MainAgent 在用户审核通过后于本会话内联执行，按 Task 顺序逐项完成并提交。

**Goal:** 将友链页从 "静态网格 + 压扁 logo" 重构为 "顶部无限滚动 Logo 墙 + 下方详情卡片"，彻底解决 logo 变形，提升视觉表现。

**Architecture:** 新增一个纯展示型 Vue SFC 组件 `FriendLinkMarquee.vue` 承载滚动 Logo 墙（纯 CSS 动画、无缝循环、hover 暂停、边缘渐隐、灰度变彩）；改造 `friend-link/index.vue` 引入该组件，并把卡片 logo 由 `el-image`（默认 fill 拉伸）改为原生 `<img>` + `object-fit: contain` + 固定高度 / 自适应宽度。零后端改动、零新依赖。

**Tech Stack:** Vue 3 `<script setup lang="ts">` + SCSS（scoped）+ Element Plus（仅保留页面其余部分）+ 现有全局 CSS 变量。

**Spec:** `docs/friend-link-display-spec.md`（同目录，评审时一并阅读）

## Global Constraints



* 数据模型保持 `FriendLink { logo, link, name, description }` 不变；不新增 i18n 词条。

* 所有颜色 / 间距 / 字号使用现有 CSS 变量：`--bg`、`--bg-elevated`、`--border`、`--text-primary`、`--text-muted`、`--accent`、`--radius-sm`、`--sp-1/2/3/6/8`、`--fs-12` 等，禁止硬编码色值。

* 滚动墙条目间距必须用 `margin-right`（非 flex gap），以保证 `translateX(-50%)` 无缝循环对齐。

* 支持 `prefers-reduced-motion: reduce`：动画禁用、隐藏复制段、条目居中换行。

* 全部改动仅限下述两个文件；不动后端、路由、管理后台。

* 完成后必须通过 `cd web && npm run type-check` 与本地浏览器验证。



***

### Task 1: 创建滚动 Logo 墙组件

**Files:**



* Create: `web/src/components/pages/FriendLinkMarquee.vue`

**Interfaces:**



* Consumes: `FriendLink` 类型（来自 `@/api/friend-link`，字段 `logo/link/name/description`）

* Produces: 默认导出组件，props `items: FriendLink[]`；不向外 emit（点击跳转内部处理）

- [ ] **Step 1: 创建组件文件**



```
<template>
  <div class="friend-link-marquee">
    <div class="marquee-track">
      <!-- 内容复制为两组，配合 translateX(-50%) 实现无缝循环 -->
      <div v-for="(group, g) in [0, 1]" :key="g" class="marquee-group">
        <div
          v-for="(item, i) in items"
          :key="`${item.name}-${g}-${i}`"
          class="marquee-item"
          :title="item.name"
          @click="open(item.link)"
        >
          <img class="marquee-logo" :src="item.logo" :alt="item.name" loading="lazy"/>
          <span class="marquee-name">{{ item.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FriendLink } from '@/api/friend-link'

defineProps<{ items: FriendLink[] }>()

const open = (link: string) => {
  window.open(link)
}
</script>

> 实现注记：无缝循环采用「双分组」结构（`v-for="[0,1]"` 渲染两个 `.marquee-group`），
> 替代初版 `nth-child(n + calc(var(--count) + 1))` 方案——后者会被 Sass 解析器
> 以 "Expected a number" 拒绝（`calc` 不允许出现在 `nth-child` 的 `An+B` 中）。
> 双分组下 `translateX(-50%)` 恰好对齐一组宽度，循环无缝；减弱动效时直接
> `.marquee-group:nth-child(2) { display: none; }` 隐藏复制组。

<style scoped lang="scss">
.friend-link-marquee {
  overflow: hidden;
  margin-bottom: var(--sp-6);
  padding: var(--sp-2) 0;
  -webkit-mask-image: linear-gradient(90deg, transparent, #000 10%, #000 90%, transparent);
  mask-image: linear-gradient(90deg, transparent, #000 10%, #000 90%, transparent);

  &:hover .marquee-track {
    animation-play-state: paused;
  }

  .marquee-track {
    display: flex;
    width: max-content;
    animation: marquee-scroll 32s linear infinite;
  }

  .marquee-group {
    display: flex;
  }

  .marquee-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--sp-1);
    margin-right: var(--sp-8);
    cursor: pointer;

    .marquee-logo {
      height: 44px;
      max-width: 140px;
      object-fit: contain;
      filter: grayscale(1);
      opacity: 0.75;
      transition: filter 200ms ease-out, opacity 200ms ease-out;
    }

    .marquee-name {
      max-width: 140px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      font-size: var(--fs-12);
      color: var(--text-muted);
    }

    &:hover .marquee-logo {
      filter: grayscale(0);
      opacity: 1;
    }
  }

  @media (max-width: 768px) {
    .marquee-logo {
      height: 36px;
    }
    .marquee-item {
      margin-right: var(--sp-6);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .marquee-track {
      animation: none;
      flex-wrap: wrap;
      justify-content: center;
    }
    // 隐藏复制出来的第二组，只展示一轮
    .marquee-group:nth-child(2) {
      display: none;
    }
  }
}

@keyframes marquee-scroll {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(-50%);
  }
}
</style>
```



* [ ] **Step 2: 语法自检**



```
cd /Users/xiaoyuwang/projects/MyBlog/web && npx vue-tsc --noEmit -p tsconfig.app.json 2>&1 | grep -i "FriendLinkMarquee" || echo "FriendLinkMarquee: no type errors"
```

Expected: 输出 `FriendLinkMarquee: no type errors`（若整仓其他文件本就有错误，仅需确认无本组件相关错误）。



* [ ] **Step 3: 提交**



```
cd /Users/xiaoyuwang/projects/MyBlog && git add web/src/components/pages/FriendLinkMarquee.vue && git commit -m "feat(friend-link): add infinite marquee logo wall component"
```



***

### Task 2: 改造友链页（集成滚动墙 + 修复卡片 logo 变形）

**Files:**



* Modify: `web/src/views/web/friend-link/index.vue`

**Interfaces:**



* Consumes: Task 1 的 `FriendLinkMarquee` 组件（默认导出，props `items: FriendLink[]`）

* Produces: 重构后的友链页（滚动墙 + 卡片网格）

- [ ] **Step 1: 模板中引入滚动墙、卡片 logo 改为原生 img**

将 `index.vue` 的 `<template>` 修改为：



```
<template>
  <div class="friend-link-page">
    <web-navbar :noScroll="true"/>
    <div class="page">
      <h1 class="page-title">{{ t('pages.friendLink.title') }}</h1>
      <friend-link-marquee v-if="friendLinkList.length" :items="friendLinkList"/>
      <div class="list">
        <div v-for="(item, i) in friendLinkList" :key="item.name" class="link-card"
             :style="{ animationDelay: Math.min(i, 5) * 60 + 'ms' }"
             @click="handleFriendLinkJumps(item.link)">
          <img class="card-logo" :src="item.logo" :alt="item.name" loading="lazy"/>
          <div class="card-body">
            <h3 class="card-name">{{ item.name }}</h3>
            <p class="card-desc">{{ item.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
```



* [ ] **Step 2: script 中引入组件**

将 `<script setup lang="ts">` 首行处新增：



```
import FriendLinkMarquee from "@/components/pages/FriendLinkMarquee.vue";
```

（其余逻辑 `friendLinkList`、`getFriendLinkInfo`、`handleFriendLinkJumps`、`t` 保持不变。）



* [ ] **Step 3: 样式调整 —— 卡片 logo 保持宽高比**

将 scoped 样式中 `.link-card .card-logo` 规则（现为 `el-image` 用的圆角样式）替换为：



```
    .card-logo {
      height: 48px;
      max-width: 96px;
      object-fit: contain;
      border-radius: var(--radius-sm);
      flex-shrink: 0;
    }
```

其余样式（`.list` 网格、`.link-card` hover/active、`.card-name`、`.card-desc`）保持不变。



* [ ] **Step 4: 类型检查**



```
cd /Users/xiaoyuwang/projects/MyBlog/web && npm run type-check
```

Expected: 通过（exit 0）。若本文件外存在存量类型错误，确认新增代码无新增报错即可。



* [ ] **Step 5: 提交**



```
cd /Users/xiaoyuwang/projects/MyBlog && git add web/src/views/web/friend-link/index.vue && git commit -m "fix(friend-link): preserve logo aspect ratio and integrate marquee wall"
```



***

### Task 3: 本地运行与浏览器验证

**Files:**



* 无代码改动；验证 `http://localhost/friend-link`

- [ ] **Step 1: 确认前端服务在跑**



```
lsof -iTCP:80 -sTCP:LISTEN | head -3
```

Expected: 有监听进程（如无，执行 `cd /Users/xiaoyuwang/projects/MyBlog/web && npm run dev` 后台启动）。



* [ ] **Step 2: 打开友链页逐条验收（对照 Spec §7）**



| # | 验收点                               | 操作                    |
| - | --------------------------------- | --------------------- |
| 1 | 滚动墙位于标题下方，无缝循环、无跳变                | 目视 30s                |
| 2 | logo 无拉伸变形                        | 对比重构前截图               |
| 3 | hover 暂停；单个 logo 灰→彩；移出恢复         | 鼠标操作                  |
| 4 | 点击条目 / 卡片新标签页打开                   | 点击验证                  |
| 5 | 375px 宽度无横向溢出                     | DevTools 设备模拟         |
| 6 | 控制台无报错                            | DevTools Console      |
| 7 | prefers-reduced-motion 下无动画、仅一轮居中 | DevTools Rendering 模拟 |



* [ ] **Step 3: 截图留档（重构后）**



```
cd /Users/xiaoyuwang/projects/MyBlog && mkdir -p docs/assets && echo "截图保存至 docs/assets/friend-link-after.png（浏览器手动截图）"
```

Expected: `docs/assets/friend-link-after.png` 存在（供评审对比）。



***

### Task 4: 收尾提交与交付



* [ ] **Step 1: 更新文档**

在 `docs/friend-link-display-spec.md` 顶部状态行将 `状态：待评审` 改为 `状态：已实施`（用户确认验收后执行）。



* [ ] **Step 2: 最终提交**



```
cd /Users/xiaoyuwang/projects/MyBlog && git add -A && git commit -m "docs(friend-link): mark display refactor spec as implemented"
```



* [ ] **Step 3: 向用户交付**

- 通过 `present_files` 交付：`docs/assets/friend-link-after.png`（重构后截图）、`docs/friend-link-display-spec.md`（最终版 spec）。

- 总结：改动文件清单、验收结果、遗留事项（如有）。