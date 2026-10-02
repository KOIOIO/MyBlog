<template>
  <div ref="scrollRef" class="message-list">
    <div v-if="messages.length === 0" class="empty">
      <div class="empty-icon">🤖</div>
      <p class="empty-title">{{ t('pages.agent.emptyChatTitle') }}</p>
      <p class="empty-desc">{{ t('pages.agent.emptyChatDesc') }}</p>
    </div>
    <div v-for="msg in messages" :key="msg.id" class="msg-row" :class="msg.role">
      <div class="avatar">{{ msg.role === 'assistant' ? '🤖' : '👤' }}</div>
      <div class="bubble-wrap">
        <div v-if="msg.role === 'user' && msg.article_ids && msg.article_ids.length > 0" class="ref-tag">
          📄 {{ t('pages.agent.referencedArticles', {n: msg.article_ids.length}) }}
        </div>
        <div class="bubble">
          <span v-if="msg.content" class="content">{{ msg.content }}</span>
          <span v-else-if="streaming" class="typing">▍</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {nextTick, ref, watch} from "vue";
import {useI18n} from "vue-i18n";
import type {AgentMessage} from "@/api/agent";

const props = defineProps<{
    messages: AgentMessage[];
    streaming: boolean;
}>();

const {t} = useI18n();
const scrollRef = ref<HTMLElement | null>(null);

const scrollToBottom = (): void => {
    nextTick(() => {
        if (scrollRef.value) {
            scrollRef.value.scrollTop = scrollRef.value.scrollHeight;
        }
    });
};

watch(() => [props.messages, props.streaming], scrollToBottom, {deep: true});
</script>

<style scoped lang="scss">
.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;

  .empty {
    margin: auto;
    text-align: center;
    color: var(--text-secondary);

    .empty-icon { font-size: 48px; }
    .empty-title { margin: 12px 0 4px; font-size: 16px; color: var(--text-body); }
    .empty-desc { font-size: 13px; }
  }

  .msg-row {
    display: flex;
    gap: 10px;
    max-width: 82%;

    &.user {
      align-self: flex-end;
      flex-direction: row-reverse;

      .bubble {
        background: var(--el-color-primary);
        color: #fff;
      }
    }

    &.assistant {
      align-self: flex-start;
    }
  }

  .avatar {
    width: 34px;
    height: 34px;
    border-radius: 50%;
    background: var(--el-fill-color-light);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    flex-shrink: 0;
  }

  .bubble-wrap {
    min-width: 0;
  }

  .ref-tag {
    font-size: 12px;
    color: var(--el-color-primary);
    margin-bottom: 4px;
  }

  .bubble {
    padding: 10px 14px;
    border-radius: 12px;
    background: var(--el-fill-color-light);
    color: var(--text-body);
    line-height: 1.6;
    font-size: 14px;
    word-break: break-word;
    white-space: pre-wrap;

    .typing { animation: blink 1s step-end infinite; }
  }
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
</style>
