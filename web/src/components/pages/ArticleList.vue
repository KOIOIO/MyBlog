<template>
  <div class="article-list">
    <!-- 筛选工具栏 -->
    <div class="toolbar">
      <div class="toolbar-search">
        <el-input v-model="articleSearchRequest.query" :placeholder="t('components.articleList.searchPlaceholder')" prefix-icon="Search" maxlength="50"
                  clearable
                  @change="changeArticleSearchItem"/>
      </div>
      <div class="toolbar-row">
        <span class="toolbar-label">{{ t('components.articleList.categoryLabel') }}</span>
        <div class="chip-group">
          <button type="button" class="chip" :class="{ 'is-selected': articleSearchRequest.category === ''}"
                  @click="selectCategory('')">{{ t('common.all') }}</button>
          <button v-for="item in categoryArr" :key="item" type="button" class="chip"
                  :class="{ 'is-selected': articleSearchRequest.category === item}"
                  @click="selectCategory(item)">{{ categoryLabel(item) }}</button>
        </div>
      </div>
      <div class="toolbar-row">
        <span class="toolbar-label">{{ t('components.articleList.tagLabel') }}</span>
        <div class="chip-group">
          <button type="button" class="chip" :class="{ 'is-selected': articleSearchRequest.tag === ''}"
                  @click="selectTag('')">{{ t('common.all') }}</button>
          <button v-for="item in tagArr" :key="item" type="button" class="chip"
                  :class="{ 'is-selected': articleSearchRequest.tag === item}"
                  @click="selectTag(item)">{{ tagLabel(item) }}</button>
        </div>
      </div>
      <div class="toolbar-row">
        <span class="toolbar-label">{{ t('components.articleList.sortLabel') }}</span>
        <button type="button" class="sort-btn press" @click="handleSortClick();changeArticleSearchItem()">
          <el-icon :color="downColor"><component is="SortDown"/></el-icon>
          <el-icon :color="upColor"><component is="SortUp"/></el-icon>
        </button>
        <div class="chip-group">
          <button v-for="item in sortArr" :key="item.value" type="button" class="chip"
                  :class="{ 'is-selected': articleSearchRequest.sort === item.value}"
                  @click="selectSort(item.value)">{{ item.label }}</button>
        </div>
      </div>
    </div>

    <!-- 时间线列表（仅筛选组合变化时淡入；key 不含页码 → 分页即时） -->
    <transition name="filter-fade">
      <div class="timeline" :key="filterKey">
        <div v-for="row in articleTableData" :key="row._id" class="timeline-item"
             @click="handleArticleJumps(row._id)">
          <div class="item-date">
            <span class="date-text">{{ row._source.created_at }}</span>
            <span v-if="row._source.category" class="item-category">{{ categoryLabel(row._source.category) }}</span>
          </div>
          <div class="item-main">
            <div class="item-body">
              <h2 class="item-title">
                <span v-if="row._source.is_top === 1" class="item-top-badge">{{ t('pages.article.top') }}</span>
                {{ row._source.title }}
              </h2>
              <p class="item-abstract">{{ row._source.abstract }}</p>
              <div class="item-meta">
                <span v-for="tag in row._source.tags" :key="tag" class="meta-tag">{{ tagLabel(tag) }}</span>
                <span class="meta-sep">·</span>
                <span class="meta-stat">
                  <el-icon><component is="View"/></el-icon> {{ row._source.views }}
                </span>
                <span class="meta-stat">
                  <el-icon><component is="ChatDotRound"/></el-icon> {{ row._source.comments }}
                </span>
                <span class="meta-stat">
                  <el-icon><component is="Star"/></el-icon> {{ row._source.likes }}
                </span>
              </div>
            </div>
            <div v-if="row._source.cover" class="item-cover">
              <img :src="row._source.cover" :alt="row._source.title" loading="lazy"/>
            </div>
          </div>
        </div>
      </div>
    </transition>

    <!-- 分页 -->
    <div class="pagination-wrap">
      <el-pagination
          :current-page="page"
          :page-size="page_size"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type {Hit} from "@/api/common";
