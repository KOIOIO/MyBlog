<template>
  <div class="user-comment">
    <div class="page-header">
      <div class="page-title">{{ t('system.userCenter.comment.title') }}</div>
      <div class="page-desc">{{ t('system.userCenter.comment.desc') }}</div>
    </div>
    <div class="table-data" v-for="item in userCommentTableData">
      <div class="link">
        <el-link :href="'/article/'+item.article_id">{{ item.article_id }}</el-link>
      </div>
      <div class="item">
        <comment-item :comments="[item]"/>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from "vue";
import CommentItem from "@/components/common/CommentItem.vue";
import {type Comment, commentInfo} from "@/api/comment";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const userCommentTableData = ref<Comment[]>()
const getUserCommentTableData = async () => {
  const table = await commentInfo()

  if (table.code === 0) {
    userCommentTableData.value = table.data
  }
}
getUserCommentTableData()

const layoutStore = useLayoutStore()

watch(() => layoutStore.state.shouldRefreshCommentList, (newVal) => {
  if (newVal) {
    getUserCommentTableData()
  }
})
</script>

<style scoped lang="scss">
.user-comment {
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

  .table-data {
    display: flex;
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    margin-bottom: var(--sp-3);
    transition: background-color 150ms ease-out;

    .link {
      width: 160px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-right: 1px solid var(--border);
      padding: var(--sp-3);
    }

    .item {
      flex: 1;
      padding: var(--sp-3);
    }
  }
}
</style>
