<template>
  <div class="user-star">
    <div class="page-header">
      <div class="page-title">{{ t('system.userCenter.star.title') }}</div>
      <div class="page-desc">{{ t('system.userCenter.star.desc') }}</div>
    </div>

    <el-table
        :data="articleLikesListData"
    >
      <el-table-column :label="t('common.cover')" width="100">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          <el-image :src="scope.row._source.cover" alt=""/>
        </template>
      </el-table-column>
      <el-table-column prop="_source.title" :label="t('common.title')" width="120"/>
      <el-table-column prop="_source.category" :label="t('system.userCenter.star.category')" width="80"/>
      <el-table-column :label="t('common.tags')" width="120">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          <el-tag v-for="tag in scope.row._source.tags">{{ tag }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('system.userCenter.star.abstract')">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          <el-text line-clamp="5">{{ scope.row._source.abstract }}</el-text>
        </template>
      </el-table-column>
      <el-table-column prop="_source.created_at" :label="t('system.userCenter.star.publishTime')" width="102"/>
      <el-table-column :label="t('system.userCenter.star.articleId')" width="200">
        <template #default="scope:{ row: Hit<Article>, column: any, $index: number }">
          <el-link :href="'/article/'+scope.row._id">{{ scope.row._id }}</el-link>
        </template>
      </el-table-column>
    </el-table>

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
import {type Article} from "@/api/article";
import {useRoute, useRouter} from "vue-router";
import type {Hit, PageInfo} from "@/api/common";
import {articleLikesList} from "@/api/article";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const articleLikesListData = ref<Hit<Article>[]>()
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const articleLikesListRequest = reactive<PageInfo>({
  page: 1,
  page_size: 10,
})

const route = useRoute()
const router = useRouter()

onMounted(() => {
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})

const getArticleLikesListData = async () => {
  articleLikesListRequest.page = page.value
  articleLikesListRequest.page_size = page_size.value

  const table = await articleLikesList(articleLikesListRequest)

  if (table.code === 0) {
    articleLikesListData.value = table.data.list
    total.value = table.data.total

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        page: articleLikesListRequest.page,
        page_size: articleLikesListRequest.page_size,
      },
    })
  }
}

watch(() => route.query, (newQuery) => {
  articleLikesListRequest.page = Number(newQuery.page) || 1
  articleLikesListRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getArticleLikesListData()
})

const handleSizeChange = (val: number) => {
  page_size.value = val
  getArticleLikesListData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getArticleLikesListData()
}
</script>

<style scoped lang="scss">
.user-star {
  .page-header {
    margin-bottom: var(--sp-5);

    .page-title {
      font-size: var(--fs-24);
      font-weight: 600;
      color: var(--text-primary);
      line-height: var(--lh-title);
    }

    .page-desc {
      margin-top: var(--sp-1);
      font-size: var(--fs-14);
      color: var(--text-muted);
    }
  }

  .el-table {
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);

    .el-image {
      height: 48px;
      border-radius: var(--radius-sm);
    }
  }

  .el-pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--sp-5);
  }
}
</style>
