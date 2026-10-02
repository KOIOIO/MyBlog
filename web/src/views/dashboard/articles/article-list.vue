<template>
  <div class="article-list">
    <div class="title">
      <el-row>{{ t('dashboard.articles.list.title') }}</el-row>
      <el-button-group>
        <el-button type="success" icon="Plus" @click="handleToPublishArticle">
          {{ t('dashboard.articles.list.new') }}
        </el-button>

        <el-button type="danger" icon="Delete" @click="articleBulkDeleteVisible = true;handleIdsToDelete()">
          {{ t('common.batchDelete') }}
        </el-button>

        <el-dialog
            v-model="articleBulkDeleteVisible"
            width="500"
            align-center
            destroy-on-close
        >
          <template #header>
            {{ t('dashboard.articles.list.deleteTitle') }}
          </template>
          {{ t('dashboard.common.deleteConfirm', { count: idsToDelete?.length ?? 0 }) }}
          <template #footer>
            <el-button type="primary" @click="handleBulkDelete(idsToDelete)">
              {{ t('common.confirm') }}
            </el-button>
            <el-button @click="articleBulkDeleteVisible = false">{{ t('common.cancel') }}</el-button>
          </template>
        </el-dialog>
      </el-button-group>
    </div>

    <div class="article-list-request">
      <el-form :inline="true" :model="articleListRequest">
        <el-form-item :label="t('dashboard.articles.list.filter.title')">
          <el-input v-model="articleListRequest.title" :placeholder="t('dashboard.articles.list.filter.titlePh')" clearable/>
        </el-form-item>
        <el-form-item :label="t('dashboard.articles.list.filter.category')">
          <el-input v-model="articleListRequest.category" :placeholder="t('dashboard.articles.list.filter.categoryPh')" clearable/>
        </el-form-item>
        <el-form-item :label="t('dashboard.articles.list.filter.abstract')">
          <el-input v-model="articleListRequest.abstract" :placeholder="t('dashboard.articles.list.filter.abstractPh')" clearable/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="getArticleTableData">{{ t('common.query') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
        ref="multipleArticleTableRef"
        :data="articleTableData"
    >
      <el-table-column type="selection" width="60"/>
      <el-table-column :label="t('dashboard.articles.list.column.cover')" width="100">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          <el-image :src="scope.row._source.cover" alt=""/>
        </template>
      </el-table-column>
      <el-table-column prop="_source.title" :label="t('dashboard.articles.list.column.title')" width="120"/>
      <el-table-column :label="t('dashboard.articles.list.column.category')" width="80">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          {{ categoryLabel(scope.row._source.category) }}
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.articles.list.column.tag')" width="120">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          <el-tag v-for="tag in scope.row._source.tags" :key="tag">{{ tagLabel(tag) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.articles.list.column.abstract')">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          <el-text line-clamp="5">{{ scope.row._source.abstract }}</el-text>
        </template>
      </el-table-column>
      <el-table-column prop="_source.created_at" :label="t('dashboard.articles.list.column.publishedAt')" width="102"/>
      <el-table-column :label="t('dashboard.articles.list.column.articleId')" width="220">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          <el-link :href="'/article/'+scope.row._id">{{ scope.row._id }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.articles.list.column.actions')" width="220">
        <template #default="scope:{ row: any, column: any, $index: number }">
          <el-button
              v-if="scope.row._source.is_top === 1"
              type="primary" link
              @click="handleSetTop(scope.row, false)"
          >
            {{ t('dashboard.articles.list.cancelTop') }}
          </el-button>
          <el-button
              v-else
              type="primary" link
              @click="handleSetTop(scope.row, true)"
          >
            {{ t('dashboard.articles.list.top') }}
          </el-button>
          <el-button
              type="warning"
              @click="layoutStore.state.articleUpdateVisible=true;articleInfo=scope.row"
          >
            {{ t('dashboard.articles.list.update') }}
          </el-button>
          <el-button
              type="danger"
              @click="articleDeleteVisible=true;articleInfo=scope.row"
          >
            {{ t('common.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
        v-model="articleUpdateVisible"
        width="500"
        align-center
        destroy-on-close
        :before-close="articleUpdateVisibleSynchronization"
    >
      <template #header>
        {{ t('dashboard.articles.list.updateTitle') }}
      </template>
      <article-update-form :article=articleInfo />
      <template #footer>
      </template>
    </el-dialog>

    <el-dialog
        v-model="articleDeleteVisible"
        width="500"
        align-center
        destroy-on-close
    >
      <template #header>
        {{ t('dashboard.articles.list.deleteTitle') }}
      </template>
      {{ t('dashboard.common.deleteConfirm', { count: 1 }) }}
      <template #footer>
        <el-button type="primary" @click="handleDelete(articleInfo._id)">
          {{ t('common.confirm') }}
        </el-button>
        <el-button @click="articleDeleteVisible = false">{{ t('common.cancel') }}</el-button>
      </template>
    </el-dialog>

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
</template>

<script setup lang="ts">
import {nextTick, onMounted, reactive, ref, watch} from "vue";
import {
  type Article,
  articleDelete, type ArticleDeleteRequest,
  articleList,
  type ArticleListRequest,
  articleSetTop, type ArticleSetTopRequest
} from "@/api/article";
import {useLayoutStore} from "@/stores/layout";
import {ElMessage} from "element-plus";
import {useRoute, useRouter} from "vue-router";
import {useI18n} from "vue-i18n";
import ArticleUpdateForm from "@/components/forms/ArticleUpdateForm.vue";
import type {Hit} from "@/api/common";
import {type Tag, useTagStore} from "@/stores/tag";
import {categoryLabel, tagLabel} from "@/i18n/meta";

const {t} = useI18n()

const multipleArticleTableRef = ref()
const articleTableData = ref<Hit<Article>[]>()
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const tagStore = useTagStore()

const handleToPublishArticle = () => {
  const newTag: Tag = {
    title: t("menu.articles.publish"),
    name: "article-publish"
  }
  const exists = tagStore.state.tags.some(tag => tag.name === newTag.name);
  if (exists) {
    return;
  }
  tagStore.state.tags.push(newTag);
  router.push({name:'article-publish'})
}

const layoutStore = useLayoutStore()

const articleBulkDeleteVisible = ref(false)
let idsToDelete: string[]

const handleIdsToDelete = () => {
  idsToDelete = []

  const rows: Hit<Article>[] = multipleArticleTableRef.value.getSelectionRows()
  rows.forEach(row => {
    idsToDelete.push(row._id)
  })
}

const handleBulkDelete = async (ids: string[]) => {
  const requestData: ArticleDeleteRequest = {
    ids: ids
  }
  const res = await articleDelete(requestData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  articleBulkDeleteVisible.value = false
  layoutStore.state.shouldRefreshArticleTable = true
}

const articleListRequest = reactive<ArticleListRequest>({
  title: null,
  category: null,
  abstract: null,
  page: 1,
  page_size: 10,
})

const route = useRoute()
const router = useRouter()

onMounted(() => {
  articleListRequest.title = route.query.title as string || null
  articleListRequest.category = route.query.category as string || null
  articleListRequest.abstract = route.query.abstract as string || null
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})

const getArticleTableData = async () => {
  if (articleListRequest.title === "") {
    articleListRequest.title = null
  }
  if (articleListRequest.category === "") {
    articleListRequest.category = null
  }
  if (articleListRequest.abstract === "") {
    articleListRequest.abstract = null
  }

  articleListRequest.page = page.value
  articleListRequest.page_size = page_size.value

  const table = await articleList(articleListRequest)

  if (table.code === 0) {
    articleTableData.value = table.data.list
    total.value = table.data.total

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        title: articleListRequest.title,
        category: articleListRequest.category,
        abstract: articleListRequest.abstract,
        page: articleListRequest.page,
        page_size: articleListRequest.page_size,
      },
    })
  }
}

watch(() => route.query, (newQuery) => {
  articleListRequest.title = newQuery.title as string || null
  articleListRequest.category = newQuery.category as string || null
  articleListRequest.abstract = newQuery.abstract as string || null
  articleListRequest.page = Number(newQuery.page) || 1
  articleListRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getArticleTableData()
})

let articleInfo: Hit<Article>
const articleDeleteVisible = ref(false)

const articleUpdateVisible = ref(layoutStore.state.articleUpdateVisible)
watch(
    () => layoutStore.state.articleUpdateVisible,
    (newValue) => {
      articleUpdateVisible.value = newValue
    }
)

const articleUpdateVisibleSynchronization = () => {
  layoutStore.state.articleUpdateVisible = false
}

const handleSetTop = async (row: Hit<Article>, isTop: boolean) => {
  const req: ArticleSetTopRequest = {
    id: row._id,
    is_top: isTop,
  }
  const res = await articleSetTop(req)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    await getArticleTableData()
  }
}

const handleDelete = async (id: string) => {
  let ids: string[] = []
  ids.push(id)

  const requestData: ArticleDeleteRequest = {
    ids: ids
  }

  const res = await articleDelete(requestData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  articleDeleteVisible.value = false
  layoutStore.state.shouldRefreshArticleTable = true
}

watch(() => layoutStore.state.shouldRefreshArticleTable, (newVal) => {
  if (newVal) {
    getArticleTableData()
    layoutStore.state.shouldRefreshArticleTable = false
  }
})

const handleSizeChange = (val: number) => {
  page_size.value = val
  getArticleTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getArticleTableData()
}
</script>

<style scoped lang="scss">
.article-list {
  .title {
    display: flex;
    align-items: center;
    margin-bottom: var(--sp-5);

    .el-row {
      font-size: var(--fs-24);
      font-weight: 600;
      color: var(--text-primary);
    }

    .el-button-group {
      margin-left: auto;

      .el-button {
        margin-left: var(--sp-3);
      }

      /* 新建文章：accent 实心白字（文字色沿用 EP 默认白） */
      .el-button--success {
        --el-button-bg-color: var(--accent);
        --el-button-border-color: var(--accent);
        --el-button-hover-bg-color: var(--accent);
        --el-button-hover-border-color: var(--accent);
        --el-button-active-bg-color: var(--accent);
        --el-button-active-border-color: var(--accent);
      }
    }
  }

  .article-list-request {
    background-color: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    padding: var(--sp-4) var(--sp-5);
    margin-bottom: var(--sp-5);

    .el-form {
      display: flex;
      flex-wrap: wrap;
      align-items: center;

      .el-form-item {
        margin-right: var(--sp-4);
        margin-bottom: 0;
      }

      /* 输入框 40px 高，细线圆角 */
      .el-input__wrapper {
        height: 40px;
        border-radius: var(--radius-sm);
      }
    }
  }

  .el-table {
    --el-table-border-color: var(--border);

    .el-image {
      height: 48px;
      width: 48px;
      border-radius: var(--radius-sm);
    }
  }

  /* 操作列：文字按钮（更新 accent / 删除 danger） */
  .el-table .cell {
    .el-button + .el-button {
      margin-left: var(--sp-4);
    }

    .el-button--warning {
      --el-button-bg-color: transparent;
      --el-button-border-color: transparent;
      --el-button-text-color: var(--accent);
      --el-button-hover-bg-color: var(--accent-weak);
      --el-button-hover-text-color: var(--accent);
      --el-button-hover-border-color: transparent;
      --el-button-active-bg-color: var(--accent-weak);
      --el-button-active-text-color: var(--accent);
      padding: 4px 0;
    }

    .el-button--danger {
      --el-button-bg-color: transparent;
      --el-button-border-color: transparent;
      --el-button-text-color: var(--el-color-danger);
      --el-button-hover-bg-color: var(--el-color-danger-light-9);
      --el-button-hover-text-color: var(--el-color-danger);
      --el-button-hover-border-color: transparent;
      --el-button-active-bg-color: var(--el-color-danger-light-9);
      --el-button-active-text-color: var(--el-color-danger);
      padding: 4px 0;
    }
  }

  .el-pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--sp-5);
  }
}
</style>
