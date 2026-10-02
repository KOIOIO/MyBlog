<template>
  <div class="article">
    <web-navbar :noScroll="true"/>
    <div class="page">
      <div class="container">
        <main class="main-content">
          <!-- 文章头部 -->
          <header class="article-header">
            <div v-if="articleInfo.cover" class="article-cover">
              <img :src="articleInfo.cover" :alt="articleInfo.title" loading="lazy"/>
            </div>
            <div class="meta-line">
              <span v-if="articleInfo.category" class="meta-category">{{ categoryLabel(articleInfo.category) }}</span>
              <span class="meta-sep">·</span>
              <span>{{ articleInfo.created_at }}</span>
            </div>
            <h1 class="article-title">{{ articleInfo.title }}</h1>
            <p v-if="articleInfo.abstract" class="article-abstract">{{ articleInfo.abstract }}</p>
          </header>

          <!-- 正文 -->
          <div class="article-body">
            <MdPreview :id="mdID" :modelValue="articleInfo.content"/>
          </div>

          <!-- 标签 + 互动 -->
          <div class="article-footer">
            <div class="tag-row">
              <el-tag v-for="item in articleInfo.tags" :key="item" effect="plain" round>{{ tagLabel(item) }}</el-tag>
            </div>
            <div class="action-row">
              <span class="action-stat">
                <el-icon><component is="View"/></el-icon> {{ articleInfo.views }}
              </span>
              <span class="action-stat">
                <el-icon><component is="ChatDotRound"/></el-icon> {{ articleInfo.comments }}
              </span>
              <span class="action-stat action-like" @click="handelLike">
                <el-icon :key="isLike" :class="{ 'like-pop': isLike }">
                  <component v-if="!isLike" is="Star"/>
                  <component v-else is="StarFilled"/>
                </el-icon>
                {{ articleInfo.likes }}
              </span>
            </div>
          </div>

          <!-- 评论输入 -->
          <div class="comment-box">
            <div class="comment-head">
              <h3 class="comment-title">{{ t('pages.article.commentTitle') }}</h3>
              <span class="comment-count">{{ comments.length }}</span>
            </div>
            <div class="comment-input-wrap">
              <el-input v-model="content" :autosize="{ minRows: 3, maxRows: 8 }" type="textarea"
                        :placeholder="t('pages.article.commentPlaceholder')" maxlength="320"/>
              <div class="comment-actions">
                <el-popover width="448" trigger="click">
                  <template #reference>
                    <el-avatar class="emoji-trigger">😊</el-avatar>
                  </template>
                  <template #default>
                    <div class="emoji-panel">
                      <span
                          v-for="emoji in emojis"
                          :key="emoji"
                          class="emoji-item"
                          @click="content=content+emoji"
                      >{{ emoji }}</span>
                    </div>
                  </template>
                </el-popover>
                <div class="button-group">
                  <el-button size="small" type="primary" @click="submitComment">{{ t('pages.article.postComment') }}</el-button>
                  <el-button size="small" @click="content=''">{{ t('common.cancel') }}</el-button>
                </div>
              </div>
            </div>
            <el-text class="login-tip">{{ t('pages.article.loginTip') }}</el-text>
          </div>

          <!-- 评论列表 -->
          <div class="comment-list">
            <comment-item :comments="comments"/>
          </div>
        </main>

        <!-- 右侧目录 -->
        <aside class="aside-content">
          <div ref="catalogRef" class="catalog" @click="onCatalogClick">
            <div class="catalog-title">{{ t('pages.article.catalogTitle') }}</div>
            <MdCatalog :editorId="mdID" :scrollElement="scrollElement" :offsetTop="100" :scrollElementOffsetTop="80"/>
          </div>
          <el-anchor class="comment-link" :marker="false">
            <el-anchor-link href="#comment">{{ t('pages.article.backToComment') }}</el-anchor-link>
          </el-anchor>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useRoute} from "vue-router";
