<template>
  <div class="forum-comments">
    <div class="page-head">
      <div class="page-head__text">
        <div class="page-head__title">{{ t('dashboard.forum.comments.title') }}</div>
        <div class="page-head__desc">{{ userStore.isAdmin ? t('dashboard.forum.comments.descAdmin') : t('dashboard.forum.comments.descUser') }}</div>
      </div>
    </div>

    <div class="filter-bar">
      <el-form :inline="true" :model="query">
        <el-form-item :label="t('dashboard.forum.comments.filter.content')">
          <el-input v-model="query.content" :placeholder="t('dashboard.forum.comments.filter.contentPh')" clearable/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="handleSearch">{{ t('common.query') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="tableData">
      <el-table-column prop="id" :label="t('dashboard.forum.comments.column.id')" width="70"/>
      <el-table-column prop="post_id" :label="t('dashboard.forum.comments.column.postId')" width="90"/>
      <el-table-column :label="t('dashboard.forum.comments.column.author')" width="140">
        <template #default="scope">
          {{ scope.row.user?.username || '-' }}
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.forum.comments.column.content')" min-width="280" show-overflow-tooltip>
        <template #default="scope">
          {{ scope.row.content }}
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.forum.comments.column.level')" width="90">
        <template #default="scope">
          <el-tag size="small" :type="scope.row.parent_id === 0 ? 'info' : 'warning'">
            {{ scope.row.parent_id === 0 ? t('dashboard.forum.comments.column.levelRoot') : t('dashboard.forum.comments.column.levelReply') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="t('dashboard.forum.comments.column.publishedAt')" width="170"/>
      <el-table-column :label="t('dashboard.forum.comments.column.actions')" width="100" fixed="right">
        <template #default="scope">
          <el-link type="danger" :underline="false" @click="handleDelete(scope.row)">{{ t('common.delete') }}</el-link>
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
import type {ForumComment} from "@/api/forum";
import {useUserStore} from "@/stores/user";

interface ForumManageCommentsData {
  list: ForumComment[];
  total: number;
}

const {t} = useI18n()
const userStore = useUserStore()
const tableData = ref<ForumComment[]>([])
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const query = reactive<{ content: string }>({
  content: '',
})

const getTableData = async () => {
  const res = await service({
    url: '/forum/manageComments',
    method: 'get',
    params: {
      page: page.value,
      page_size: page_size.value,
      content: query.content || undefined,
    },
  }) as ApiResponse<ForumManageCommentsData>

  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
  }
}

const handleSearch = () => {
  page.value = 1
  getTableData()
}

const doDelete = async (ids: number[]) => {
  const res = await service({
    url: '/forum/comment',
    method: 'delete',
    data: {ids},
  }) as ApiResponse<undefined>
  if (res.code === 0) {
    ElMessage.success(res.msg || t('dashboard.forum.comments.deleteSuccess'))
    getTableData()
  }
}

const handleDelete = (row: ForumComment) => {
  ElMessageBox.confirm(t('dashboard.forum.comments.confirmDelete'), t('dashboard.forum.comments.tip'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
  }).then(() => doDelete([row.id])).catch(() => {})
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
.forum-comments {
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

  .el-table {
    --el-table-border-color: var(--border);
  }

  .el-pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--sp-5);
  }
}
</style>
