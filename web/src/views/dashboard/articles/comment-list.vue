<template>
  <div class="comment-list">
    <div class="title">
      <el-row>{{ t('dashboard.comments.title') }}</el-row>
      <el-button-group>

        <el-button type="danger" icon="Delete" @click="commentBulkDeleteVisible = true;handleIdsToDelete()">
          {{ t('common.batchDelete') }}
        </el-button>

        <el-dialog
            v-model="commentBulkDeleteVisible"
            width="500"
            align-center
            destroy-on-close
        >
          <template #header>
            {{ t('dashboard.comments.deleteTitle') }}
          </template>
          {{ t('dashboard.common.deleteConfirm', { count: idsToDelete?.length ?? 0 }) }}
          <template #footer>
            <el-button type="primary" @click="handleBulkDelete(idsToDelete)">
              {{ t('common.confirm') }}
            </el-button>
            <el-button @click="commentBulkDeleteVisible = false">{{ t('common.cancel') }}</el-button>
          </template>
        </el-dialog>
      </el-button-group>
    </div>

    <div class="comment-list-request">
      <el-form :inline="true" :model="commentListRequest">
        <el-form-item :label="t('dashboard.comments.filter.articleId')">
          <el-input v-model="commentListRequest.article_id" :placeholder="t('dashboard.comments.filter.articleIdPh')" clearable/>
        </el-form-item>
        <el-form-item :label="t('dashboard.comments.filter.userUuid')">
          <el-input v-model="commentListRequest.user_uuid" :placeholder="t('dashboard.comments.filter.userUuidPh')" clearable/>
        </el-form-item>
        <el-form-item :label="t('dashboard.comments.filter.content')">
          <el-input v-model="commentListRequest.content" :placeholder="t('dashboard.comments.filter.contentPh')" clearable/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="getCommentTableData">{{ t('common.query') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
        ref="multipleCommentTableRef"
        :data="commentTableData"
    >
      <el-table-column type="selection" width="60"/>
      <el-table-column :label="t('dashboard.comments.column.articleId')" width="200">
        <template #default="scope:{ row: Comment, column: any, $index: number }">
          <el-link :href="'/article/'+scope.row.article_id">{{ scope.row.article_id }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.comments.column.user')" width="80">
        <template #default="scope:{ row: Comment, column: any, $index: number }">
          <user-card-popover :uuid="scope.row.user_uuid"/>
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.comments.column.content')">
        <template #default="scope:{ row: Comment, column: any, $index: number }">
          <MdPreview class="content" :modelValue="scope.row.content"/>
        </template>
      </el-table-column>
      <el-table-column :label="t('dashboard.comments.column.actions')" width="100">
        <template #default="scope:{ row: Comment, column: any, $index: number }">
          <el-button
              type="danger"
              @click="commentDeleteVisible=true;commentInfo=scope.row"
          >
            {{ t('common.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
        v-model="commentDeleteVisible"
        width="500"
        align-center
        destroy-on-close
    >
      <template #header>
        {{ t('dashboard.comments.deleteTitle') }}
      </template>
      {{ t('dashboard.common.deleteConfirm', { count: 1 }) }}
      <template #footer>
        <el-button type="primary" @click="handleDelete(commentInfo.id)">
          {{ t('common.confirm') }}
        </el-button>
        <el-button @click="commentDeleteVisible = false">{{ t('common.cancel') }}</el-button>
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
  type Comment,
  commentDelete, type CommentDeleteRequest,
  commentList,
  type CommentListRequest
} from "@/api/comment";
import {useLayoutStore} from "@/stores/layout";
import {ElMessage} from "element-plus";
import {useRoute, useRouter} from "vue-router";
import {useI18n} from "vue-i18n";
import UserCardPopover from "@/components/common/UserCardPopover.vue";
import {MdPreview} from "md-editor-v3";

const {t} = useI18n()

const multipleCommentTableRef = ref()
const commentTableData = ref<Comment[]>()
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const layoutStore = useLayoutStore()

const commentBulkDeleteVisible = ref(false)
let idsToDelete: number[]

const handleIdsToDelete = () => {
  idsToDelete = []

  const rows: Comment[] = multipleCommentTableRef.value.getSelectionRows()
  rows.forEach(row => {
    idsToDelete.push(row.id)
  })
}

const handleBulkDelete = async (ids: number[]) => {
  const requestData: CommentDeleteRequest = {
    ids: ids
  }
  const res = await commentDelete(requestData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  commentBulkDeleteVisible.value = false
  layoutStore.state.shouldRefreshCommentTable = true
}

const commentListRequest = reactive<CommentListRequest>({
  article_id: null,
  user_uuid: null,
  content: null,
  page: 1,
  page_size: 10,
})

const route = useRoute()
const router = useRouter()

onMounted(() => {
  commentListRequest.article_id = route.query.article_id as string || null
  commentListRequest.user_uuid = route.query.user_uuid as string || null
  commentListRequest.content = route.query.content as string || null
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})

const getCommentTableData = async () => {
  if (commentListRequest.article_id === "") {
    commentListRequest.article_id = null
  }
  if (commentListRequest.user_uuid === "") {
    commentListRequest.user_uuid = null
  }
  if (commentListRequest.content === "") {
    commentListRequest.content = null
  }

  commentListRequest.page = page.value
  commentListRequest.page_size = page_size.value

  const table = await commentList(commentListRequest)

  if (table.code === 0) {
    commentTableData.value = table.data.list
    total.value = table.data.total

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        article_id: commentListRequest.article_id,
        user_uuid: commentListRequest.user_uuid,
        content: commentListRequest.content,
        page: commentListRequest.page,
        page_size: commentListRequest.page_size,
      },
    })
  }
}

watch(() => route.query, (newQuery) => {
  commentListRequest.article_id = newQuery.article_id as string || null
  commentListRequest.user_uuid = newQuery.user_uuid as string || null
  commentListRequest.content = newQuery.content as string || null
  commentListRequest.page = Number(newQuery.page) || 1
  commentListRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getCommentTableData()
})

let commentInfo: Comment
const commentDeleteVisible = ref(false)

const handleDelete = async (id: number) => {
  let ids: number[] = []
  ids.push(id)

  const requestData: CommentDeleteRequest = {
    ids: ids
  }

  const res = await commentDelete(requestData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  commentDeleteVisible.value = false
  layoutStore.state.shouldRefreshCommentTable = true
}

watch(() => layoutStore.state.shouldRefreshCommentTable, (newVal) => {
  if (newVal) {
    getCommentTableData()
    layoutStore.state.shouldRefreshCommentTable = false
  }
})

const handleSizeChange = (val: number) => {
  page_size.value = val
  getCommentTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getCommentTableData()
}
</script>

<style scoped lang="scss">
.comment-list {
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
    }
  }

  .comment-list-request {
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

    /* 评论摘要：限高收起 */
    .content {
      max-height: 96px;
      overflow: hidden;
      color: var(--text-body);
      font-size: var(--fs-14);
    }
  }

  /* 操作列：文字删除按钮 */
  .el-table .cell .el-button--danger {
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

  .el-pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--sp-5);
  }
}
</style>