import {type Article, articleCategory, articleSearch, type ArticleSearchRequest, articleTags} from "@/api/article";
import {computed, reactive, ref} from "vue";
import {useI18n} from "vue-i18n";
import {categoryLabel, tagLabel} from "@/i18n/meta";

const {t} = useI18n();

const articleSearchRequest = reactive<ArticleSearchRequest>({
  query: "",
  category: "",
  tag: "",
  sort: "",
  order: "desc",
  page: 1,
  page_size: 10,
})

const categoryArr = ref<string[]>([])
const tagArr = ref<string[]>([])
const sortArr = computed(() => [
  {label: t('components.articleList.sortDefault'), value: ""},
  {label: t('components.articleList.sortTime'), value: "time"},
  {label: t('components.articleList.sortComment'), value: "comment"},
  {label: t('components.articleList.sortView'), value: "view"},
  {label: t('components.articleList.sortLike'), value: "like"},
])

const downColor = computed(() => {
  return articleSearchRequest.order === "desc" ? "blue" : "gray"
})
const upColor = computed(() => {
  return articleSearchRequest.order === "desc" ? "gray" : "blue"
})

const handleSortClick = () => {
  articleSearchRequest.order = articleSearchRequest.order === "desc" ? "asc" : "desc"
}

// chip 单选互斥：手工判断，仅当值真正变化时才触发请求（对齐原 el-radio-group 语义）
const selectCategory = (val: string) => {
  if (articleSearchRequest.category === val) return
  articleSearchRequest.category = val
  changeArticleSearchItem()
}
const selectTag = (val: string) => {
  if (articleSearchRequest.tag === val) return
  articleSearchRequest.tag = val
  changeArticleSearchItem()
}
const selectSort = (val: string) => {
  if (articleSearchRequest.sort === val) return
  articleSearchRequest.sort = val
  changeArticleSearchItem()
}

// 筛选组合 key（不含页码）：仅筛选变化触发 transition 淡入，分页切换保持即时
const filterKey = computed(() =>
  [articleSearchRequest.category, articleSearchRequest.tag, articleSearchRequest.sort, articleSearchRequest.order].join("|")
)

const getArticleCategory = async () => {
  const res = await articleCategory()
  if (res.code === 0) {
    res.data.forEach((item) => {
      categoryArr.value.push(item.category)
    })
  }
}

getArticleCategory()

const getArticleTags = async () => {
  const res = await articleTags()
  if (res.code === 0) {
    res.data.forEach((item) => {
      tagArr.value.push(item.tag)
    })
  }
}

getArticleTags()

const page = ref(1)
const page_size = ref(10)
const total = ref(0)
const articleTableData = ref<Hit<Article>[]>()

const getArticleSearchTableData = async () => {
  articleSearchRequest.page = page.value;
  articleSearchRequest.page_size = page_size.value;

  const table = await articleSearch(articleSearchRequest)

  if (table.code === 0) {
    articleTableData.value = table.data.list;
    total.value = table.data.total;
  }
}

getArticleSearchTableData()

const changeArticleSearchItem = () => {
  getArticleSearchTableData()
}

const handleArticleJumps = (id: string) => {
  window.open("/article/" + id)
}

const handleSizeChange = (val: number) => {
  page_size.value = val
  getArticleSearchTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getArticleSearchTableData()
}
</script>

