<template>
  <el-card class="recent-comments">
    <el-row class="title">{{ t('components.recentComments.title') }}</el-row>
    <comment-item :comments="comments"/>
  </el-card>
</template>

<script setup lang="ts">
import CommentItem from "@/components/common/CommentItem.vue";
import {ref, watch} from "vue";
import {type Comment, commentNew} from "@/api/comment";
import {useLayoutStore} from "@/stores/layout";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const comments=ref<Comment[]>([])

const getRecentCommentInfo = async ()=>{
  const res = await commentNew()
  if (res.code===0){
    comments.value=res.data
  }
}

getRecentCommentInfo()

const layoutStore = useLayoutStore()

watch(() => layoutStore.state.shouldRefreshCommentList, (newVal) => {
  if (newVal) {
    getRecentCommentInfo()
  }
})
</script>

<style scoped lang="scss">
.recent-comments {
  margin-bottom: var(--sp-4);
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  animation: kf-fade-up var(--dur-mid) var(--ease-out) backwards;

  .title {
    font-size: var(--fs-16);
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: var(--sp-3);
  }
}

</style>
