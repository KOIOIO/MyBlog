<template>
  <div class="search-page">
    <web-navbar :noScroll="true"/>
    <div class="page">
      <div class="search-area">
        <div class="search-box">
          <el-input v-model="articleSearchRequest.query" :placeholder="t('pages.search.searchPlaceholder')" prefix-icon="Search"
                    maxlength="50" clearable size="large"
                    @change="changeArticleSearchItem"/>
        </div>

        <div class="filter-row">
          <span class="filter-label">{{ t('pages.search.categoryLabel') }}</span>
          <div class="chip-group">
            <button type="button" class="chip" :class="{ 'is-selected': articleSearchRequest.category === ''}"
                    @click="selectCategory('')">{{ t('common.all') }}</button>
            <button v-for="item in categoryArr" :key="item" type="button" class="chip"
                    :class="{ 'is-selected': articleSearchRequest.category === item}"
                    @click="selectCategory(item)">{{ categoryLabel(item) }}</button>
          </div>
        </div>

        <div class="filter-row">
          <span class="filter-label">{{ t('pages.search.tagLabel') }}</span>
          <div class="chip-group">
            <button type="button" class="chip" :class="{ 'is-selected': articleSearchRequest.tag === ''}"
                    @click="selectTag('')">{{ t('common.all') }}</button>
            <button v-for="item in tagArr" :key="item" type="button" class="chip"
                    :class="{ 'is-selected': articleSearchRequest.tag === item}"
                    @click="selectTag(item)">{{ tagLabel(item) }}</button>
          </div>
        </div>

        <div class="filter-row">
          <span class="filter-label">{{ t('pages.search.sortLabel') }}</span>
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

      <div class="results">
        <div v-for="(row, idx) in articleTableData" :key="row._id" class="result-item"
             :style="{ animationDelay: Math.min(idx, 7) * 30 + 'ms' }"
             @click="handleArticleJumps(row._id)">
          <h3 class="result-title">{{ row._source.title }}</h3>
          <p class="result-abstract">{{ row._source.abstract }}</p>
          <div class="result-meta">
            <span v-if="row._source.category" class="meta-category">{{ categoryLabel(row._source.category) }}</span>
            <span class="meta-date">{{ row._source.created_at }}</span>
            <span class="meta-stat">
              <el-icon><component is="View"/></el-icon> {{ row._source.views }}
            </span>
            <span class="meta-stat">
              <el-icon><component is="Star"/></el-icon> {{ row._source.likes }}
            </span>
          </div>
        </div>

        <!-- 空态：当前无 i18n key（pages.search.empty*），待 i18n 分片补文案后启用；
             此处仅保留结果行 stagger，空列表不渲染额外结构 -->
      </div>

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
  </div>
</template>

<script setup lang="ts">
import WebNavbar from "@/components/layout/WebNavbar.vue";
import type {Hit} from "@/api/common";
import {type Article, articleCategory, articleSearch, type ArticleSearchRequest, articleTags} from "@/api/article";
import {computed, nextTick, onMounted, reactive, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";
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

const route = useRoute()
const router = useRouter()

const categoryArr = ref<string[]>([])
const tagArr = ref<string[]>([])
const sortArr = computed(() => [
  {label: t('pages.search.sortDefault'), value: ""},
  {label: t('pages.search.sortTime'), value: "time"},
  {label: t('pages.search.sortComment'), value: "comment"},
  {label: t('pages.search.sortView'), value: "view"},
  {label: t('pages.search.sortLike'), value: "like"},
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

// chip 单选互斥：仅值真正变化时触发请求（对齐原 el-radio-group 语义；URL query 回写逻辑不变）
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

onMounted(() => {
  articleSearchRequest.query = route.query.query as string || ""
  articleSearchRequest.category = route.query.category as string || ""
  articleSearchRequest.tag = route.query.tag as string || ""
  articleSearchRequest.sort = route.query.sort as string || ""
  articleSearchRequest.order = route.query.order as string || "desc"
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})


const getArticleSearchTableData = async () => {
  articleSearchRequest.page = page.value;
  articleSearchRequest.page_size = page_size.value;

  const table = await articleSearch(articleSearchRequest)

  if (table.code === 0) {
    articleTableData.value = table.data.list;
    total.value = table.data.total;
  }

  await router.push({
    path: router.currentRoute.value.path,
    query: {
      query: articleSearchRequest.query,
      category: articleSearchRequest.category,
      tag: articleSearchRequest.tag,
      sort: articleSearchRequest.sort,
      order: articleSearchRequest.order,
      page: articleSearchRequest.page,
      page_size: articleSearchRequest.page_size,
    }
  })
}

watch(() => route.query, (newQuery) => {
  articleSearchRequest.query = newQuery.query as string || ""
  articleSearchRequest.category = newQuery.category as string || ""
  articleSearchRequest.tag = newQuery.tag as string || ""
  articleSearchRequest.sort = newQuery.sort as string || ""
  articleSearchRequest.order = newQuery.order as string || "desc"
  articleSearchRequest.page = Number(newQuery.page) || 1
  articleSearchRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getArticleSearchTableData()
})

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
.search-page {
  background-color: var(--bg);
  min-height: 100vh;

  .page {
    max-width: 640px;
    margin: 0 auto;
    padding: calc(70px + var(--sp-7)) var(--sp-4) var(--sp-9);
  }

  .search-area {
    margin-bottom: var(--sp-7);

    .search-box {
      margin-bottom: var(--sp-5);

      :deep(.el-input__wrapper) {
        height: 56px;
        border-radius: 12px;
        padding-left: var(--sp-4);
        transition: box-shadow 150ms cubic-bezier(.16,1,.3,1);
      }

    }

    .filter-row {
      display: flex;
      align-items: center;
      gap: var(--sp-3);
      margin-bottom: var(--sp-2);
      flex-wrap: wrap;

      .filter-label {
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

  .results {
    .result-item {
      padding: var(--sp-4) 0;
      border-bottom: 1px solid var(--border);
      cursor: pointer;
      transition: background-color 150ms ease-out, transform 80ms ease-out;
      padding-left: var(--sp-2);
      padding-right: var(--sp-2);
      border-radius: var(--radius-sm);
      animation: kf-fade-up var(--dur-base) var(--ease-out) backwards;

      &:hover {
        background-color: var(--bg-elevated);
      }

      &:active {
        transform: translateY(1px);
      }

      .result-title {
        font-size: var(--fs-18);
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 var(--sp-1);
        transition: color 150ms ease-out;
      }

      &:hover .result-title {
        color: var(--accent);
      }

      .result-abstract {
        font-size: var(--fs-14);
        color: var(--text-body);
        line-height: var(--lh-body);
        margin: 0 0 var(--sp-2);
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }

      .result-meta {
        display: flex;
        align-items: center;
        gap: var(--sp-3);
        font-size: var(--fs-12);
        color: var(--text-muted);

        .meta-category {
          color: var(--accent);
        }

        .meta-stat {
          display: inline-flex;
          align-items: center;
          gap: 2px;
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