<style scoped lang="scss">
.article-list {
  .toolbar {
    margin-bottom: var(--sp-6);

    .toolbar-search {
      max-width: 360px;
      margin-bottom: var(--sp-4);
    }

    .toolbar-row {
      display: flex;
      align-items: center;
      gap: var(--sp-3);
      margin-bottom: var(--sp-2);

      .toolbar-label {
        font-size: var(--fs-12);
        color: var(--text-muted);
        min-width: 32px;
        flex-shrink: 0;
      }

      .sort-btn {
        width: 32px;
        padding: 0;
        border: none;
        background: transparent;
      }
    }

    /* chip 组布局：idle/hover/选中(is-selected)/圆点/pop/press/focus 全部由全局 .chip 提供 */
    .chip-group {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: var(--sp-2);
    }
  }


  /* 筛选切换内容过渡：仅 opacity 100ms 入场（§7/§10；key 不含页码 → 分页即时） */
  .filter-fade-enter-active {
    transition: opacity var(--dur-base) var(--ease-out);
  }
  .filter-fade-leave-active {
    transition: opacity var(--dur-fast) var(--ease-in);
  }
  .filter-fade-enter-from,
  .filter-fade-leave-to {
    opacity: 0;
  }

  .timeline {
    border-top: 1px solid var(--border);

    .item-main {
      display: flex;
      align-items: flex-start;
      gap: var(--sp-6);
    }

    .item-body {
      flex: 1;
      min-width: 0;
    }

    .item-cover {
      flex-shrink: 0;
      width: 168px;
      margin-top: var(--sp-1);
      border-radius: var(--radius-md);
      overflow: hidden;
    }

    .item-cover img {
      width: 100%;
      height: auto;
      aspect-ratio: 16 / 9;
      object-fit: cover;
      border-radius: var(--radius-md);
      border: 1px solid var(--border);
      display: block;
      transition: transform var(--dur-mid) var(--ease-out);
    }

    .timeline-item:hover .item-cover img {
      transform: scale(1.02);
    }

    @media (max-width: 640px) {
      .item-cover { display: none; }
    }

    .timeline-item {
      padding: var(--sp-5) var(--sp-3);
      border-bottom: 1px solid var(--border);
      cursor: pointer;
      transition: background-color var(--dur-base) var(--ease-out), transform var(--dur-instant) var(--ease-out);

      &:hover {
        background-color: var(--bg-elevated);
      }

      &:active {
        transform: translateY(1px);
      }

      .item-date {
        display: flex;
        align-items: center;
        gap: var(--sp-3);
        margin-bottom: var(--sp-1);

        .date-text {
          font-size: var(--fs-12);
          color: var(--text-muted);
          font-family: var(--font-mono);
          white-space: nowrap;
        }

        .item-category {
          font-size: var(--fs-12);
          color: var(--accent);
          background: var(--accent-weak);
          padding: 1px var(--sp-2);
          border-radius: var(--radius-sm);
        }
      }

      .item-body {
        min-width: 0;

        .item-title {
          font-size: var(--fs-20);
          font-weight: 600;
          color: var(--text-primary);
          line-height: var(--lh-title);
          margin: 0 0 var(--sp-1);
          transition: color var(--dur-base) var(--ease-out);
        }

        .item-top-badge {
          display: inline-block;
          margin-right: var(--sp-2);
          padding: 0 var(--sp-1);
          font-size: var(--fs-12);
          font-weight: 600;
          line-height: 1.5;
          color: var(--accent);
          background: var(--accent-weak);
          border-radius: 6px;
          vertical-align: 3px;
        }

        &:hover .item-title {
          color: var(--accent);
        }

        .item-abstract {
          font-size: var(--fs-14);
          color: var(--text-body);
          line-height: var(--lh-body);
          margin: 0 0 var(--sp-2);
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          overflow: hidden;
        }

        .item-meta {
          display: flex;
          align-items: center;
          flex-wrap: wrap;
          gap: var(--sp-2);
          font-size: var(--fs-12);
          color: var(--text-muted);

          .meta-tag {
            color: var(--text-muted);
          }

          .meta-sep {
            color: var(--border);
          }

          .meta-stat {
            display: inline-flex;
            align-items: center;
            gap: 2px;

            .el-icon {
              font-size: 14px;
            }
          }
        }
      }
    }
  }

  .pagination-wrap {
    display: flex;
    justify-content: center;
    margin-top: var(--sp-6);
  }
}
</style>
