# MyBlog UI 重构・设计提示词手册

> 基于 8 个标杆内容站调研提炼：Lee Robinson、Josh Comeau、Ghost Casper、Medium、Substack、Awwwards 获奖站（Kusaludhara）、少数派、2025-2026 排版趋势。
> 适用对象：Vue 3 + TypeScript 技术博客，含前台（首页 / 文章 / 搜索 / 关于 / 友链 / 新闻）与管理后台。



***

## 一、先选风格方向（三选一）

### 方向 A：极简编辑风（Minimal Editorial）—— 推荐技术博客首选

**灵感来源**：Lee Robinson ([leerob.com](https://leerob.com)) + Medium + Ghost Casper

**底层逻辑**：技术内容的核心价值是文字本身。极简风通过 "减法" 让读者注意力 100% 落在内容上；单栏时间线列表消除了传统博客 "轮播 + 侧栏" 的信息噪音，65ch 阅读宽度符合人眼最佳扫读范围（每行 65-70 字符是排版学百年验证的黄金值）。

**设计系统 Token**：



```
颜色（浅色）：

&#x20; \--bg:            #FAFAFA      页面底色，非纯白，减少眩光

&#x20; \--bg-elevated:   #FFFFFF      卡片/浮层

&#x20; \--text-primary:  #18181B      标题（zinc-900）

&#x20; \--text-body:     #3F3F46      正文（zinc-700）

&#x20; \--text-muted:    #71717A      辅助/元信息（zinc-500）

&#x20; \--border:        #E4E4E7      分割线（zinc-200）

&#x20; \--accent:        #2563EB      强调色（蓝-600），仅用于链接/hover/选中

颜色（深色）：

&#x20; \--bg:            #09090B      近黑（zinc-950），非纯黑，更柔和

&#x20; \--bg-elevated:   #18181B      卡片（zinc-900）

&#x20; \--text-primary:  #FAFAFA      标题（zinc-50）

&#x20; \--text-body:     #D4D4D8      正文（zinc-300），对比度约 11:1，符合 WCAG AA

&#x20; \--text-muted:    #71717A      辅助（zinc-500）

&#x20; \--border:        #27272A      分割线（zinc-800）

&#x20; \--accent:        #60A5FA      强调色（蓝-400），深色下提亮

字体：

&#x20; \--font-sans:  "Inter", "Geist", -apple-system, "PingFang SC", "Microsoft YaHei", sans-serif

&#x20; \--font-serif: "Source Serif Pro", "Noto Serif SC", Georgia, serif   仅文章标题/引用

&#x20; \--font-mono:  "JetBrains Mono", "Fira Code", monospace              代码块

字号阶梯（px）：

&#x20; 12 / 14 / 16(body) / 18 / 20 / 24 / 30 / 36 / 48

&#x20; 正文行高 1.75，标题行高 1.25，段落间距 1.5rem

间距（8px 栅格）：

&#x20; 4 / 8 / 12 / 16 / 24 / 32 / 48 / 64 / 96

&#x20; 内容最大宽度 720px（文章）/ 1120px（列表页）

圆角：6px（按钮/标签）/ 12px（卡片）/ 0（分割线风格可选）

阴影：浅色用 0 1px 2px rgba(0,0,0,0.04)；深色几乎不用阴影，靠 border 分层
```

**视觉关键词**：大量留白、单栏时间线、细线分割、无衬线为主、serif 点缀标题、悬停下划线动画、主题切换按钮。



***

### 方向 B：深色精致风（Refined Dark）—— 适合极客 / 开发者个人品牌

**灵感来源**：Josh Comeau + Awwwards Kusaludhara (#0A0A0A / #F2F2F2)

**底层逻辑**：开发者群体天然偏好深色模式（长时间阅读护眼、代码与背景融合）。但 "深色 ≠ 反色"—— 纯黑底 + 纯白字会产生光晕效应（halation），所以用近黑底 + 米白字；通过微妙的渐变光晕（glow）和玻璃拟态（glassmorphism）营造空间感，而不是靠边框。

**设计系统 Token**：



```
颜色：

&#x20; \--bg:            #0A0A0A      近黑底色

&#x20; \--bg-elevated:   #141414      卡片，比底色亮 5%

&#x20; \--bg-glass:      rgba(255,255,255,0.03)   玻璃层

&#x20; \--text-primary:  #F2F2F2      近白标题

&#x20; \--text-body:     #C4C4C4      正文，降亮度护眼

&#x20; \--text-muted:    #6B6B6B      辅助

&#x20; \--border:        rgba(255,255,255,0.08)   极淡边框

&#x20; \--accent:        #A78BFA      紫罗兰（violet-400），呼应你现有的紫色基因

&#x20; \--accent-glow:   radial-gradient(circle, rgba(167,139,250,0.15), transparent 70%)

&#x20; \--code-bg:       #0D0D0D      代码块，比页面再深一档

字体：

&#x20; \--font-sans:  "Inter", "Geist", "PingFang SC", sans-serif

&#x20; \--font-mono:  "JetBrains Mono", monospace

&#x20; 标题可用 font-weight: 600 + letter-spacing: -0.02em（负字距营造高级感）

特殊效果：

&#x20; \- 首屏 Hero 区：背景叠加 radial-gradient 紫色光晕，随鼠标轻微移动（parallax）

&#x20; \- 卡片 hover：border-color 过渡到 accent，同时 translateY(-2px)

&#x20; \- 代码块：左上角三个圆点（mac 窗口风），语法高亮用 One Dark 主题

&#x20; \- 滚动条：自定义细滚动条，宽度 6px，滑块 accent 色
```

**视觉关键词**：近黑、紫罗兰光晕、玻璃拟态、负字距大标题、mac 风代码块、鼠标跟随光晕、无轮播。



***

### 方向 C：杂志内容风（Magazine Content）—— 适合内容量大、多分类的博客

**灵感来源**：Substack + 少数派 + Ghost 付费主题（EstudioPatagon）

**底层逻辑**：当文章数量多、分类丰富时，纯时间线会让读者迷失。杂志风用 "网格 + 封面图 + 分类标签" 建立视觉索引，让读者像翻杂志一样浏览；首屏 Featured 文章用大封面建立品牌调性，下方网格承载长尾内容。

**设计系统 Token**：



```
颜色（浅色为主，深色可选）：

&#x20; \--bg:            #FDFCFB      暖白，纸张感

&#x20; \--bg-elevated:   #FFFFFF

&#x20; \--text-primary:  #1C1917      暖黑（stone-900）

&#x20; \--text-body:     #44403C      （stone-700）

&#x20; \--text-muted:    #78716C      （stone-500）

&#x20; \--border:        #E7E5E4      （stone-200）

&#x20; \--accent:        #B45309      琥珀色（amber-700），温暖、有纸质感

&#x20; \--category-tag-bg: #F5F5F4    标签底色

字体：

&#x20; \--font-serif: "Source Serif Pro", "Noto Serif SC", Georgia, serif   标题主力

&#x20; \--font-sans:  "Inter", "PingFang SC", sans-serif                    正文/UI

&#x20; 文章标题用 serif 36-48px，营造出版物质感

布局：

&#x20; \- 首屏：1 篇 Featured 大封面（左图右文，占满宽度）

&#x20; \- 第二层：2-3 篇次精选，中等卡片

&#x20; \- 第三层：文章网格，3 列，每卡含封面缩略图+分类标签+标题+摘要+日期

&#x20; \- 侧栏：仅在文章详情页出现，含目录（TOC）、作者卡、相关推荐

&#x20; \- 首页无侧栏，全宽内容
```

**视觉关键词**：暖白纸感、serif 大标题、封面图驱动、分类标签、Featured 首屏、无轮播、阅读进度条。



***

## 二、通用设计原则（三个方向都必须遵守）



1. **消灭轮播图**：现有首页的 Carousel 是 "老 UI" 的标志性元素。轮播图点击率极低（第二张以后 <1%），且占据首屏黄金位置。用 Featured 大卡或时间线替代。

2. **行宽控制**：文章正文 `max-width: 65ch`（约 720px），永远不要让文字铺满宽屏。这是可读性的第一原则。

3. **垂直节奏**：段落 `margin-bottom: 1.5rem`，标题 `margin-top: 2.5rem`，用一致的间距创造阅读呼吸感。

4. **深色模式不是反色**：深色背景用近黑（#09090B / #0A0A0A），文字用米白（#D4D4D8 / #F2F2F2），禁止 #000 + #FFF 组合。

5. **强调色克制**：全站只有一个 accent 色，仅用于链接、hover、选中态、按钮。不要用彩色装饰元素。

6. **动效服务于功能**：hover 过渡 150-200ms ease-out；页面切换用淡入；禁止自动播放动画、禁止弹跳 / 旋转等花哨效果。

7. **图标统一**：用 `lucide-vue-next`（线性、2px 描边、24px 基准），不要混用 Element Plus 图标和 emoji。

8. **中文排版**：中英文之间自动加空格（可用 `pangu.js`）；标点悬挂；代码块内英文用 mono，中文回退到 sans。



***

## 三、逐页面提示词（可直接喂给 AI 设计 / 编码工具）

> 使用方式：复制对应页面的 prompt，加上你选的方向（A/B/C）的设计系统 Token，即可生成。

### 3.1 首页（Home / 文章列表）



```
为一个 Vue 3 技术博客设计首页，采用【方向 A：极简编辑风】。

布局结构（从上到下，单栏居中，最大宽度 1120px）：

1\. 顶部导航栏：固定顶部，滚动时背景从透明过渡到半透明毛玻璃。左侧 Logo（文字 "MyBlog"，font-weight 700），中间导航链接（首页 / 分类 / 标签 / 归档 / 关于），右侧搜索图标 + 主题切换按钮 + 登录头像。

2\. Hero 区（可选，方向 B 才有）：大标题 "分享技术，记录思考"，副标题一句话，背景紫色光晕。

3\. 文章列表：时间线式，每篇文章一个条目，包含：

&#x20;  \- 左侧：发布日期（小字，text-muted，格式 "2026-09-28"）

&#x20;  \- 右侧：分类标签（小胶囊，accent 色文字+透明底）+ 文章标题（text-primary，20px，font-weight 600，hover 时颜色变 accent）+ 摘要（text-body，14px，最多两行，超出省略）+ 元信息行（阅读时长 · 标签）

&#x20;  \- 条目之间用 1px border-bottom 分割，无卡片背景

&#x20;  \- hover 时整行背景轻微变色（bg-elevated）

4\. 分页：底部居中，"上一页 / 页码 / 下一页"，当前页 accent 色。

5\. 页脚：三列（关于本站 / 友情链接 / 联系方式），底部版权行，text-muted 小字。

交互：

\- 文章标题 hover：颜色过渡到 accent，左侧日期也同步高亮

\- 导航链接 hover：底部出现 2px accent 下划线，从左到右展开动画

\- 滚动到文章列表时，条目逐个淡入上移（stagger 50ms）

禁止：轮播图、侧边栏、彩色渐变背景、卡片阴影堆叠。
```

### 3.2 文章详情页（Article Detail）



```
为 Vue 3 技术博客设计文章详情页，采用【方向 A】。

布局：

1\. 顶部导航栏（同首页，滚动时收起为紧凑模式）。

2\. 文章头部（居中，max-width 720px）：

&#x20;  \- 分类标签 + 发布日期 + 阅读时长（一行，text-muted，14px）

&#x20;  \- 文章标题（serif 或 sans 加粗，36-48px，行高 1.25，letter-spacing -0.01em）

&#x20;  \- 作者信息行：头像（32px 圆形）+ 作者名 + 关注按钮

&#x20;  \- 封面图（可选，全宽出血，圆角 12px，margin-top 32px）

3\. 正文（居中，max-width 720px，font-size 16px，line-height 1.75）：

&#x20;  \- 段落 margin-bottom 1.5rem

&#x20;  \- h2：28px，margin-top 2.5rem，margin-bottom 1rem，font-weight 700

&#x20;  \- h3：22px，margin-top 2rem

&#x20;  \- 引用块：左侧 3px accent 竖线，padding-left 1rem，text-muted，斜体

&#x20;  \- 代码块：bg 比页面深一档，圆角 8px，padding 1.5rem，font-mono 14px，行号可选，右上角复制按钮

&#x20;  \- 行内代码：accent 色背景 10% 透明度，accent 色文字，padding 2px 6px，圆角 4px

&#x20;  \- 图片：全宽，圆角 8px，caption 居中 text-muted 14px

&#x20;  \- 表格：全宽，border-collapse，header 行 bg-elevated，行间 border-bottom

4\. 文章底部：

&#x20;  \- 标签行（标签胶囊）

&#x20;  \- 点赞/收藏按钮组（图标 + 数字，hover 放大 1.1x）

&#x20;  \- 作者介绍卡（头像 + 简介 + 关注按钮，bg-elevated，圆角 12px）

&#x20;  \- 相关文章推荐（2-3 篇横向卡片，含标题+分类+日期）

5\. 评论区：嵌套评论，每条含头像+用户名+时间+内容+回复按钮，输入框在顶部。

6\. 右侧悬浮目录（TOC，仅桌面端，position sticky，top 100px，宽度 220px）：

&#x20;  \- 自动提取 h2/h3，当前阅读段落高亮 accent 色

&#x20;  \- 滚动时平滑跟随

7\. 阅读进度条：页面顶部 2px 细线，随滚动填充 accent 色。

交互：

\- 选中文字时弹出 "复制 / 分享" 小浮层

\- 代码块复制成功后按钮变对勾 1.5 秒

\- 图片点击放大查看（lightbox）

禁止：正文区域出现广告横幅、侧边栏小部件、自动播放视频。
```

### 3.3 搜索页（Search）



```
为 Vue 3 技术博客设计搜索页，采用【方向 A】。

布局：

1\. 顶部导航栏。

2\. 搜索区（居中，max-width 640px，margin-top 80px）：

&#x20;  \- 大搜索框：高度 56px，圆角 12px，border 1px，左侧搜索图标，placeholder "搜索文章、标签、作者..."，输入时右侧出现清除按钮

&#x20;  \- 搜索框下方：热门搜索标签（横向排列，点击即搜）+ 搜索历史（可清除）

3\. 搜索结果（max-width 640px，margin-top 48px）：

&#x20;  \- 结果统计行："找到 N 篇相关文章"（text-muted，14px）

&#x20;  \- 结果列表：每篇含标题（关键词高亮 accent 色）+ 摘要（关键词高亮）+ 分类+日期+阅读时长

&#x20;  \- 无结果时：插画 + "没有找到相关内容" + 推荐热门文章

4\. 高级筛选（可选展开）：按分类筛选、按时间范围筛选、排序方式（相关度/最新/最热）。

交互：

\- 输入防抖 300ms，实时显示结果

\- 键盘 "/" 快捷键聚焦搜索框

\- 结果项 hover：背景 bg-elevated，标题变 accent

\- 支持 ↑↓ 键选择结果，Enter 跳转

禁止：搜索结果用卡片网格（列表更易扫读）、分页改为无限滚动。
```

### 3.4 归档 / 分类 / 标签页（Archive / Category / Tag）



```
为 Vue 3 技术博客设计归档页，采用【方向 A】。

布局（单栏居中，max-width 800px）：

1\. 顶部导航栏。

2\. 页面标题："归档" 或 "分类：前端" 或 "标签：Vue3"（36px，font-weight 700）。

3\. 统计行：共 N 篇文章（text-muted）。

4\. 按年份分组的时间线：

&#x20;  \- 年份标题（24px，font-weight 700，sticky 顶部，背景与页面同色，padding-top 32px）

&#x20;  \- 该年份下的文章列表：每行 = 日期（MM-DD，text-muted，等宽字体）+ 标题（hover accent）+ 分类标签

&#x20;  \- 月份之间可用细分割线

5\. 标签云页替代方案：所有标签按文章数量排序，字号映射数量（越多越大），hover 显示文章数。

交互：

\- 年份标题滚动时吸顶

\- 点击标签跳转到标签页

\- 回到顶部按钮（滚动超过 500px 时出现，右下角圆形按钮）

禁止：用日历组件展示归档（日历只适合每日更新的博客，技术博客是不定期更新）。
```

### 3.5 关于页（About）



```
为 Vue 3 技术博客设计关于页，采用【方向 A】。

布局（单栏居中，max-width 720px）：

1\. 顶部导航栏。

2\. 个人介绍区：

&#x20;  \- 头像（120px 圆形，带 subtle 边框）

&#x20;  \- 姓名（32px，font-weight 700）

&#x20;  \- 一句话简介（text-body，18px）

&#x20;  \- 社交链接行（GitHub / Twitter / 邮箱 / RSS，图标 + 文字，hover accent）

3\. 关于本站（正文排版，同文章详情页的正文样式）：

&#x20;  \- 建站初衷、技术栈、更新频率

&#x20;  \- 可用时间线展示建站历程（2024 建站 → 2025 重构 → 2026 UI 改版）

4\. 技能/兴趣标签云。

5\. 友情链接区：好友博客卡片网格（头像 + 站名 + 一句话描述，hover 边框 accent）。

交互：

\- 社交图标 hover：translateY(-2px) + 颜色过渡

\- 时间线节点 hover：节点放大

禁止：放简历式的长表格、过多的自拍照。
```

### 3.6 管理后台（Dashboard）



```
为 Vue 3 技术博客设计管理后台，采用【方向 A 的深色变体】。

布局（经典后台三段式，但视觉现代化）：

1\. 左侧边栏（width 240px，bg-elevated，border-right）：

&#x20;  \- 顶部：Logo + 站点名

&#x20;  \- 导航分组：内容管理（文章列表 / 写文章 / 评论 / 图片）、用户管理（用户列表）、系统（配置 / 友链 / 广告 / 反馈 / 登录日志）

&#x20;  \- 每个菜单项：图标 + 文字，选中态左侧 3px accent 竖条 + 文字 accent + 背景 accent 5% 透明度

&#x20;  \- 可折叠为图标模式（width 64px）

2\. 顶部栏（height 56px，bg-elevated，border-bottom）：

&#x20;  \- 左侧：折叠按钮 + 面包屑

&#x20;  \- 右侧：搜索 + 通知铃铛 + 用户头像下拉

3\. 主内容区（padding 24px，bg 页面底色）：

&#x20;  \- 数据概览首页：4 个统计卡片（文章数 / 评论数 / 用户数 / 浏览量），每卡含图标 + 数字 + 环比小箭头

&#x20;  \- 图表区：ECharts 浏览量趋势（折线，accent 色，area 渐变透明）+ 分类分布（环形图）

&#x20;  \- 文章列表页：表格（无竖线，仅横线 border-bottom，hover 行 bg-elevated），操作列用文字按钮（编辑/删除，删除用红色）

&#x20;  \- 写文章页：左右分栏（左 md-editor-v3 编辑器，右实时预览），顶部标题输入框（无边框，28px，placeholder "输入文章标题"），底部发布栏（分类选择 / 标签输入 / 封面上传 / 保存草稿 / 发布按钮）

4\. 通用组件规范：

&#x20;  \- 按钮：主按钮 accent 实心白字，次按钮 border 透明底，危险按钮红色

&#x20;  \- 表单：label 14px text-muted，input height 40px 圆角 8px border，focus 时 border accent + 2px accent 15% 外发光

&#x20;  \- 表格分页：右下角，简洁样式

&#x20;  \- 弹窗：圆角 12px，标题 18px，底部操作按钮右对齐

交互：

\- 侧边栏折叠动画 200ms ease

\- 表格行 hover 背景过渡

\- 发布成功后 toast 提示（右上角，3 秒自动消失）

\- 表单校验错误：input border 红色 + 下方红色小字提示

禁止：用 Element Plus 默认主题色（蓝绿）不做定制、表格用斑马纹、按钮加渐变。
```



***

## 四、技术落地建议

### 4.1 替换 UI 库

现有 Element Plus 组件风格偏 "中后台企业风"，与内容站气质不符。建议：



* **前台**：完全不用组件库，手写 CSS（用 CSS Variables 实现主题切换）。组件少（按钮 / 输入框 / 卡片 / 标签），手写成本低且风格可控。

* **后台**：保留 Element Plus 但深度定制主题（覆盖 CSS 变量），或迁移到 `Naive UI`（更现代、主题定制更灵活）。

### 4.2 推荐依赖



```
{

&#x20; "lucide-vue-next": "^0.400.0",   // 统一图标

&#x20; "vueuse": "^10.0.0",             // useDark / useToggle 主题切换

&#x20; "pangu": "^4.0.0",               // 中英文自动加空格

&#x20; "highlight.js": "^11.0.0"        // 代码高亮（替代 md-editor 自带的）

}
```

### 4.3 主题切换实现原理

用 `data-theme` 属性挂在 `<html>` 上，所有颜色走 CSS Variable：



```
:root { --bg: #FAFAFA; --text: #18181B; }

\[data-theme="dark"] { --bg: #09090B; --text: #FAFAFA; }
```

初始化时读 `localStorage`，没有则跟随 `prefers-color-scheme`，在 `<head>` 内联一段 JS 避免闪烁（FOUC）。

### 4.4 迁移步骤



1. 先建 `design-tokens.css`，定义所有变量（三个方向选一套）。

2. 重构全局样式（reset、字体、滚动条、选中色）。

3. 从文章详情页开始重构（内容页是博客的核心，做好了其他页有参照）。

4. 再做首页 → 搜索 → 归档 → 关于。

5. 最后重构管理后台。

6. 每完成一页用 Lighthouse 跑性能和可访问性。



***

## 五、一页速查卡



| 维度   | 方向 A 极简编辑 | 方向 B 深色精致    | 方向 C 杂志内容     |
| ---- | --------- | ------------ | ------------- |
| 底色   | #FAFAFA   | #0A0A0A      | #FDFCFB 暖白    |
| 强调色  | #2563EB 蓝 | #A78BFA 紫    | #B45309 琥珀    |
| 标题字体 | Sans 加粗   | Sans 负字距     | Serif         |
| 首页布局 | 单栏时间线     | 单栏 + 光晕 Hero | Featured + 网格 |
| 文章侧栏 | 无（仅 TOC）  | 无（仅 TOC）     | 有（TOC + 作者）   |
| 轮播图  | 无         | 无            | 无             |
| 适合场景 | 技术博客首选    | 极客个人品牌       | 内容量大          |

> **最终建议**
>
> ：你的项目是技术博客，当前已有紫色基因，推荐 
>
> **方向 A（浅色为主 + 深色模式）**
>
>  作为主方案，把紫色作为 accent 色保留。如果想更有辨识度，选 
>
> **方向 B**
>
> 。方向 C 适合文章数超过 200 篇后再考虑。