import {type Article, articleInfoByID} from "@/api/article";
import router from "@/router";
import {computed, onMounted, onUnmounted, ref, watch, nextTick} from "vue";
import {MdPreview, MdCatalog} from 'md-editor-v3';
import 'md-editor-v3/lib/style.css';
import 'md-editor-v3/lib/preview.css';
import WebNavbar from "@/components/layout/WebNavbar.vue";
import CommentItem from "@/components/common/CommentItem.vue";
import {articleIsLike, articleLike, type ArticleLikeRequest} from "@/api/article";
import {type Comment, commentCreate, type CommentCreateRequest, commentInfoByArticleID} from "@/api/comment";
import {useLayoutStore} from "@/stores/layout";
import {useUserStore} from "@/stores/user";
import {useI18n} from "vue-i18n";
import {categoryLabel, tagLabel} from "@/i18n/meta";

const {t} = useI18n();

const mdID = "md-id"

const articleInfo = ref<Article>({
  created_at: '',
  updated_at: '',
  cover: '',
  title: '',
  keyword: '',
  category: '',
  tags: [],
  abstract: '',
  content: '',
  comments: 0,
  views: 0,
  likes: 0,
  is_top: 0,
})

const scrollElement = document.documentElement

const route = useRoute()

const articleID = computed(() => route.params.id)

const getArticleInfo = async () => {
  const res = await articleInfoByID(articleID.value as string)
  if (res.code === 0) {
    articleInfo.value = res.data
  } else {
    await router.push({name: "404"})
  }
}

getArticleInfo()

const isLike = ref(false)

const getIsLikeInfo = async () => {
  const req: ArticleLikeRequest = {
    article_id: articleID.value as string
  }
  const res = await articleIsLike(req)
  if (res.code === 0) {
    isLike.value = res.data
  }
}

if (useUserStore().state.userInfo.role_id !== 0) {
  getIsLikeInfo()
}

const handelLike = async () => {
  const req: ArticleLikeRequest = {
    article_id: articleID.value as string
  }
  const res = await articleLike(req)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    articleInfo.value.likes += isLike.value ? -1 : 1
    isLike.value = !isLike.value
  }
}

const content = ref('')

const emojis = ['😀','😁','😂','🤣','😊','😇','🙂','😉','😍','🥰','😘','😋','😛','🤪','🤗','🤭','🤔','😏','🙄','😬','😮','😲','🥺','😢','😭','😤','😡','🤯','😱','😰','🥳','😎','🤓','👍','👎','👏','🙏','💪','🔥','❤️'];

const submitComment = async () => {
  const commentCreateRequest: CommentCreateRequest = {
    article_id: articleID.value as string,
    p_id: null,
    content: content.value,
  }
  const res = await commentCreate(commentCreateRequest)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    content.value = ''
    layoutStore.state.shouldRefreshCommentList = true
  }
}

const comments = ref<Comment[]>([])

const getArticleCommentsInfo = async () => {
  comments.value = []
  const res = await commentInfoByArticleID(articleID.value as string)
  if (res.code === 0) {
    comments.value = res.data
  }
}

const layoutStore = useLayoutStore()
watch(() => layoutStore.state.shouldRefreshCommentList, (newVal) => {
  if (newVal) {
    getArticleCommentsInfo()
    layoutStore.state.shouldRefreshCommentList = false
  }
})

/* ============================================================
   §6 目录栏重构：连续阅读进度指示条 + 自动居中 + 点击反馈 + 首载 stagger
   - 保留 MdCatalog 承担锚点 / active 文字判定；
   - 滚动(rAF)时测量正文标题与目录项位置，连续插值驱动 .md-editor-catalog-indicator 的 top；
   - 若标题无 id 或 DOM 测量不可行，自动降级（不挂指示条，MdCatalog 仍高亮文字）。
   ============================================================ */
const catalogRef = ref<HTMLElement | null>(null)

let tocHeadings: HTMLElement[] = []
let tocLinks: HTMLElement[] = []
let tocIndicator: HTMLElement | null = null
let tocReady = false
let lastActiveIdx = -1
let tocRaf = 0
let tocPending = false
let resizeObs: ResizeObserver | null = null

const READ_LINE = 100 // 阅读线距视口顶部（navbar 56 + 余量）

const measureDocTops = () =>
  tocHeadings.map((h) => h.getBoundingClientRect().top + window.scrollY)

