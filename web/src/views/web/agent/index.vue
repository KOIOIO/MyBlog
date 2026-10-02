<template>
  <div class="agent-page">
    <web-navbar :noScroll="true"/>
    <div class="page">
      <div class="container">
        <!-- 未登录引导 -->
        <div v-if="!userStore.isLoggedIn" class="login-guide">
          <div class="guide-card">
            <div class="guide-icon">🤖</div>
            <h2 class="guide-title">{{ t('pages.agent.title') }}</h2>
            <p class="guide-desc">{{ t('pages.agent.loginTip') }}</p>
            <el-button type="primary" size="large" @click="requireLogin">{{ t('pages.agent.loginBtn') }}</el-button>
          </div>
        </div>

        <!-- 对话工作区 -->
        <div v-else class="workspace">
          <!-- 左侧会话列表 -->
          <el-drawer
              v-model="drawerVisible"
              :title="t('pages.agent.history')"
              size="280px"
              :with-header="false"
              class="agent-drawer"
          >
            <div class="side-inner">
              <ConversationSide :conversations="agentStore.conversations"
                                :current-id="agentStore.currentConversationId"
                                @select="onSelectConversation"
                                @remove="onRemoveConversation"/>
            </div>
          </el-drawer>

          <aside class="side" :class="{collapsed: sideCollapsed}">
            <div class="side-inner">
              <el-button class="new-chat-btn" type="primary" plain style="width: 100%" @click="startNewChat">
                ＋ {{ t('pages.agent.newChat') }}
              </el-button>
              <ConversationSide :conversations="agentStore.conversations"
                                :current-id="agentStore.currentConversationId"
                                @select="onSelectConversation"
                                @remove="onRemoveConversation"/>
            </div>
          </aside>

          <!-- 右侧对话区 -->
          <div class="chat-panel">
            <div class="chat-head">
              <el-button text class="collapse-btn" @click="sideCollapsed = !sideCollapsed">
                <el-icon><Expand v-if="sideCollapsed"/><Fold v-else/></el-icon>
              </el-button>
              <span class="chat-title">{{ currentTitle }}</span>
            </div>

            <MessageList :messages="agentStore.messages" :streaming="agentStore.streaming"/>

            <div class="input-area">
              <div v-if="selectedArticles.length > 0" class="article-tags">
                <el-tag v-for="id in selectedArticles" :key="id" closable type="primary" @close="removeArticle(id)">
                  #{{ id }}
                </el-tag>
              </div>
              <div class="input-row">
                <el-button text class="pick-btn" @click="openPicker">
                  <el-icon><Document/></el-icon>
                  {{ t('pages.agent.selectArticle') }}
                </el-button>
                <el-input
                    v-model="inputText"
                    type="textarea"
                    :rows="1"
                    :autosize="{minRows: 1, maxRows: 5}"
                    resize="none"
                    :placeholder="t('pages.agent.placeholder')"
                    @keydown.enter.exact.prevent="onSend"
                />
                <el-button v-if="agentStore.streaming" type="danger" class="send-btn" @click="agentStore.stopStreaming()">
                  {{ t('pages.agent.stop') }}
                </el-button>
                <el-button v-else type="primary" class="send-btn" :disabled="!inputText.trim()" @click="onSend">
                  {{ t('pages.agent.send') }}
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ArticlePicker ref="pickerRef" @confirm="onPickerConfirm"/>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref, watch} from "vue";
import {ElMessage} from "element-plus";
import {Document, Expand, Fold} from "@element-plus/icons-vue";
import {useI18n} from "vue-i18n";
import WebNavbar from "@/components/layout/WebNavbar.vue";
import ArticlePicker from "@/components/agent/ArticlePicker.vue";
import MessageList from "@/components/agent/MessageList.vue";
import ConversationSide from "@/components/agent/ConversationSide.vue";
import {useAgentStore} from "@/stores/agent";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import type {AgentConversation} from "@/api/agent";

const {t} = useI18n();
const agentStore = useAgentStore();
const userStore = useUserStore();
const layoutStore = useLayoutStore();

const inputText = ref('');
const selectedArticles = ref<number[]>([]);
const pickerRef = ref<InstanceType<typeof ArticlePicker> | null>(null);
const sideCollapsed = ref(false);
const drawerVisible = ref(false);

