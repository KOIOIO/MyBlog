<template>
<div class="feedback-reply-form">
  <el-form
      :model="feedbackReplyFormData"
      :validate-on-rule-change="false"
  >
    <el-form-item :label="t('forms.feedbackReply.reply')" prop="reply">
      <el-input
          type="textarea"
          :rows="4"
          v-model="feedbackReplyFormData.reply"
          :placeholder="t('forms.feedbackReply.placeholder')"
      />
    </el-form-item>
    <el-form-item>
      <div class="button-group">
        <el-button
            type="primary"
            size="large"
            @click="submitForm"
        >{{ t('common.confirm') }}
        </el-button>
        <el-button
            size="large"
            @click="layoutStore.state.feedbackReplyVisible = false"
        >{{ t('common.cancel') }}
        </el-button>
      </div>
    </el-form-item>
  </el-form>
</div>
</template>

<script setup lang="ts">
import {defineProps, reactive} from 'vue';
import {ElMessage} from "element-plus";
import {useLayoutStore} from "@/stores/layout";
import {feedbackReply, type FeedbackReplyRequest} from "@/api/feedback";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const layoutStore = useLayoutStore()

const props = defineProps<{
  id: number;
}>();

const feedbackReplyFormData = reactive<FeedbackReplyRequest>({
  id: props.id,
  reply: '',
})

const submitForm = async () => {
  const res = await feedbackReply(feedbackReplyFormData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    layoutStore.state.shouldRefreshFeedbackTable = true
    layoutStore.state.feedbackReplyVisible = false
  }
}
</script>

<style scoped lang="scss">
.feedback-reply-form {
  .el-form {
    .el-form-item {
      .el-textarea__inner {
        background-color: var(--bg-elevated);
      }

      .button-group {
        margin-left: auto;
        display: flex;
        gap: var(--sp-2);

        .el-button {
          height: 40px;
          min-width: 88px;
          font-weight: 500;
        }
      }
    }
  }
}
</style>