const tocUpdate = () => {
  tocPending = false
  if (!tocReady || !tocIndicator || tocHeadings.length === 0 || tocLinks.length === 0) return
  const scrollY = window.scrollY || document.documentElement.scrollTop || 0
  const tops = measureDocTops()

  // active 索引：最后一个 docTop <= 阅读线的标题
  let i = -1
  for (let k = 0; k < tops.length; k++) {
    if (tops[k] <= scrollY + READ_LINE) i = k
    else break
  }
  if (i < 0) i = 0

  // 区间进度 p（在标题 i 与 i+1 之间线性插值）
  let p = 0
  if (i >= tops.length - 1) {
    p = 1
  } else {
    const span = tops[i + 1] - tops[i]
    p = span > 0 ? (scrollY + READ_LINE - tops[i]) / span : 0
    p = Math.min(1, Math.max(0, p))
  }

  // 指示条 top = lerp(目录项 i, 目录项 i+1, p)
  const topA = tocLinks[i]?.offsetTop ?? 0
  const topB = tocLinks[Math.min(i + 1, tocLinks.length - 1)]?.offsetTop ?? topA
  tocIndicator.style.top = `${topA + (topB - topA) * p}px`

  // 自动居中跟随（仅 active 变化时，手动滚动目录容器，避免整页滚动）
  if (i !== lastActiveIdx) {
    lastActiveIdx = i
    const link = tocLinks[i]
    const box = catalogRef.value
    if (link && box && box.scrollHeight > box.clientHeight) {
      const wantTop = link.offsetTop
      const wantBottom = wantTop + link.offsetHeight
      const viewTop = box.scrollTop
      const viewBottom = viewTop + box.clientHeight
      if (wantTop < viewTop) box.scrollTop = wantTop
      else if (wantBottom > viewBottom) box.scrollTop = wantBottom - link.offsetHeight
    }
  }
}

const onScroll = () => {
  if (tocPending) return
  tocPending = true
  tocRaf = requestAnimationFrame(tocUpdate)
}

const initToc = () => {
  const box = catalogRef.value
  const preview = document.querySelector('.md-editor-preview-wrapper')
  if (!box || !preview) return

  const headings = Array.from(preview.querySelectorAll('h2[id], h3[id]')) as HTMLElement[]
  const links = Array.from(box.querySelectorAll('.md-editor-catalog-link')) as HTMLElement[]
  const indicator = box.querySelector('.md-editor-catalog-indicator') as HTMLElement | null

  // 降级：标题无 id / 目录项缺失 → 不挂连续指示条（MdCatalog 仍负责 active 文字）
  if (headings.length === 0 || links.length === 0 || !indicator) {
    tocReady = false
    return
  }

  tocHeadings = headings
  tocLinks = links
  tocIndicator = indicator
  tocReady = true
  lastActiveIdx = -1

  // 首载 stagger：目录项 60ms 间隔 kf-fade-up 进入一次
  links.forEach((link, idx) => {
    link.classList.add('toc-enter')
    ;(link as HTMLElement).style.animationDelay = `${idx * 60}ms`
  })

  tocUpdate()
  window.addEventListener('scroll', onScroll, {passive: true})
  window.addEventListener('resize', onScroll, {passive: true})
  if (!resizeObs) {
    resizeObs = new ResizeObserver(() => onScroll())
    resizeObs.observe(preview)
  }
}

// 点击目录项 → 被点项播 kf-toc-flash 300ms（事件委托）
const onCatalogClick = (e: MouseEvent) => {
  const target = (e.target as HTMLElement).closest('.md-editor-catalog-link') as HTMLElement | null
  if (!target) return
  target.classList.remove('toc-flash')
  // 强制重排以重放动画
  void target.offsetWidth
  target.classList.add('toc-flash')
  setTimeout(() => target.classList.remove('toc-flash'), 320)
}

onMounted(() => {
  getArticleCommentsInfo()
  // 正文异步渲染后 MdCatalog 才生成，重试直到目录 DOM 就绪
  nextTick(() => {
    initToc()
    setTimeout(initToc, 300)
  })
})