const currentTitle = computed(() => {
    const cur = agentStore.conversations.find((c) => c.id === agentStore.currentConversationId);
    return cur ? cur.title : t('pages.agent.title');
});

const requireLogin = (): void => {
    layoutStore.state.popoverVisible = true;
    layoutStore.state.loginVisible = true;
};

const initWhenLoggedIn = async (): Promise<void> => {
    await agentStore.loadConversations();
    if (agentStore.conversations.length > 0 && agentStore.currentConversationId === 0) {
        await agentStore.openConversation(agentStore.conversations[0].id);
    }
};

onMounted(() => {
    if (userStore.isLoggedIn) {
        initWhenLoggedIn();
    }
});

watch(() => userStore.isLoggedIn, (logged) => {
    if (logged && agentStore.conversations.length === 0) {
        initWhenLoggedIn();
    }
});

const startNewChat = (): void => {
    agentStore.currentConversationId = 0;
    agentStore.messages = [];
    selectedArticles.value = [];
    inputText.value = '';
};

const onSelectConversation = (id: number): void => {
    agentStore.openConversation(id);
    drawerVisible.value = false;
};

const onRemoveConversation = async (id: number): Promise<void> => {
    await agentStore.removeConversation(id);
};

const onSend = async (): Promise<void> => {
    const content = inputText.value.trim();
    if (!content || agentStore.streaming) {
        return;
    }
    inputText.value = '';
    try {
        await agentStore.sendMessage(content, selectedArticles.value);
        selectedArticles.value = [];
    } catch (e) {
        ElMessage.error(e instanceof Error ? e.message : t('common.error'));
    }
};

const openPicker = (): void => {
    pickerRef.value?.open(selectedArticles.value);
};

const onPickerConfirm = (ids: number[]): void => {
    selectedArticles.value = ids;
};

const removeArticle = (id: number): void => {
    selectedArticles.value = selectedArticles.value.filter((x) => x !== id);
};
</script>

<style scoped lang="scss">
.agent-page {
  .page {
    .container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 16px;
    }
  }

  .login-guide {
    display: flex;
    justify-content: center;
    padding: 80px 16px;

    .guide-card {
      text-align: center;
      padding: 48px 64px;
      border-radius: 16px;
      background: var(--el-bg-color);
      box-shadow: var(--el-box-shadow-light);

      .guide-icon { font-size: 56px; }
      .guide-title { margin: 16px 0 8px; font-size: 22px; color: var(--text-body); }
      .guide-desc { margin-bottom: 24px; color: var(--text-secondary); }
    }
  }

  .workspace {
    display: flex;
    gap: 16px;
    height: calc(100vh - 140px);
    min-height: 480px;
  }

  .side {
    width: 260px;
    flex-shrink: 0;
    background: var(--el-bg-color);
    border-radius: 12px;
    box-shadow: var(--el-box-shadow-light);
    overflow: hidden;

    &.collapsed { display: none; }

    .side-inner {
      display: flex;
      flex-direction: column;
      gap: 12px;
      height: 100%;
      padding: 12px;
      box-sizing: border-box;
    }
  }

  .chat-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: var(--el-bg-color);
    border-radius: 12px;
    box-shadow: var(--el-box-shadow-light);
    overflow: hidden;

    .chat-head {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 10px 16px;
      border-bottom: 1px solid var(--el-border-color-lighter);

      .collapse-btn { display: none; }
      .chat-title { font-size: 15px; font-weight: 600; color: var(--text-body); }
    }

    .input-area {
      padding: 12px 16px;
      border-top: 1px solid var(--el-border-color-lighter);

      .article-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;
        margin-bottom: 8px;
      }

      .input-row {
        display: flex;
        align-items: flex-end;
        gap: 8px;

        .pick-btn { flex-shrink: 0; }
        .send-btn { flex-shrink: 0; }

        :deep(.el-textarea__inner) {
          padding: 8px 12px;
          font-size: 14px;
          line-height: 1.5;
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .agent-page {
    .workspace { height: calc(100vh - 120px); }
    .side {
      display: none;
      &.collapsed { display: none; }
    }
    .chat-panel .chat-head .collapse-btn { display: inline-flex; }
  }
}
</style>
