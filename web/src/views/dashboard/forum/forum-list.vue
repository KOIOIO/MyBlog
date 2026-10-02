<template>
  <div class="forum-list">
    <div class="page-head">
      <div class="page-head__text">
        <div class="page-head__title">{{ t('dashboard.forum.list.title') }}</div>
        <div class="page-head__desc">{{ userStore.isAdmin ? t('dashboard.forum.list.descAdmin') : t('dashboard.forum.list.descUser') }}</div>
      </div>
      <div class="page-head__actions">
        <el-button type="danger" icon="Delete" @click="handleBulkDelete">{{ t('common.batchDelete') }}</el-button>
      </div>
    </div>

    <div class="filter-bar">
      <el-form :inline="true" :model="query">
        <el-form-item :label="t('dashboard.forum.list.filter.title')">
          <el-input v-model="query.title" :placeholder="t('dashboard.forum.list.filter.titlePh')" clearable/>
        </el-form-item>
        <el-form-item :label="t('dashboard.forum.list.filter.category')">
          <el-select v-model="query.category" :placeholder="t('dashboard.forum.list.filter.categoryPh')" clearable style="width: 160px">
            <el-option :label="t('common.all')" value=""/>
            <el-option :label="categoryLabel('技术')" value="技术"/>
            <el-option :label="categoryLabel('生活')" value="生活"/>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleSearch">{{ t('common.query') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
        ref="tableRef"
        :data="tableData"
        @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="48"/>
      <el-table-column prop="id" :label="t('dashboard.forum.list.column.id')" width="70"/>
      <el-table-column :label="t('dashboard.forum.list.column.title')" min-width="200" show-overflow-tooltip>
        <template #default="scope">
          {{ scope.row.title }}
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.forum.list.column.author')" width="120">
        <template #default="scope">
          {{ scope.row.user?.username || '-' }}
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.forum.list.column.category')" width="80">
        <template #default="scope">
          {{ categoryLabel(scope.row.category) }}
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.forum.list.column.tag')" width="160">
        <template #default="scope">
          <el-tag v-for="tag in getTags(scope.row)" :key="tag" size="small" class="tag-item">{{ tagLabel(tag) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="like_count" :label="t('dashboard.forum.list.column.likes')" width="70"/>
      <el-table-column prop="comment_count" :label="t('dashboard.forum.list.column.comments')" width="70"/>
      <el-table-column prop="view_count" :label="t('dashboard.forum.list.column.views')" width="70"/>
      <el-table-column prop="created_at" :label="t('dashboard.forum.list.column.publishedAt')" width="170"/>
      <el-table-column :label="t('dashboard.forum.list.column.actions')" width="140" fixed="right">
        <template #default="scope">
          <el-link type="primary" :underline="false" :href="'/forum/' + scope.row.id" target="_blank">{{ t('dashboard.forum.list.view') }}</el-link>
          <el-link type="danger" :underline="false" class="op-link" @click="handleDelete(scope.row)">{{ t('common.delete') }}</el-link>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
        :current-page="page"
        :page-size="page_size"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
    />
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive, ref} from "vue";
import {useI18n} from "vue-i18n";
import {ElMessage, ElMessageBox} from "element-plus";
import service from "@/utils/request";
import type {ApiResponse} from "@/utils/request";
import type {ForumPost} from "@/api/forum";
import {useUserStore} from "@/stores/user";
import {categoryLabel, tagLabel} from "@/i18n/meta";

interface ForumManageListData {
  list: ForumPost[];
  total: number;
}

const {t} = useI18n()
const userStore = useUserStore()
const tableRef = ref()
const tableData = ref<ForumPost[]>([])
const selectedRows = ref<ForumPost[]>([])
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const query = reactive<{ title: string; category: string }>({
  title: '',
  category: '',
})

const getTableData = async () => {
  const res = await service({
    url: '/forum/manageList',
    method: 'get',
    params: {
      page: page.value,
      page_size: page_size.value,
      title: query.title || undefined,
      category: query.category || undefined,
    },
  }) as ApiResponse<ForumManageListData>

  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

// 兼容后端返回 tags 为 JSON 字符串或数组两种格式
const getTags = (row: ForumPost): string[] => {
  if (!row.tags) return []
  if (Array.isArray(row.tags)) return row.tags
  try {
    const parsed = JSON.parse(row.tags as unknown as string)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

const handleSearch = () => {
  page.value = 1
  getTableData()
}

const handleSelectionChange = (rows: ForumPost[]) => {
  selectedRows.value = rows
}

const doDelete = async (ids: number[]) => {
  const res = await service({
    url: '/forum/delete',
    method: 'delete',
    data: {ids},
  }) as ApiResponse<undefined>
  if (res.code === 0) {
    ElMessage.success(res.msg || t('dashboard.forum.list.deleteSuccess'))
    getTableData()
  }
}

const handleDelete = (row: ForumPost) => {
  ElMessageBox.confirm(t('dashboard.forum.list.confirmDeletePost', { title: row.title }), t('dashboard.forum.list.tip'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
  }).then(() => doDelete([row.id])).catch(() => {})
}

const handleBulkDelete = () => {
  if (!selectedRows.value.length) {
    ElMessage.warning(t('dashboard.forum.list.selectFirst'))
    return
  }
  const ids = selectedRows.value.map(r => r.id)
  ElMessageBox.confirm(t('dashboard.forum.list.confirmBulkDelete', { count: ids.length }), t('dashboard.forum.list.tip'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
  }).then(() => doDelete(ids)).catch(() => {})
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val: number) => {
  page_size.value = val
  page.value = 1
  getTableData()
}

onMounted(getTableData)
</script>

<style scoped lang="scss">
.forum-list {
  .page-head {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    margin-bottom: var(--sp-5);

    &__title {
      font-size: var(--fs-24);
      font-weight: 600;
      color: var(--text-primary);
      line-height: var(--lh-title);
    }

    &__desc {
      margin-top: var(--sp-1);
      font-size: var(--fs-14);
      color: var(--text-muted);
    }

    &__actions {
      margin-left: auto;
    }
  }

  .filter-bar {
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

      .el-input__wrapper {
        height: 40px;
        border-radius: var(--radius-sm);
      }
    }
  }

  .tag-item {
    margin-right: var(--sp-1);
  }

  .el-table {
    --el-table-border-color: var(--border);
  }

  .op-link {
    margin-left: var(--sp-3);
  }

  .el-pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--sp-5);
  }
}
</style>
