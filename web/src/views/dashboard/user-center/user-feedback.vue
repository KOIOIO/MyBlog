<template>
  <div class="user-feedback">
    <div class="page-header">
      <div class="page-title">{{ t('system.userCenter.feedback.title') }}</div>
      <div class="page-desc">{{ t('system.userCenter.feedback.desc') }}</div>
    </div>

    <el-table
        :data="userFeedbackTableData"
        :row-style="{height: '80px'}"
    >
      <el-table-column :label="t('common.time')" width="150">
        <template #default="scope:{ row: Feedback, column: any, $index: number }">
          {{ getTime(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="content" :label="t('common.content')"/>
      <el-table-column prop="reply" :label="t('system.userCenter.feedback.reply')"/>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {type Feedback, feedbackInfo} from "@/api/feedback";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const getTime = (date: Date): string => {
  const time = new Date(date)
  return time.toLocaleString()
}

const userFeedbackTableData = ref<Feedback[]>()
const getUserFeedbackTableData = async () => {
  const table = await feedbackInfo()

  if (table.code === 0) {
    userFeedbackTableData.value = table.data
  }
}
getUserFeedbackTableData()
</script>

<style scoped lang="scss">
.user-feedback {
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
  }
}
</style>
