<template>
  <div ref="scrollRef" class="message-list">
    <div v-if="messages.length === 0" class="empty">
      <img class="empty-icon" src="/images/agent-avatar.jpg" alt="agent"/>
      <p class="empty-title">{{ t('pages.agent.emptyChatTitle') }}</p>
      <p class="empty-desc">{{ t('pages.agent.emptyChatDesc') }}</p>
    </div>
    <div v-for="msg in messages" :key="msg.id" class="msg-row" :class="msg.role">
      <div class="avatar">
        <img v-if="msg.role === 'assistant'" class="avatar-img" src="/images/agent-avatar.jpg" alt="agent"/>
        <el-avatar v-else-if="userStore.state.userInfo.avatar" class="avatar-img" :size="34" :src="userStore.state.userInfo.avatar"/>
        <span v-else class="avatar-fallback">{{ userInitial }}</span>
      </div>
      <div class="bubble-wrap">
        <div v-if="msg.role === 'user' && msg.article_ids && msg.article_ids.length > 0" class="ref-tag">
          📄 {{ t('pages.agent.referencedArticles', {n: msg.article_ids.length}) }}
        </div>
        <div class="bubble">
          <div v-if="msg.role === 'assistant' && msg.content" class="md-body" v-html="renderMarkdown(msg.content)"></div>
          <span v-else-if="msg.content" class="content">{{ msg.content }}</span>
          <span v-else-if="streaming" class="typing">▍</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, nextTick, onMounted, onUnmounted, ref, watch} from "vue";
import {useI18n} from "vue-i18n";
import {useUserStore} from "@/stores/user";
import mermaid from "mermaid";
import {renderMarkdown} from "@/utils/markdown";
import type {AgentMessage} from "@/api/agent";

const props = defineProps<{
    messages: AgentMessage[];
    streaming: boolean;
}>();

const {t} = useI18n();
const userStore = useUserStore();
const scrollRef = ref<HTMLElement | null>(null);

const userInitial = computed(() => {
    const name = userStore.state.userInfo.username || 'U';
    return name.charAt(0).toUpperCase();
});

const scrollToBottom = (): void => {
    nextTick(() => {
        if (scrollRef.value) {
            scrollRef.value.scrollTop = scrollRef.value.scrollHeight;
        }
    });
};

/**
 * 将消息中未处理的 ```mermaid 代码块渲染为图表。
 * 流式输出期间内容不断变化，由 MutationObserver 防抖触发；
 * mermaid 对已处理节点（data-processed）幂等跳过，失败仅告警保留源码。
 */
const renderMermaid = async (): Promise<void> => {
    const nodes = scrollRef.value?.querySelectorAll('.mermaid:not([data-processed="true"])');
    if (!nodes || nodes.length === 0) {
        return;
    }
    const isDark = document.documentElement.classList.contains('dark');
    mermaid.initialize({
        startOnLoad: false,
        securityLevel: 'strict',
        theme: isDark ? 'dark' : 'default',
    });
    try {
        await mermaid.run({nodes: Array.from(nodes) as HTMLElement[]});
    } catch (e) {
        // 单个图解析失败不影响其他消息；保留源码文本便于排查
        console.warn('mermaid render failed:', e);
    }
};

let mmObserver: MutationObserver | null = null;
let mmTimer: number | undefined;

const scheduleMermaid = (): void => {
    if (mmTimer) {
        window.clearTimeout(mmTimer);
    }
    mmTimer = window.setTimeout(() => {
        renderMermaid();
    }, 400);
};

onMounted(() => {
    if (scrollRef.value) {
        mmObserver = new MutationObserver(scheduleMermaid);
        mmObserver.observe(scrollRef.value, {childList: true, subtree: true, characterData: true});
    }
});

onUnmounted(() => {
    if (mmObserver) {
        mmObserver.disconnect();
        mmObserver = null;
    }
    if (mmTimer) {
        window.clearTimeout(mmTimer);
    }
});

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

    .empty-icon {
      width: 72px;
      height: 72px;
      border-radius: 50%;
      object-fit: cover;
      margin: 0 auto;
      display: block;
      box-shadow: var(--el-box-shadow-light);
    }
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
    overflow: hidden;

    .avatar-img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      border-radius: 50%;
      display: block;
    }

    .avatar-fallback {
      color: var(--el-color-primary);
      font-weight: 600;
    }
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

    .md-body {
      white-space: normal;

      > :first-child { margin-top: 0; }
      > :last-child { margin-bottom: 0; }

      p { margin: 0.5em 0; }
      h1, h2, h3, h4 { margin: 0.8em 0 0.4em; font-weight: 600; line-height: 1.3; }
      h1 { font-size: 1.4em; }
      h2 { font-size: 1.25em; }
      h3 { font-size: 1.15em; }
      h4 { font-size: 1.05em; }

      ul, ol { margin: 0.5em 0; padding-left: 1.5em; }
      li { margin: 0.25em 0; }

      blockquote {
        margin: 0.5em 0;
        padding: 4px 12px;
        border-left: 3px solid var(--el-color-primary);
        color: var(--text-secondary);
        background: var(--el-fill-color-blank);
        border-radius: 0 6px 6px 0;
      }

      a {
        color: var(--el-color-primary);
        text-decoration: underline;
      }

      code {
        font-family: "SF Mono", "JetBrains Mono", Consolas, "Courier New", monospace;
        font-size: 0.9em;
        background: var(--el-fill-color);
        padding: 1px 5px;
        border-radius: 4px;
      }

      pre {
        margin: 0.6em 0;
        padding: 10px 12px;
        background: var(--el-fill-color-dark);
        border-radius: 8px;
        overflow-x: auto;
        white-space: pre;
        line-height: 1.5;

        code {
          background: transparent;
          padding: 0;
        }
      }

      table {
        margin: 0.6em 0;
        border-collapse: collapse;
        font-size: 0.95em;

        th, td {
          border: 1px solid var(--el-border-color-lighter);
          padding: 6px 10px;
        }

        th {
          background: var(--el-fill-color);
          font-weight: 600;
        }
      }

      hr {
        border: none;
        border-top: 1px solid var(--el-border-color-lighter);
        margin: 0.8em 0;
      }

      /* Mermaid：未渲染时展示源码，渲染后由 mermaid 生成 SVG */
      .mermaid {
        margin: 0.6em 0;
        text-align: center;
        font-family: "SF Mono", "JetBrains Mono", Consolas, monospace;
        font-size: 12px;
        line-height: 1.5;
        white-space: pre-wrap;
        word-break: break-all;
        background: var(--el-fill-color);
        border-radius: 8px;
        padding: 10px;
        overflow-x: auto;

        svg {
          max-width: 100%;
          height: auto;
          display: inline-block;
        }
      }
    }
  }
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
</style>
