<template>
  <el-card class="feedback">
    <el-row class="title">{{ t('components.feedback.title') }}</el-row>
    <el-input type="textarea" :rows="4" v-model="feedbackCreateFormData.content" maxlength="100" show-word-limit
              :placeholder="t('components.feedback.placeholder')"></el-input>
    <div class="content">
      <el-text class="login-tip">{{ t('components.feedback.loginTip') }}</el-text>
      <div class="button-group">
        <el-button @click="submitForm" type="primary">{{ t('common.confirm') }}</el-button>
        <el-button @click="feedbackCreateFormData.content=''">{{ t('common.cancel') }}</el-button>
      </div>
    </div>
    <el-row class="title-sub">{{ t('components.feedback.listTitle') }}</el-row>
    <div class="footer">
      <div class="feedback-new" v-for="item in feedbackInfoList" :key="item.time">
        <el-row>{{ item.content }}</el-row>
        <el-row class="container">
          <div class="time">{{ item.time }}</div>
        </el-row>
        <div class="reply">
          <el-text v-if="item.reply!==''">{{ t('components.feedback.replyLabel') }}{{ item.reply }}</el-text>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import {reactive, ref, watch} from "vue";
import {feedbackCreate, type FeedbackCreateRequest, feedbackNew} from "@/api/feedback";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const feedbackCreateFormData = reactive<FeedbackCreateRequest>({
  content: '',
})

interface FeedbackNew {
  content: string;
  reply: string;
  time: string;
}

const feedbackInfoList = ref<FeedbackNew[]>([])

const shouldRefreshFeedbackInfoTable = ref(false)
watch(() => shouldRefreshFeedbackInfoTable.value, (newVal) => {
  if (newVal) {
    getFeedbackNew()
    shouldRefreshFeedbackInfoTable.value = false;
  }
})

const submitForm = async () => {
  const res = await feedbackCreate(feedbackCreateFormData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    feedbackCreateFormData.content = ''
    shouldRefreshFeedbackInfoTable.value = true
  }
}

const getFeedbackNew = async () => {
  feedbackInfoList.value = []
  const res = await feedbackNew()
  if (res.code === 0) {
    res.data.forEach(value => {
      const date = new Date(value.created_at);
      const info: FeedbackNew = {
        content: value.content,
        reply: value.reply,
        time: date.toLocaleString(),
      }
      feedbackInfoList.value.push(info)
    })
  }
}

getFeedbackNew()
</script>

<style scoped lang="scss">
.feedback {
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  animation: kf-fade-up var(--dur-mid) var(--ease-out) backwards;

  .title {
    font-size: var(--fs-16);
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: var(--sp-3);
  }

  .content {
    display: flex;
    align-items: center;
    margin: var(--sp-3) 0;

    .login-tip {
      font-size: var(--fs-12);
      color: var(--text-muted);
    }

    .button-group {
      margin-left: auto;
    }
  }

  .title-sub {
    font-size: var(--fs-14);
    font-weight: 600;
    color: var(--text-primary);
    margin: var(--sp-4) 0 var(--sp-3);
  }

  .footer {
    .feedback-new {
      border: 1px solid var(--border);
      border-radius: var(--radius-sm);
      padding: var(--sp-3);
      margin-bottom: var(--sp-3);
      font-size: var(--fs-14);
      color: var(--text-body);
      transition: background-color 150ms ease-out, transform 80ms ease-out;

      &:hover {
        background-color: var(--bg-elevated);
      }

      &:active {
        transform: translateY(1px);
      }

      .container {
        display: flex;
        border-bottom: 1px solid var(--border);
        padding-bottom: var(--sp-2);
        margin-bottom: var(--sp-2);

        .time {
          font-size: var(--fs-12);
          color: var(--text-muted);
          margin-left: auto;
        }
      }

      .reply {
        font-size: var(--fs-12);
        color: var(--text-muted);
      }
    }
  }
}

</style>