watch(() => articleInfo.value.content, () => {
  nextTick(() => {
    initToc()
    setTimeout(initToc, 300)
  })
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  window.removeEventListener('resize', onScroll)
  if (tocRaf) cancelAnimationFrame(tocRaf)
  resizeObs?.disconnect()
})
</script>

<style scoped lang="scss">
.article {
  background-color: var(--bg);
  min-height: 100vh;

  .page {
    padding: calc(70px + var(--sp-6)) var(--sp-4) var(--sp-9);
  }

  .container {
    display: flex;
    justify-content: center;
    gap: var(--sp-7);
    max-width: 1000px;
    margin: 0 auto;
  }

  .main-content {
    flex: 0 1 var(--reading-width);
    min-width: 0;
  }

  /* 文章头部 */
  .article-cover {
    border-radius: var(--radius-md);
    overflow: hidden;
    margin-bottom: var(--sp-8);
    border: 1px solid var(--border);
  }

  .article-cover img {
    width: 100%;
    height: auto;
    display: block;
  }

  .article-header {
    margin-bottom: var(--sp-7);

    .meta-line {
      font-size: var(--fs-14);
      color: var(--text-muted);
      margin-bottom: var(--sp-3);
      display: flex;
      align-items: center;
      gap: var(--sp-2);

      .meta-category {
        color: var(--accent);
      }

      .meta-sep {
        color: var(--border);
      }
    }

    .article-title {
      font-family: var(--font-serif);
      font-size: var(--fs-36);
      font-weight: 700;
      line-height: var(--lh-title);
      color: var(--text-primary);
      margin: 0 0 var(--sp-4);
      letter-spacing: -0.01em;
    }

    .article-abstract {
      font-size: var(--fs-16);
      color: var(--text-body);
      line-height: var(--lh-body);
      margin: 0;
    }
  }

  /* 正文排版 */
  .article-body {
    font-size: var(--fs-16);
    line-height: var(--lh-body);
    color: var(--text-body);

    :deep(p) {
      margin-bottom: 1.5rem;
    }

    :deep(h2) {
      font-size: var(--fs-24);
      font-weight: 700;
      margin-top: 2.5rem;
      margin-bottom: 1rem;
      color: var(--text-primary);
      scroll-margin-top: 80px;
    }

    :deep(h3) {
      font-size: var(--fs-20);
      font-weight: 600;
      margin-top: 2rem;
      margin-bottom: var(--sp-3);
      color: var(--text-primary);
      scroll-margin-top: 80px;
    }

    :deep(blockquote) {
      border-left: 3px solid var(--accent);
      padding-left: var(--sp-4);
      margin-left: 0;
      color: var(--text-muted);
      font-style: italic;
    }

    :deep(code) {
      background: var(--accent-weak);
      color: var(--accent);
      padding: 2px 6px;
      border-radius: 4px;
      font-family: var(--font-mono);
      font-size: 0.9em;
    }

    :deep(pre) {
      background: var(--bg-elevated);
      border: 1px solid var(--border);
      border-radius: 8px;
      padding: var(--sp-4);
      overflow-x: auto;

      code {
        background: transparent;
        color: var(--text-body);
        padding: 0;
      }
    }

    :deep(img) {
      max-width: 100%;
      border-radius: 8px;
    }
  }

  /* 标签 + 互动 */
  .article-footer {
    margin-top: var(--sp-7);
    padding-top: var(--sp-5);
    border-top: 1px solid var(--border);

    .tag-row {
      display: flex;
      flex-wrap: wrap;
      gap: var(--sp-2);
      margin-bottom: var(--sp-4);
    }

    .action-row {
      display: flex;
      align-items: center;
      gap: var(--sp-5);
      font-size: var(--fs-14);
      color: var(--text-muted);

      .action-stat {
        display: inline-flex;
        align-items: center;
        gap: var(--sp-1);
      }

      .action-like {
        cursor: pointer;
        transition: color 150ms ease-out;

        &:hover {
          color: var(--accent);
        }

        :deep(.el-icon.like-pop) {
          animation: kf-pop 150ms var(--ease-pop);
        }
      }
    }
  }

  /* 评论（B站式容器） */
  .comment-box {
    margin-top: var(--sp-7);
    padding-top: var(--sp-5);
    border-top: 1px solid var(--border);

    .comment-head {
      display: flex;
      align-items: baseline;
      gap: var(--sp-1);
      margin-bottom: var(--sp-3);

      .comment-title {
        font-size: var(--fs-18);
        font-weight: 600;
        color: var(--text-primary);
      }

      .comment-count {
        font-size: var(--fs-14);
        color: var(--text-muted);
      }
    }

    .comment-input-wrap {
      background: var(--bg-elevated);
      border: 1px solid var(--border);
      border-radius: var(--radius-sm);
      padding: var(--sp-3);
      transition: border-color 200ms var(--ease-out);

      &:focus-within {
        border-color: var(--accent);
      }

      :deep(.el-textarea__inner) {
        background: transparent;
        box-shadow: none !important;
      }

      .comment-actions {
        display: flex;
        align-items: center;
        margin-top: var(--sp-2);

        .emoji-trigger {
          background: transparent;
          font-size: var(--fs-20);
          cursor: pointer;
        }

        .emoji-panel {
          display: flex;
          flex-wrap: wrap;
          gap: var(--sp-1);
          padding: var(--sp-2);

          .emoji-item {
            width: 34px;
            height: 34px;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            font-size: var(--fs-20);
            border-radius: 6px;
            cursor: pointer;
            transition: background-color 120ms ease;

            &:hover {
              background-color: var(--bg-elevated);
            }
          }
        }

        .button-group {
          margin-left: auto;
        }
      }
    }

    .login-tip {
      display: block;
      margin-top: var(--sp-2);
      font-size: var(--fs-12);
      color: var(--text-muted);
    }
  }

  .comment-list {
    margin-top: var(--sp-5);
  }

  /* 右侧目录 */
  .aside-content {
    width: 220px;
    flex-shrink: 0;
    position: sticky;
    top: 100px;
    align-self: flex-start;

    .catalog {
      position: relative;
      max-height: calc(100vh - 220px);
      overflow-y: auto;
      padding-bottom: var(--sp-4);
      border-bottom: 1px solid var(--border);

      .catalog-title {
        font-size: var(--fs-12);
        font-weight: 500;
        color: var(--text-muted);
        text-transform: uppercase;
        letter-spacing: 0.05em;
        margin-bottom: var(--sp-3);
      }

      /* §6 active 项：accent 文字 + 700，不做整行背景块 */
      :deep(.md-editor-catalog-active > span) {
        color: var(--accent);
        font-weight: 700;
        transition: color 150ms var(--ease-out);
      }

      /* §6 连续指示条：accent、2px 宽，top 150ms 区间内平滑 */
      :deep(.md-editor-catalog-indicator) {
        background-color: var(--accent);
        width: 2px;
        height: 18px;
        border-radius: 2px;
        transition: top 150ms var(--ease-out);
      }

      /* §6 hover：文字 accent + 右侧「跳转」箭头滑入 */
      :deep(.md-editor-catalog-link) {
        position: relative;

        > span {
          padding-right: 18px;
          transition: color 150ms var(--ease-out);

          &::after {
            content: "→";
            position: absolute;
            right: 2px;
            top: 50%;
            transform: translateY(-50%) translateX(2px);
            opacity: 0;
            font-size: 12px;
            color: var(--accent);
            transition: opacity 150ms var(--ease-out),
              transform 150ms var(--ease-out);
          }
        }

        &:hover > span {
          color: var(--accent);

          &::after {
            opacity: 1;
            transform: translateY(-50%) translateX(0);
          }
        }
      }

      /* §6 首载 stagger 进入 */
      :deep(.md-editor-catalog-link.toc-enter) {
        animation: kf-fade-up 150ms var(--ease-out) both;
      }

      /* §6 点击确认反馈 */
      :deep(.md-editor-catalog-link.toc-flash > span) {
        animation: kf-toc-flash 300ms var(--ease-out);
      }
    }

    .comment-link {
      margin-top: var(--sp-4);
    }
  }
}

:deep(.el-popover.el-popper) {
  .el-image {
    height: 50px;
    width: 50px;
  }
}
</style>
